package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/iansmith/parigot/apiimpl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"

	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	testmsg "github.com/iansmith/parigot/g/msg/test/v1"
	"github.com/iansmith/parigot/g/queue/v1"
	"github.com/iansmith/parigot/g/test/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/backdoor"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var callImpl = syscall.NewCallImpl()

const testQueueName = "test_queue"

func main() {
	lib.FlagParseCreateEnv()

	// The queue we will use
	if _, err := callImpl.Require1("queue", "QueueService"); err != nil {
		panic("myLogServer:ready: error in attempt to require queue.QueueService: " + err.Error())
	}
	// The logger we will use
	if _, err := callImpl.Require1("log", "Log"); err != nil {
		panic("myLogServer:ready: error in attempt to require log.Log: " + err.Error())
	}

	// you need to put Require and Export calls in here, but put Run() call in Ready()
	if _, err := callImpl.Export1("test", "TestService"); err != nil {
		panic("myLogServer:ready: error in attempt to export test.TestService: " + err.Error())
	}
	test.RunTestService(&myTestServer{})
}

type myTestServer struct {
	suite     map[string]*suiteInfo
	suiteExec map[string]string
	testQid   lib.Id

	queueSvc queue.QueueService
	logger   log.LogService

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

func (m *myTestServer) Ready() bool {
	// initialization needs to be done here, not in main
	m.suite = make(map[string]*suiteInfo)
	m.suiteExec = make(map[string]string)

	var err error

	if _, err = callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		panic("myTestServer: ready: error in attempt to signal Run: " + err.Error())
	}
	m.logger, err = log.LocateLogService()
	if err != nil {
		panic("myTestServer: ready: error in attempt to get logger: " + err.Error())
	}
	m.queueSvc, err = queue.LocateQueueService(m.logger)
	if err != nil {
		panic("myTestServer: ready: error in attempt to get queue service: " + err.Error())
	}
	qid, deets := m.findOrCreateQueue(testQueueName)
	if qid.IsErrorType() {
		panic("myTestServer: unable to create the test queue: " + deets)
	}

	return true
}

func (m *myTestServer) findOrCreateQueue(name string) (lib.Id, string) {
	req := queuemsg.LocateRequest{}
	req.QueueName = testQueueName
	resp, err := m.queueSvc.Locate(&req)
	if err != nil {
		panic("myTestServer: ready: error in attempt to get queue service: " + err.Error())
	}
	qid := lib.Unmarshal(resp.Id)
	if qid.IsErrorType() && qid.ErrorCode() != lib.QueueNotFound {
		return qid, err.Error()
	}
	if !qid.IsErrorType() {
		return qid, ""
	}
	// it's a not found, so create it
	createReq := queuemsg.CreateQueueRequest{}
	createReq.QueueName = name
	createResp, err := m.queueSvc.CreateQueue(&createReq)
	if err != nil {
		return lib.NewQueueError(lib.QueueInternalError), err.Error()
	}
	qid = lib.Unmarshal(createResp.GetId())
	return qid, ""
}

