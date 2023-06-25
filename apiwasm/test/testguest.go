package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/iansmith/parigot/apishared/id"
	qlib "github.com/iansmith/parigot/apiwasm/queue/lib"
	guestsyscall "github.com/iansmith/parigot/apiwasm/syscall"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/queue/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	test "github.com/iansmith/parigot/g/test/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/client"

	"google.golang.org/protobuf/types/known/anypb"
)

const testQueueName = "test_queue"

func main() {
	ctx := pcontext.CallTo(pcontext.GuestContext(pcontext.NewContextWithContainer(context.Background(), "[testwasm]main")), "[testwasm].main")
	myId := test.MustRegisterUnderTest(ctx)
	queue.MustRequire(ctx, myId)
	test.MustExportTest(ctx)

	server := &myTestServer{
		myId: myId,
	}

	binding := test.MustWaitSatisfiedTest(ctx, myId, server)
	launchFuture := test.LaunchTest(ctx, myId, server)
	launchFuture.Handle(func(ok bool) {
		if !ok {
			pcontext.Errorf(ctx, "ready check failed for test service")
			guestsyscall.Exit(1)
		}
	})
	var kerr syscall.KernelErr
	for {
		test.RunTest(ctx, binding, 1000, server)
	}
	pcontext.Fatalf(ctx, "error while waiting for test service calls: %s", syscall.KernelErr_name[int32(kerr)])

}

