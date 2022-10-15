package codegen

import (
	"log"
	"strconv"
	"strings"
)

// This file does lots of option processing.  Typically, one is using these options
// to "match" the binary interface of somebody else's WASM module.
const (
	fileOptionForAbi         = "543210"
	serviceOptionForWasmName = "543210"
	messageOptionForWasmName = "543210"
	methodOptionForWasmName  = "543210"
	fieldOptionForWasmName   = "543210"

	serviceOptionNoPackage        = "543211"
	serviceOptionAlwaysPullParams = "543212"

	methodOptionPullParams = "543211"
	methodOptionCallsAbi   = "543212"

	messageOptionNoPackage = "543211"
)

// options to map converts the text string that is the options for a given level
// of the proto file and parses into a map.  Note that you can have file options
// service options, field options, etc.
func optionsToMap(s string) map[string]string {
	parts := strings.Split(s, " ")
	result := make(map[string]string)
	for _, opt := range parts {
		if strings.TrimSpace(opt) == "" {
			continue
		}
		assign := strings.Split(opt, ":")
		if len(assign) != 2 {
			log.Printf("unable to understand option: %s ", opt)
			continue
		}
		k := assign[0]
		v := assign[1]
		result[k] = v
	}
	return result
}

// IsAbi checks to see if the options given on the file has an option which looks like:
// option (parigot.abi) = true. This is a free standing function so it can be called
// when you not yet in code generation and thus don't have access to the GenInfo.
func IsAbi(s string) bool {
	_, b := isBooleanOptionPresent(s, fileOptionForAbi)
	return b
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
	return "", false
}

// isWasmServiceName looks for the option wasm_service_name inside the given string.
func isWasmServiceName(s string) (string, bool) {
	return isStringOptionPresent(s, serviceOptionForWasmName)
}

// isWasmMessageName looks for the option wasm_message_name inside the given string.
func isWasmMessageName(s string) (string, bool) {
	return isStringOptionPresent(s, messageOptionForWasmName)
}

func hasNoPackageOption(s string) bool {
	_, b := isBooleanOptionPresent(s, serviceOptionNoPackage)
	return b
}

func alwaysPullParamsOption(s string) bool {
	_, b := isBooleanOptionPresent(s, serviceOptionAlwaysPullParams)
	return b
}

func pullParamsOption(s string) bool {
	_, b := isBooleanOptionPresent(s, methodOptionPullParams)
	return b
}

func abiCallOption(s string) bool {
	_, b := isBooleanOptionPresent(s, methodOptionCallsAbi)
	return b
}