func (m *myTestServer) AddTestSuite(_ *protosupportmsg.Pctx, in proto.Message) (proto.Message, error) {
	req := in.(*testmsg.AddTestSuiteRequest)

	infoList, err := newSuiteInfo(req)
	if err != nil { // really should not happen
		logError("unable to create suite info", err)
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
	logDebug(fmt.Sprintf("addSuiteResp.Succeeded:%#v", resp.Succeeded))

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

func (m *myTestServer) Start(_ *protosupportmsg.Pctx, in proto.Message) (proto.Message, error) {
	req := in.(*testmsg.StartRequest)
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
func (m *myTestServer) Background() {
	if !m.started {
		return
	}
	req := queuemsg.ReceiveRequest{
		Id:           lib.Marshal[protosupportmsg.QueueId](m.testQid),
		MessageLimit: 1,
	}
	resp, err := m.queueSvc.Receive(&req)
	if err != nil {
		m.logger.Log(&logmsg.LogRequest{
			Stamp:   timestamppb.Now(),
			Level:   logmsg.LogLevel_LOG_LEVEL_ERROR,
			Message: fmt.Sprintf("unable to pull test message from queue:%v", err),
		})
		return
	}
	msg := resp.Message[0]
	aload := msg.GetPayload()
	payload := testmsg.QueuePayload{}
	err = aload.UnmarshalTo(&payload)
	if err != nil {
		m.logger.Log(&logmsg.LogRequest{
			Stamp:   timestamppb.Now(),
			Level:   logmsg.LogLevel_LOG_LEVEL_ERROR,
			Message: fmt.Sprintf("unable to unmarshal test message from queue:%v", err),
		})
		return
	}
	m.logger.Log(&logmsg.LogRequest{
		Stamp:   timestamppb.Now(),
		Level:   logmsg.LogLevel_LOG_LEVEL_INFO,
		Message: fmt.Sprintf("got test from queue: %s,%s", payload.Name, payload.FuncName),
	})

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

func logDebug(msg string) {
	backdoor.Log(&logmsg.LogRequest{
		Message: msg,
		Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
		Stamp:   timestamppb.Now(),
	}, false, true, false, nil)
}

func logError(msg string, err error) {
	backdoor.Log(&logmsg.LogRequest{
		Message: fmt.Sprintf("%s:%s", msg, err.Error()),
		Level:   logmsg.LogLevel_LOG_LEVEL_ERROR,
		Stamp:   timestamppb.Now(),
	}, false, true, false, nil)
}

var badLocate = &test.UnderTestServiceClient{}

func (m *myTestServer) runTests(fullTestName, execPackageSvc string) (lib.Id, string) {
	call := syscall.NewCallImpl()
	print(fmt.Sprintf("run tests0\n"))
	locate := make(map[string]*test.UnderTestServiceClient)
	var client *test.UnderTestServiceClient
	part := strings.Split(fullTestName, ".")

	pkg, svc := splitPkgAndService(strings.Join(part[:len(part)-1], "."))
	//name := part[len(part)-1]
	print(fmt.Sprintf("run tests2 %s,%s\n", fullTestName, execPackageSvc))
	loc, ok := locate[execPackageSvc]
	if ok && loc == badLocate {
		print(fmt.Sprintf("run tests3a, not able to locate\n"))
		return lib.NewTestError(lib.TestErrorServiceNotFound), "unable to locate " + execPackageSvc
	}
	if !ok {
		execPkg, execSvc := splitPkgAndService(execPackageSvc)
		client = m.locateClient(execPkg, execSvc, call)
		locate[execPackageSvc] = client
	} else {
		client = loc
	}
	print(fmt.Sprintf("run tests4, got a client\n"))
	if client == badLocate {
		return lib.NewTestError(lib.TestErrorServiceNotFound), "unable to locate " + execPackageSvc
	}
	resp, err := client.Exec(&testmsg.ExecRequest{Package: pkg, Service: svc, Name: part[len(part)-1]})
	if err != nil {
		print(fmt.Sprintf("xxx run tests %v\n", err.Error()))
		return lib.NewTestError(lib.TestErrorServiceNotFound), err.Error()
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

func (m *myTestServer) locateClient(pkg, svc string, call lib.Call) *test.UnderTestServiceClient {
	print(fmt.Sprintf("xxx locate test pkg=%s svc=%s\n", pkg, svc))
	req := &syscallmsg.LocateRequest{
		PackageName: pkg,
		ServiceName: svc,
	}
	print(fmt.Sprintf("xxx locate client 1\n"))
	resp, err := call.Locate(req)
	if err != nil {
		print("error caught in the locate code:" + err.Error() + "\n")
		return badLocate
	}
	print(fmt.Sprintf("xxx locate client 2\n"))
	if resp.GetServiceId() == nil {
		print("locate failed for " + pkg + "." + svc + "\n")
		return badLocate
	}
	print(fmt.Sprintf("xxx locate client 3\n"))
	service := lib.Unmarshal(resp.GetServiceId())
	cs := lib.NewClientSideService(service, "testService", nil, callImpl)

	print(fmt.Sprintf("xxx locate client 4\n"))
	return &test.UnderTestServiceClient{
		ClientSideService: cs,
		Call:              syscall.NewCallImpl(),
	}
}

// makeSendRequest creates a SendRequest and all the internal objects required
// make it work correctly in the test queue.
func makeSendRequest(qid lib.Id, name, funcName string) (*queuemsg.SendRequest, error) {
	qidM := lib.Marshal[protosupportmsg.QueueId](qid)
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
