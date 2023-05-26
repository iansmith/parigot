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
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	testmsg "github.com/iansmith/parigot/g/msg/test/v1"
	"github.com/iansmith/parigot/g/queue/v1"
	"github.com/iansmith/parigot/g/test/v1"
	lib "github.com/iansmith/parigot/lib/go"

	"google.golang.org/protobuf/types/known/anypb"
)

const testQueueName = "test_queue"

func main() {
	lib.FlagParseCreateEnv()

	// The queue we will use
	queue.RequireQueueServiceOrPanic(context.Background())
	test.ExportTestServiceOrPanic()
	test.RunTestService(&myTestServer{})
}

type myTestServer struct {
	suite     map[string]*suiteInfo
	suiteExec map[string]string
	testQid   id.Id

	queueSvc queue.QueueServiceClient

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

func newSuiteInfo(req *testmsg.AddTestSuiteRequest) ([]*suiteInfo, error) {
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
	return infoList, nil
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
	m.queueSvc = queue.LocateQueueServiceOrPanic(ctx)
	qid, err := m.findOrCreateQueue(ctx, testQueueName)
	if err != nil {
		panic("myTestServer: unable to create the test queue: " + err.Error())
	}
	m.testQid = qid

	return true
}

func (m *myTestServer) findOrCreateQueue(ctx context.Context, name string) (id.Id, error) {
	req := queuemsg.LocateRequest{}
	req.QueueName = testQueueName
	resp, err := m.queueSvc.Locate(&req)
	if err != nil {
		panic("myTestServer: ready: error in attempt to get queue by name: " + err.Error())
	}
	qid := id.Unmarshal(resp.Id)
	if qid.IsError() {
		return nil, id.NewPerrorFromId("Unmarshal:", qid)
	}
	// it's a not found, so create it
	createReq := queuemsg.CreateQueueRequest{}
	createReq.QueueName = name
	createResp, err := m.queueSvc.CreateQueue(&createReq)
	if err != nil {
		return nil, id.NewPerrorFromId("CreateQueue", id.NewQueueError(id.QueueInternalError))
	}
	qid = id.Unmarshal(createResp.GetId())
	return qid, nil
}

func (m *myTestServer) AddTestSuite(ctx context.Context, req *testmsg.AddTestSuiteRequest) (*testmsg.AddTestSuiteResponse, error) {

	infoList, err := newSuiteInfo(req)
	if err != nil { // really should not happen
		pcontext.Errorf(ctx, "unable to create suite info: %s", err.Error())
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

	return resp, nil
}

func contains(list []string, cand string) bool {
	for _, member := range list {
		if member == cand {
			return true
		}
	}
	return false
}

func (m *myTestServer) Start(ctx context.Context, req *testmsg.StartRequest) (*testmsg.StartResponse, error) {
	var err error
	var regexpFail bool
	var suiteString, nameString string

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
		}, err
	}
	count := 0
	m.started = true // lets go
	if m.haveSuite {
		for _, suiteInfo := range m.suite {
			name := suiteInfoToSuiteName(suiteInfo)
			match := m.suiteRegex.MatchString(name)
			if match {
				count += len(suiteInfo.funcName)
				m.addAllTests(suiteInfo)
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
					sReq, err := makeSendRequest(m.testQid, name, m.suiteExec[suiteName])
					if err != nil {
						return nil, err
					}
					resp, err := m.queueSvc.Send(sReq)
					if err != nil {
						return nil, err
					}
					if len(resp.Succeed) != 1 {
						return nil, fmt.Errorf("unable to send correct number of messages: %d", len(resp.Succeed))
					}
				}
			}
		}
	}
	if !m.haveSuite && !m.haveName {
		for _, suite := range m.suite {
			count += len(suite.funcName)
			m.addAllTests(suite)
		}
	}
	resp := &testmsg.StartResponse{
		NumTest: int32(count),
	}
	return resp, nil
}
func (m *myTestServer) Background(ctx context.Context) {
	if !m.started {
		return
	}
	req := queuemsg.ReceiveRequest{
		Id:           id.Marshal[protosupportmsg.QueueId](m.testQid),
		MessageLimit: 1,
	}
	resp, err := m.queueSvc.Receive(&req)
	if err != nil {
		pcontext.Errorf(ctx, "unable to receive from queue: %s, retrying...", err.Error())
		return
	}
	msg := resp.Message[0]
	aload := msg.GetPayload()
	payload := testmsg.QueuePayload{}
	err = aload.UnmarshalTo(&payload)
	if err != nil {
		pcontext.Errorf(ctx, "unable to unmarshal queue message payload: %s, retrying...", err.Error())
		return
	}
	pcontext.Logf(ctx, pcontext.Info, false, "got test from queue %s,%s", payload.Name, payload.FuncName)
}

