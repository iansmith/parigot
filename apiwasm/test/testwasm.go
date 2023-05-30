package main

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/apiwasm/syscall"
	pcontext "github.com/iansmith/parigot/context"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	testmsg "github.com/iansmith/parigot/g/msg/test/v1"
	"github.com/iansmith/parigot/g/queue/v1"
	queueg "github.com/iansmith/parigot/g/queue/v1"
	"github.com/iansmith/parigot/g/test/v1"
	testg "github.com/iansmith/parigot/g/test/v1"
	lib "github.com/iansmith/parigot/lib/go"

	"google.golang.org/protobuf/types/known/anypb"
)

const testQueueName = "test_queue"

func main() {
	lib.FlagParseCreateEnv()
	ctx := pcontext.CallTo(pcontext.ServerWasmContext(pcontext.NewContextWithContainer(context.Background(), "[testwasm]main")), "[testwasm].main")

	myId := testg.MustRegisterTestService(ctx)
	// The queue we will use
	queueg.MustRequireQueueService(ctx, myId)
	testg.MustExportTestService(ctx)
	testg.RunTestService(ctx, &myTestServer{})
}

type myTestServer struct {
	suite     map[string]*suiteInfo
	suiteExec map[string]string
	testQid   queueg.QueueId

	queueSvc queueg.QueueServiceClient

	started    bool
	haveName   bool
	haveSuite  bool
	suiteRegex *regexp.Regexp
	nameRegex  *regexp.Regexp
}

type suiteInfo struct {
	pkg, service string
	funcName     []string
}

func newSuiteInfo(req *testmsg.AddTestSuiteRequest) ([]*suiteInfo, testg.TestErrId) {
	infoList := []*suiteInfo{}

	for _, suite := range req.GetSuite() {
		result := &suiteInfo{
			pkg:     suite.GetPackagePath(),
			service: suite.GetService(),
		}
		result.funcName = make([]string, len(suite.GetFunctionName()))
		for i := 0; i < len(suite.GetFunctionName()); i++ {
			result.funcName[i] = suite.GetFunctionName()[i]
		}
		infoList = append(infoList, result)
	}
	return infoList, testg.TestErrIdNoErr
}