type myTestServer struct {
	myId      id.ServiceId
	suite     map[string]*suiteInfo
	suiteExec map[string]string
	testQid   queue.QueueId

	queueSvc queue.Client

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

func newSuiteInfo(req *test.AddTestSuiteRequest) ([]*suiteInfo, test.TestErr) {
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
	return infoList, test.TestErr_NoError
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

func (m *myTestServer) Ready(ctx context.Context, sid id.ServiceId) *lib.BaseFuture[bool] {
	// initialization can be done here, not just in main
	m.suite = make(map[string]*suiteInfo)
	m.suiteExec = make(map[string]string)

	pcontext.Debugf(ctx, "myTestServer ready called")
	defer func() {
		if r := recover(); r != nil {
			log.Printf("xxxx MustLocate queue died: %v", r)
		}
	}()
	m.queueSvc = queue.MustLocate(ctx, sid)
	log.Printf("got a queue, with a cs of %s", m.queueSvc.(*queue.Client_).String())
	baseBool := lib.NewBaseFuture[bool]()
	// note that find or create queue returns an id future
	// and it will be zero value if things failed upstream.
	qlib.FindOrCreateQueue(ctx, m.queueSvc, testQueueName).
		Handle(func(i queue.QueueId) {
			if i.IsZeroOrEmptyValue() {
				pcontext.Errorf(ctx, "unable find or create queue in test service")
				baseBool.Set(false)
				return
			}
			m.testQid = i
			baseBool.Set(true)
		})

	return baseBool
}

func (m *myTestServer) AddTestSuite(ctx context.Context, req *test.AddTestSuiteRequest) (*test.AddTestSuiteResponse, test.TestErr) {

	infoList, err := newSuiteInfo(req)
	if err != test.TestErr_NoError { // really should not happen
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

	resp := &test.AddTestSuiteResponse{
		Succeeded: success,
	}
	pcontext.Debugf(ctx, "", "addSuiteResp.Succeeded:%#v", resp.Succeeded)

	return resp, test.TestErr_NoError
}

func contains(list []string, cand string) bool {
	for _, member := range list {
		if member == cand {
			return true
		}
	}
	return false
}

func (m *myTestServer) Start(ctx context.Context, req *test.StartRequest) (*test.StartResponse, test.TestErr) {
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
		return &test.StartResponse{
			RegexFailed: regexpFail,
		}, test.TestErr_RegexpFailed
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
					if err != test.TestErr_NoError {
						return nil, err
					}
					f := m.queueSvc.Send(ctx, sReq)
					f.Success(func(resp *queue.SendResponse) {
						if len(resp.Succeed) != 1 {
							pcontext.Errorf(ctx, "failed send to the queue, expected only one send to succed")
						}
					})
					f.Failure(func(err queue.QueueErr) {
						pcontext.Errorf(ctx, "failed send to the queue: %s", queue.QueueErr_name[int32(err)])
					})
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
	resp := &test.StartResponse{
		NumTest: int32(count),
	}
	return resp, test.TestErr_NoError
}
func (m *myTestServer) Background(ctx context.Context) {
	if !m.started {
		return
	}
	req := queue.ReceiveRequest{
		Id:           m.testQid.Marshal(),
		MessageLimit: 1,
	}
	f := m.queueSvc.Receive(ctx, &req)
	f.Failure(func(err queue.QueueErr) {
		pcontext.Errorf(ctx, "unable to receive from queue: %s", queue.QueueErr_name[int32(err)])
	})
	f.Success(func(resp *queue.ReceiveResponse) {
		msg := resp.Message[0]
		aload := msg.GetPayload()
		payload := test.QueuePayload{}
		unmarsh := aload.UnmarshalTo(&payload)
		if unmarsh != nil {
			pcontext.Errorf(ctx, "unable to unmarshal queue message payload: %s", unmarsh.Error())
			return
		}
		pcontext.Infof(ctx, "got test from queue %s,%s", payload.Name, payload.FuncName)
	})
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

func (m *myTestServer) runTests(ctx context.Context, fullTestName, execPackageSvc string) *test.FutureUnderTestExec {
	pcontext.Debugf(ctx, "run tests0")
	locate := make(map[string]test.ClientUnderTest)
	var client test.ClientUnderTest
	part := strings.Split(fullTestName, ".")

	pkg, svc := splitPkgAndService(strings.Join(part[:len(part)-1], "."))
	//name := part[len(part)-1]
	pcontext.Debugf(ctx, "run tests2 %s,%s\n", fullTestName, execPackageSvc)
	loc, ok := locate[execPackageSvc]
	if !ok {
		var err test.TestErr
		execPkg, execSvc := splitPkgAndService(execPackageSvc)
		client, err = m.locateClient(ctx, execPkg, execSvc)
		if err != test.TestErr_NoError {
			pcontext.Errorf(ctx, "unable to test %s.%s, cannot locate it", execPkg, execSvc)
			f := test.NewFutureUnderTestExec()
			f.CompleteCall(nil, int32(err))
			return f
		}

	} else {
		client = loc
	}
	pcontext.Debugf(ctx, "run tests4, got a client")
	f := client.UnderTestExec(ctx, &test.ExecRequest{
		Package: pkg,
		Service: svc,
		Name:    part[len(part)-1],
	})
	f.Failure(func(err test.TestErr) {
		pcontext.Errorf(ctx, "xxx run tests %v", test.TestErr_name[int32(err)])
	})
	f.Success(func(resp *test.ExecResponse) {
		pcontext.Debugf(ctx, "xxx run tests %s.%s.%s (skipped? %v, success? %v)",
			resp.GetPackage(), resp.GetService(), resp.GetName(), resp.GetSkipped(), resp.GetSuccess())
	})
	return f
}

func splitPkgAndService(s string) (string, string) {
	part := strings.Split(s, ".")

	pkg := strings.Join(part[:len(part)-1], ".")
	svc := part[len(part)-1]
	return pkg, svc

}

func (m *myTestServer) locateClient(ctx context.Context, pkg, svc string) (test.ClientUnderTest, test.TestErr) {
	cs, err := client.LocateDynamic(ctx, pkg, svc, m.myId)
	if err != syscall.KernelErr_NoError {
		return nil, test.TestErr_DynamicLocate
	}
	return &test.ClientUnderTest_{
		BaseService: cs,
	}, test.TestErr_NoError
}

// makeSendRequest creates a SendRequest and all the internal objects required
// make it work correctly in the test queue.
func makeSendRequest(ctx context.Context, qid queue.QueueId, name, funcName string) (*queue.SendRequest, test.TestErr) {

	qidM := qid.Marshal()
	payload := test.QueuePayload{
		Name:     name,
		FuncName: funcName,
	}
	a := anypb.Any{}
	if err := a.MarshalFrom(&payload); err != nil {
		pcontext.Errorf(ctx, "unable to marshal payload creating send request: %s", err.Error())
		return nil, test.TestErr_Marshal
	}
	msg := queue.QueueMsg{
		Payload: &a,
		Id:      nil,
	}
	req := queue.SendRequest{
		Id:  qidM,
		Msg: []*queue.QueueMsg{&msg},
	}
	return &req, test.TestErr_NoError
}
