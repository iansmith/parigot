package codegen

import (
	"log"
	"strconv"
	"strings"
)

// This file does lots of option processing.
const (
	serviceOptionForWasmName          = "543210"
	serviceOptionForReverseAPI        = "543214"
	serviceOptionImplementsReverseAPI = "543215"

	messageOptionForWasmName = "543210"
	methodOptionForWasmName  = "543210"
	fieldOptionForWasmName   = "543210"

	parigotOptionForEnum = "543211"

	parigotOptionForHostFuncName = "543212"
	parigotOptionForErrorIdName  = "543213"
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

// isReverseAPI looks for the option that expresses that the currently being defined
// service is not to be generated associated with the code that defines it.
func isWasmServiceReverseAPI(s string) bool {
	_, ok := isBooleanOptionPresent(s, serviceOptionForReverseAPI)
	return ok
}

// implementsReverseAPI returns the type that this service truly implements,
// not the code which is defined inline.
func implementsReverseAPI(s string) string {
	s, ok := isStringOptionPresent(s, serviceOptionImplementsReverseAPI)
	if !ok {
		return ""
	}
	return s
}

// isWasmServiceName looks for the option wasm_service_name inside the given string.
func isWasmServiceName(s string) (string, bool) {
	return isStringOptionPresent(s, serviceOptionForWasmName)
}

// isWasmMessageName looks for the option wasm_message_name inside the given string.
func isWasmMessageName(s string) (string, bool) {
	return isStringOptionPresent(s, messageOptionForWasmName)
}
