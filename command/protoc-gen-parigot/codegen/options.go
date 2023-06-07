package codegen

import (
	"log"
	"strconv"
	"strings"
)

// This file does lots of option processing.
const (
	serviceOptionForWasmName = "543210"
	messageOptionForWasmName = "543210"
	methodOptionForWasmName  = "543210"
	fieldOptionForWasmName   = "543210"

	parigotOptionForService = "543211"
	parigotOptionForMessage = "543211"
	parigotOptionForEnum    = "543211"

	// serviceOptionNoPackage        = "543211"
	// serviceOptionAlwaysPullParams = "543212"
	// serviceOptionAlwaysPullOutput = "543213"
	// serviceOptionKernel           = "543214"

	// methodOptionPullParams = "543211"
	// methodOptionPullOutput = "543213"

	// if these flags are used in conjunction with a command line flag the generator should output
	// code for the test methods mentioned in the service defintion.
	// serviceOptionServiceTest = "543215"
	// methodOptionMethodTest   = "543216"
	// serviceOptionErrId       = "543217"
)

// options to map converts the text string that is the options for a given level
// of the proto file and parses into a map.  Note that you can have file options
// service options, field options, etc.
func optionsToMap(s string) map[string]string {
	if s == "<nil>" {
		return nil
	}
	parts := strings.Split(s, " ")
	result := make(map[string]string)
	for _, opt := range parts {
		if strings.TrimSpace(opt) == "" {
			continue
		}
		assign := strings.Split(opt, ":")
		if len(assign) != 2 {
			log.Printf("unable to understand option: '%s' from original '%s' 1st:'%s' (%d)", opt, s, parts[0], len(assign))
			continue
		}
		k := assign[0]
		v := assign[1]
		result[k] = v
	}
	return result
}

// isBooleanOptionPresent does all the string futzing to find an desired option or return false beacuse
// it isn't there.  It returns the option as the first parameter, but its not likely
// you'll care.
func isBooleanOptionPresent(s, target string) (string, bool) {
	m := optionsToMap(s)
	text, ok := m[target]
	if ok {
		value, err := strconv.Atoi(text)
		if err != nil {
			panic("bad value supplied to us by protobuf compiler for our option:" + err.Error())
		}
		return text, value != 0
	}
	return "", false
}

// isStringOptionPresent does all the string futzing to find an desired option or return false because
// it isn't there.  It returns the value of the option as the first parameter.
func isStringOptionPresent(s, target string) (string, bool) {
	m := optionsToMap(s)
	text, ok := m[target]
	if ok {
		return text, true
	}
	return text, false
}

// isWasmServiceName looks for the option wasm_service_name inside the given string.
func isWasmServiceName(s string) (string, bool) {
	return isStringOptionPresent(s, serviceOptionForWasmName)
}

// isWasmServiceErrId looks for the option wasm_err_id inside the given string.
// func isWasmServiceErrId(s string) (string, bool) {
// 	return isStringOptionPresent(s, serviceOptionErrId)
// }

// isWasmMessageName looks for the option wasm_message_name inside the given string.
func isWasmMessageName(s string) (string, bool) {
	return isStringOptionPresent(s, messageOptionForWasmName)
}

func isServiceMarkedParigot(s string) bool {
	_, ok := isBooleanOptionPresent(s, messageOptionForWasmName)
	return ok
}
func IsEnumMarkedParigot(s string) bool {
	s, ok := isBooleanOptionPresent(s, parigotOptionForEnum)
	if !ok {
		return false
	}
	return s != ""
}

// func hasNoPackageOption(s string) bool {
// 	_, b := isBooleanOptionPresent(s, serviceOptionNoPackage)
// 	return b
// }

// func hasServiceTestOption(s string) bool {
// 	_, b := isBooleanOptionPresent(s, serviceOptionServiceTest)
// 	return b
// }
// func hasMethodTestOption(s string) bool {
// 	_, b := isBooleanOptionPresent(s, methodOptionMethodTest)
// 	return b
// }

// func hasKernelOption(s string) bool {
// 	_, b := isBooleanOptionPresent(s, serviceOptionKernel)
// 	return b
// }

// func alwaysPullParamsOption(s string) bool {
// 	_, b := isBooleanOptionPresent(s, serviceOptionAlwaysPullParams)
// 	return b
// }

// func alwaysPullOutputOption(s string) bool {
// 	_, b := isBooleanOptionPresent(s, serviceOptionAlwaysPullOutput)
// 	return b
// }

// func pullParamsOption(s string) bool {
// 	_, b := isBooleanOptionPresent(s, methodOptionPullParams)
// 	return b
// }
// func pullOutputOption(s string) bool {
// 	_, b := isBooleanOptionPresent(s, methodOptionPullOutput)
// 	return b
// }