func (s *suiteInfo) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s.%s[", s.pkg, s.service))
	for i, f := range s.funcName {
		buf.WriteString(fmt.Sprintf("%s", f))
		if i != len(s.funcName)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func (m *myTestServer) Ready(ctx context.Context) bool {
	// initialization needs to be done here, not in main
	m.suite = make(map[string]*suiteInfo)
	m.suiteExec = make(map[string]string)

	pcontext.Debugf(ctx, "Ready", "myTestServer ready called")
	test.WaitTestServiceOrPanic()
	m.queueSvc = queueg.MustLocateQueueService(ctx)
	qid, err := m.findOrCreateQueue(ctx, testQueueName)
	if err.IsError() {
		pcontext.Errorf(ctx, "myTestServer: failed to extract queue ID ", err.Short())
		return false
	}
	m.testQid = qid

	return true
}

func (m *myTestServer) findOrCreateQueue(ctx context.Context, name string) (queueg.QueueId, queueg.QueueErrId) {
	req := queuemsg.LocateRequest{}
	req.QueueName = testQueueName
	resp, err := m.queueSvc.Locate(ctx, &req)
	if err.IsError() {
		pcontext.Errorf(ctx, "myTestServer: ready: error in attempt to get queue by name: ", err.Short())
		return queueg.ZeroValueQueueId(), err
	}
	qid, idErr := queue.UnmarshalQueueId(resp.Id)
	if idErr.IsError() {
		pcontext.Errorf(ctx, "myTestServer: failed to extract queue ID ", err.Short())
		return queueg.ZeroValueQueueId(), err
	}
	// it's a not found, so create it
	createReq := queuemsg.CreateQueueRequest{}
	createReq.QueueName = name
	createResp, err := m.queueSvc.CreateQueue(ctx, &createReq)
	if err.IsError() {
		return queueg.ZeroValueQueueId(), err
	}
	qid, unmarsh := queueg.UnmarshalQueueId(createResp.GetId())
	if unmarsh.IsError() {
		return queueg.ZeroValueQueueId(), queueg.NewQueueErrId(queue.QueueErrIdUnmarshalError)
	}
	return qid, queueg.QueueErrIdNoErr
}

func (m *myTestServer) AddTestSuite(ctx context.Context, req *testmsg.AddTestSuiteRequest) (*testmsg.AddTestSuiteResponse, testg.TestErrId) {

	infoList, err := newSuiteInfo(req)
	if err.IsError() { // really should not happen
		return nil, err
	}
	success := make(map[string]bool)
	for _, suite := range infoList {
		var oldSuite *suiteInfo
		suiteName := suiteInfoToSuiteName(suite)
		seenSuiteBefore := false
		alreadyHave, ok := m.suite[suiteName]
		if ok {
			seenSuiteBefore = true
			oldSuite = alreadyHave
		} else {
			m.suite[suiteName] = suite
		}
		m.suiteExec[suiteName] = fmt.Sprintf("%s.%s", req.GetExecPackage(), req.GetExecService())
		for _, testName := range suite.funcName {
			test := fmt.Sprintf("%s.%s", suiteName, testName)
			if seenSuiteBefore {
				if contains(oldSuite.funcName, testName) {
					success[test] = false
					continue
				}
			}
			success[test] = true
			continue
		}
	}

	resp := &testmsg.AddTestSuiteResponse{
		Succeeded: success,
	}
	pcontext.Debugf(ctx, "", "addSuiteResp.Succeeded:%#v", resp.Succeeded)

	return resp, testg.TestErrIdNoErr
}

func contains(list []string, cand string) bool {
	for _, member := range list {
		if member == cand {
			return true
		}
	}
	return false
}

func (m *myTestServer) Start(ctx context.Context, req *testmsg.StartRequest) (*testmsg.StartResponse, testg.TestErrId) {
	var regexpFail bool
	var suiteString, nameString string
	var err error
	if req.GetFilterSuite() != "" {
		suiteString = req.GetFilterSuite()
		m.haveSuite = true
	}
	if m.haveSuite {
		m.suiteRegex, err = regexp.Compile(suiteString)
		if err != nil {
			regexpFail = true
		}
	} else {
		if req.GetFilterSuite() != "" {
			nameString = req.GetFilterName()
			m.haveName = true
		}
		if m.haveName {
			m.nameRegex, err = regexp.Compile(nameString)
			if err != nil {
				regexpFail = true
			}
		}
	}
	if regexpFail {
		return &testmsg.StartResponse{
			RegexFailed: regexpFail,
		}, testg.NewTestErrId(TestErrorrRegexpFailed)
	}
	count := 0
	m.started = true // lets go
	if m.haveSuite {
		for _, suiteInfo := range m.suite {
			name := suiteInfoToSuiteName(suiteInfo)
			match := m.suiteRegex.MatchString(name)
			if match {
				count += len(suiteInfo.funcName)
				m.addAllTests(pcontext.CallTo(ctx, "addAllTests"), suiteInfo)
			}
		}
	}
	if m.haveName {
		for _, suiteInfo := range m.suite {
			suiteName := suiteInfoToSuiteName(suiteInfo)
			for _, name := range suiteInfo.funcName {
				testName := fmt.Sprintf("%s.%s", suiteName, name)
				match := m.nameRegex.MatchString(testName)
				if match {
					count++
					sReq, err := makeSendRequest(pcontext.CallTo(ctx, "makeSendRequest"), m.testQid, name, m.suiteExec[suiteName])
					if err.IsError() {
						return nil, err
					}
					resp, errResp := m.queueSvc.Send(ctx, sReq)
					if errResp.IsError() {
						pcontext.Errorf(ctx, "failed send to the queue: %s", errResp.Short())
						return nil, err
					}
					if len(resp.Succeed) != 1 {
						return nil, testg.NewTestErrId(TestErrorInternal)
					}
				}
			}
		}
	}
	if !m.haveSuite && !m.haveName {
		for _, suite := range m.suite {
			count += len(suite.funcName)
			m.addAllTests(pcontext.CallTo(ctx, "addAlTests"), suite)
		}
	}
	resp := &testmsg.StartResponse{
		NumTest: int32(count),
	}
	return resp, testg.TestErrIdNoErr
}
func (m *myTestServer) Background(ctx context.Context) {
	if !m.started {
		return
	}
	req := queuemsg.ReceiveRequest{
		Id:           m.testQid.Marshal(),
		MessageLimit: 1,
	}
	resp, err := m.queueSvc.Receive(ctx, &req)
	if err.IsError() {
		pcontext.Errorf(ctx, "unable to receive from queue: %s", err.Short())
		return
	}
	msg := resp.Message[0]
	aload := msg.GetPayload()
	payload := testmsg.QueuePayload{}
	unmarsh := aload.UnmarshalTo(&payload)
	if err.IsError() {
		pcontext.Errorf(ctx, "unable to unmarshal queue message payload: %s, retrying...", unmarsh.Error())
		return
	}
	pcontext.Infof(ctx, "got test from queue %s,%s", payload.Name, payload.FuncName)
}

func (m *myTestServer) addAllTests(ctx context.Context, info *suiteInfo) {
	base := suiteInfoToSuiteName(info)
	for _, name := range info.funcName {
		testName := fmt.Sprintf("%s.%s", base, name)
		makeSendRequest(pcontext.CallTo(ctx, "makeSendRequest"), m.testQid, base, testName)
	}
}

func suiteAndTestToTestName(info *suiteInfo, fn string) string {
	return fmt.Sprintf("%s.%s", suiteInfoToSuiteName(info), fn)
}

func suiteInfoToSuiteName(info *suiteInfo) string {
	return fmt.Sprintf("%s.%s", info.pkg, info.service)
}

func (m *myTestServer) runTests(ctx context.Context, fullTestName, execPackageSvc string) testg.TestErrId {
	pcontext.Debugf(ctx, "run tests0")
	locate := make(map[string]test.UnderTestServiceClient)
	var client test.UnderTestServiceClient
	part := strings.Split(fullTestName, ".")

	pkg, svc := splitPkgAndService(strings.Join(part[:len(part)-1], "."))
	//name := part[len(part)-1]
	pcontext.Debugf(ctx, "run tests2 %s,%s\n", fullTestName, execPackageSvc)
	loc, ok := locate[execPackageSvc]
	if !ok {
		var tId testg.TestErrId
		execPkg, execSvc := splitPkgAndService(execPackageSvc)
		client, tId = m.locateClient(pcontext.CallTo(ctx, "locateClient"), execPkg, execSvc)
		if tId.IsError() {
			return tId
		}
		locate[execPackageSvc] = client
	} else {
		client = loc
	}
	pcontext.Debugf(ctx, "run tests4, got a client")
	resp, err :=
		client.Exec(ctx, &testmsg.ExecRequest{
			Package: pkg,
			Service: svc,
			Name:    part[len(part)-1],
		})
	if err.IsError() {
		pcontext.Errorf(ctx, "xxx run tests %v", err.Short())
		return testg.NewTestErrId(TestErrorServiceNotFound)
	}
	pcontext.Debugf(ctx, "xxx run tests %s.%s.%s (skipped? %v, success? %v)",
		resp.GetPackage(), resp.GetService(), resp.GetName(), resp.GetSkipped(), resp.GetSuccess())
	return testg.TestErrIdNoErr
}
func splitPkgAndService(s string) (string, string) {
	part := strings.Split(s, ".")

	pkg := strings.Join(part[:len(part)-1], ".")
	svc := part[len(part)-1]
	return pkg, svc

}

func (m *myTestServer) locateClient(ctx context.Context, pkg, svc string) (testg.UnderTestServiceClient, testg.TestErrId) {
	pcontext.Debugf(ctx, "xxx locate test pkg=%s svc=%s\n", pkg, svc)
	req := &syscallmsg.LocateRequest{
		PackageName: pkg,
		ServiceName: svc,
	}
	pcontext.Debugf(ctx, "xxx locate client 1")
	resp, err := syscall.Locate(req)
	if err.IsError() {
		pcontext.Errorf(ctx, "locate failed for %s.%s", pkg, svc)
		return nil, testg.NewTestErrId(TestErrorServiceNotFound)
	}
	pcontext.Debugf(ctx, "xxx locate client 3")
	service, idErr := id.UnmarshalServiceId(resp.GetServiceId())
	if idErr.IsError() {
		return nil, testg.NewTestErrId(TestErrorMarshaling)
	}
	cs := lib.NewClientSideService(ctx, service, "testService")

	pcontext.Debugf(ctx, "locateClient", "xxx locate client 4")
	return &testg.UnderTestServiceClient_{
		ClientSideService: cs,
	}, testg.TestErrIdNoErr
}

// makeSendRequest creates a SendRequest and all the internal objects required
// make it work correctly in the test queue.
func makeSendRequest(ctx context.Context, qid queueg.QueueId, name, funcName string) (*queuemsg.SendRequest, testg.TestErrId) {

	qidM := qid.Marshal()
	payload := testmsg.QueuePayload{
		Name:     name,
		FuncName: funcName,
	}
	a := anypb.Any{}
	if err := a.MarshalFrom(&payload); err != nil {
		pcontext.Errorf(ctx, "unable to marshal payload creating send request: %s", err.Error())
		return nil, testg.NewTestErrId(TestErrorMarshaling)
	}
	msg := queuemsg.QueueMsg{
		Payload: &a,
		Id:      nil,
	}
	req := queuemsg.SendRequest{
		Id:  qidM,
		Msg: []*queuemsg.QueueMsg{&msg},
	}
	return &req, testg.TestErrIdNoErr
}
