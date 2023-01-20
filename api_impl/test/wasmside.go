package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/iansmith/parigot/api_impl/syscall"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	testmsg "github.com/iansmith/parigot/g/msg/test/v1"
	"github.com/iansmith/parigot/g/test/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/backdoor"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var callImpl = syscall.NewCallImpl()

func main() {
	lib.FlagParseCreateEnv()

	// you need to put Require and Export calls in here, but put Run() call in Ready()
	if _, err := callImpl.Export1("test", "TestService"); err != nil {
		panic("myLogServer:ready: error in attempt to export api.Log: " + err.Error())
	}
	test.RunTestService(&myTestServer{})
}

type myTestServer struct {
	// within a test server, the requests are handled sequentially
	suite   map[string]*suiteInfo
	test    map[string]string
	started bool
}

type suiteInfo struct {
	pkg, service string
	nameToFunc   map[string]string
}

func newSuiteInfo(req *testmsg.AddTestSuiteRequest) ([]*suiteInfo, error) {
	infoList := []*suiteInfo{}

	for _, suite := range req.GetSuite() {
		result := &suiteInfo{
			pkg:     suite.GetPackagePath(),
			service: suite.GetService(),
		}
		result.nameToFunc = make(map[string]string)
		for name, fn := range suite.GetNameToFunc() {
			result.nameToFunc[name] = fn
		}
		infoList = append(infoList, result)
	}
	return infoList, nil
}

func (s *suiteInfo) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s.%s{", s.pkg, s.service))
	for k, v := range s.nameToFunc {
		fmt.Sprintf("[%s=%s]", k, v)
	}
	buf.WriteString("}")
	return buf.String()
}

func (m *myTestServer) Ready() bool {
	// initialization needs to be done here, not in main
	m.suite = make(map[string]*suiteInfo)
	m.test = make(map[string]string)

	if _, err := callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		panic("myTestServer: ready: error in attempt to signal Run: " + err.Error())
	}
	return true
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
		for testName, fn := range suite.nameToFunc {
			print(fmt.Sprintf("xxx in testName loop %s fn=%s", testName, fn))
			test := fmt.Sprintf("%s.%s", suiteName, testName)
			if seenSuiteBefore {
				if _, seenTest := oldSuite.nameToFunc[testName]; seenTest {
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

func (m *myTestServer) Start(_ *protosupportmsg.Pctx, in proto.Message) (proto.Message, error) {
	req := in.(*testmsg.StartRequest)
	suiteRegexpString, nameRegexpString := "", ""
	haveSuite, haveName := false, false
	var err error
	var suiteRegex, nameRegex *regexp.Regexp
	var regexpFail bool

	if req.GetFilterSuite() != "" {
		suiteRegexpString = req.GetFilterSuite()
		haveSuite = true
	}
	if haveSuite {
		suiteRegex, err = regexp.Compile(suiteRegexpString)
		if err != nil {
			regexpFail = true
		}
	} else {
		if req.GetFilterSuite() != "" {
			nameRegexpString = req.GetFilterName()
			haveName = true
		}
		if haveName {
			nameRegex, err = regexp.Compile(nameRegexpString)
			if err != nil {
				regexpFail = true
			}
		}
	}
	count := 0
	m.started = true // lets go
	if haveSuite {
		for _, suiteInfo := range m.suite {
			name := suiteInfoToSuiteName(suiteInfo)
			match := suiteRegex.MatchString(name)
			if match {
				count += len(suiteInfo.nameToFunc)
				m.addAllTests(suiteInfo)
			}
		}
	}
	if haveName {
		for _, suiteInfo := range m.suite {
			suiteName := suiteInfoToSuiteName(suiteInfo)
			for name, fn := range suiteInfo.nameToFunc {
				testName := fmt.Sprintf("%s.%s", suiteName, name)
				match := nameRegex.MatchString(testName)
				if match {
					count++
					m.test[name] = fn
				}
			}
		}
	}
	if !haveSuite && !haveName {
		for _, suite := range m.suite {
			count += len(suite.nameToFunc)
			m.addAllTests(suite)
		}
	}
	resp := &testmsg.StartResponse{
		RegexFailed: regexpFail,
		NumTest:     int32(count),
	}
	logDebug(fmt.Sprintf("all tests to run:%#v", m.test))
	m.runTests()
	return resp, nil
}

func (m *myTestServer) addAllTests(info *suiteInfo) {
	base := suiteInfoToSuiteName(info)
	for name, fn := range info.nameToFunc {
		testName := fmt.Sprintf("%s.%s", base, name)
		m.test[testName] = fn
	}
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

func (m *myTestServer) runTests() {
	call := syscall.NewCallImpl()

	locate := make(map[string]*test.UnderTestServiceClient)
	var client *test.UnderTestServiceClient
	for svcPart, testPart := range m.test {
		part := strings.Split(svcPart, ".")
		print(fmt.Sprintf("xxx test part svcPart='%s' testPart='%s' %+v\n", svcPart, testPart, part))

		pkg := strings.Join(part[:len(part)-2], ".")
		svc := part[len(part)-2]
		//name := part[len(part)-1]
		loc, ok := locate[svcPart]
		if ok && loc == badLocate {
			continue
		}
		if !ok {
			client = m.locateClient(pkg, svc, call)
			locate[svcPart] = client
		} else {
			client = locate[svcPart]
		}
		if client == badLocate {
			continue
		}
		resp, err := client.Exec(&testmsg.ExecRequest{Package: pkg, Service: svc, Name: testPart})
		if err != nil {
			print(fmt.Sprintf("xxx run tests %v\n", err.Error()))
			continue
		}
		print("xxx run tests %s.%s.%s (skipped? %v, success? %v)",
			resp.GetPackage(), resp.GetService(), resp.GetName(), resp.GetSkipped(), resp.GetSuccess())
	}
}

func (m *myTestServer) locateClient(pkg, svc string, call lib.Call) *test.UnderTestServiceClient {
	print(fmt.Sprintf("xxx locate test pkg=%s svc=%s\n", pkg, svc))
	req := &syscallmsg.LocateRequest{
		PackageName: pkg,
		ServiceName: svc,
	}
	resp, err := call.Locate(req)
	if err != nil {
		print("error caught in the locate code:" + err.Error() + "\n")
		return badLocate
	}
	if resp.GetServiceId() == nil {
		print("locate failed for " + pkg + "." + svc + "\n")
		return badLocate
	}
	service := lib.Unmarshal(resp.GetServiceId())
	cs := lib.NewClientSideService(service, "testService", nil, callImpl)

	return &test.UnderTestServiceClient{
		ClientSideService: cs,
		Call:              syscall.NewCallImpl(),
	}
}