func (m *myTestServer) addAllTests(info *suiteInfo) {
	base := suiteInfoToSuiteName(info)
	for _, name := range info.funcName {
		testName := fmt.Sprintf("%s.%s", base, name)
		makeSendRequest(m.testQid, base, testName)
	}
}

func suiteAndTestToTestName(info *suiteInfo, fn string) string {
	return fmt.Sprintf("%s.%s", suiteInfoToSuiteName(info), fn)
}

func suiteInfoToSuiteName(info *suiteInfo) string {
	return fmt.Sprintf("%s.%s", info.pkg, info.service)
}

func (m *myTestServer) runTests(ctx context.Context, fullTestName, execPackageSvc string) (id.Id, string) {
	print(fmt.Sprintf("run tests0\n"))
	locate := make(map[string]test.UnderTestServiceClient)
	var client test.UnderTestServiceClient
	part := strings.Split(fullTestName, ".")

	pkg, svc := splitPkgAndService(strings.Join(part[:len(part)-1], "."))
	//name := part[len(part)-1]
	print(fmt.Sprintf("run tests2 %s,%s\n", fullTestName, execPackageSvc))
	loc, ok := locate[execPackageSvc]
	if !ok {
		execPkg, execSvc := splitPkgAndService(execPackageSvc)
		client = m.locateClient(ctx, execPkg, execSvc)
		locate[execPackageSvc] = client
	} else {
		client = loc
	}
	print(fmt.Sprintf("run tests4, got a client\n"))
	resp, err := client.Exec(&testmsg.ExecRequest{Package: pkg, Service: svc, Name: part[len(part)-1]})
	if err != nil {
		print(fmt.Sprintf("xxx run tests %v\n", err.Error()))
		return id.NewTestError(id.TestErrorServiceNotFound), err.Error()
	}
	print("xxx run tests %s.%s.%s (skipped? %v, success? %v)",
		resp.GetPackage(), resp.GetService(), resp.GetName(), resp.GetSkipped(), resp.GetSuccess())
	return nil, ""
}
func splitPkgAndService(s string) (string, string) {
	part := strings.Split(s, ".")

	pkg := strings.Join(part[:len(part)-1], ".")
	svc := part[len(part)-1]
	return pkg, svc

}

func (m *myTestServer) locateClient(ctx context.Context, pkg, svc string) test.UnderTestServiceClient {
	print(fmt.Sprintf("xxx locate test pkg=%s svc=%s\n", pkg, svc))
	req := &syscallmsg.LocateRequest{
		PackageName: pkg,
		ServiceName: svc,
	}
	print(fmt.Sprintf("xxx locate client 1\n"))
	resp, err := syscall.Locate(req)
	if err != nil {
		panic("locate failed for " + pkg + "." + svc + "\n")
	}
	print(fmt.Sprintf("xxx locate client 3\n"))
	service := id.Unmarshal(resp.GetServiceId())
	cs := lib.NewClientSideService(ctx, service, "testService")

	pcontext.Debugf(ctx, "locateClient", "xxx locate client 4")
	return &test.UnderTestServiceClient_{
		ClientSideService: cs,
	}
}

// makeSendRequest creates a SendRequest and all the internal objects required
// make it work correctly in the test queue.
func makeSendRequest(qid id.Id, name, funcName string) (*queuemsg.SendRequest, error) {
	qidM := id.Marshal[protosupportmsg.QueueId](qid)
	payload := testmsg.QueuePayload{
		Name:     name,
		FuncName: funcName,
	}
	a := anypb.Any{}
	if err := a.MarshalFrom(&payload); err != nil {
		return nil, fmt.Errorf("unable to marshal test payload: %v", err)
	}
	msg := queuemsg.QueueMsg{
		Payload: &a,
		Id:      qidM,
	}
	req := queuemsg.SendRequest{
		Id:  qidM,
		Msg: []*queuemsg.QueueMsg{&msg},
	}
	return &req, nil
}
