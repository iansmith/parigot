package abi

import (
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
	"google.golang.org/protobuf/types/descriptorpb"
	"log"
	"text/template"
)

var abiFuncMap = template.FuncMap{
	"inputParamNameSet":       inputParamNameSet,
	"isInputEmpty":            isInputEmpty3,
	"isOutputEmpty":           isOutputEmpty4,
	"outputTypeName":          outputTypeName,
	"outputZeroDefaultReturn": outputZeroDefaultReturn,
}

const typeString = "TYPE_STRING"

var validWASMTypes = []string{
	typeString, "TYPE_INT32", "TYPE_INT64", "TYPE_FLOAT", "TYPE_DOUBLE", "TYPE_BOOL"}

var strictWASMTypes = []string{"TYPE_INT32", "TYPE_INT64", "TYPE_FLOAT", "TYPE_DOUBLE"}

var globalDescriptor *descriptorpb.FileDescriptorProto

func isValidWASMType(x string) bool {
	for _, t := range validWASMTypes {
		if t == x {
			return true
		}
	}
	return false
}
func isStrictWASMType(x string) bool {
	for _, t := range strictWASMTypes {
		if t == x {
			return true
		}
	}
	return false
}

// because of func map restrictions on return values, have to do this
func isInputEmpty3(m *descriptorpb.MethodDescriptorProto) bool {
	_, _, r, _ := isInputOutputEmpty(m)
	return r
}

// because of func map restrictions on return values, have to do this
func isOutputEmpty4(m *descriptorpb.MethodDescriptorProto) bool {
	_, _, _, r := isInputOutputEmpty(m)
	return r
}

func isInputOutputEmpty(m *descriptorpb.MethodDescriptorProto) (*descriptorpb.DescriptorProto,
	*descriptorpb.DescriptorProto, bool, bool) {
	inTarget := util.LastSegmentOfFQProtoPackageName(m.GetInputType())
	outTarget := util.LastSegmentOfFQProtoPackageName(m.GetOutputType())
	inEmpty := false
	outEmpty := false
	var in *descriptorpb.DescriptorProto
	var out *descriptorpb.DescriptorProto
	for _, msg := range globalDescriptor.GetMessageType() {
		if msg.GetName() == inTarget {
			if len(msg.Field) == 0 {
				inEmpty = true
			}
			in = msg
		}
		if msg.GetName() == outTarget {
			if len(msg.Field) == 0 {
				outEmpty = true
			}
			if len(msg.Field) > 1 {
				log.Printf("unable to process multiple output values in %s", msg.GetName())
			}
			out = msg
		}
	}
	if !inEmpty && in == nil {
		log.Fatalf("unable to find type named %s from the input parameter of %s",
			m.GetInputType(), m.GetName())
	}
	if !outEmpty && out == nil {
		log.Fatalf("unable to find type named %s from the output parameter of %s",
			m.GetOutputType(), m.GetName())
	}
	return in, out, inEmpty, outEmpty
}
func inputParamNameSet(noUsages bool, noParams bool, m *descriptorpb.MethodDescriptorProto) string {
	input, _, empty, _ := isInputOutputEmpty(m)
	if empty {
		return ""
	}
	result := ""
	for _, f := range input.Field {
		if !isValidWASMType(f.Type.String()) {
			log.Fatalf("currenly, abi functions must be one of the basic wasm types (int32,int64,float,double) or \"string\": %s", f.Type.String())
		}
	}

	for i, f := range input.Field {
		n := f.GetName() + " "
		if noUsages {
			n = "_ "
		}
		if noParams {
			n = ""
		}
		switch f.Type.String() {
		case "TYPE_STRING":
			result = n + "string"
		case "TYPE_INT32":
			result += n + "int32"
		case "TYPE_INT64":
			result += n + "int64"
		case "TYPE_FLOAT":
			result += n + "float32"
		case "TYPE_DOUBLE":
			result += n + "float64"
		case "TYPE_BOOL":
			result += n + "bool"
		}
		if i != len(input.Field)-1 {
			result += ","
		}
	}
	return result
}

func outputZeroDefaultReturn(m *descriptorpb.MethodDescriptorProto) string {
	_, output, _, empty := isInputOutputEmpty(m)
	if empty {
		panic("should not be computing the return value for a function with no return declared") //should not happen, checked earlier
	}
	for _, f := range output.Field {
		if !isStrictWASMType(f.Type.String()) {
			panic("currently, can only generate return values for basic wasm types") //should not happen, checked earlier
		}
	}
	result := "return "
	for i, f := range output.Field {
		switch f.Type.String() {
		case "TYPE_INT32":
			result += "int32(0)"
		case "TYPE_INT64":
			result += "int64(0)"
		case "TYPE_FLOAT":
			result += "float32(0.0)"
		case "TYPE_DOUBLE":
			result += "float64(0.0)"
		case "TYPE_BOOL":
			result += "int32(0) //bool"
		case "TYPE_STRING":
			result += "\"\" // string"
		}
		if i != len(output.Field)-1 {
			result += ","
		}
	}
	return result

}

func outputTypeName(m *descriptorpb.MethodDescriptorProto) string {
	_, output, _, empty := isInputOutputEmpty(m)
	if empty {
		return ""
	}
	for _, f := range output.Field {
		if !isStrictWASMType(f.Type.String()) {
			log.Fatalf("currenly, ABI functions must return one of the four basic types in WASM (int32,int64,float or double): %s", f.Type.String())
		}
	}
	result := ""
	for i, f := range output.Field {
		switch f.Type.String() {
		case "TYPE_INT32":
			result += " int32"
		case "TYPE_INT64":
			result += " int64"
		case "TYPE_FLOAT":
			result += " float32"
		case "TYPE_DOUBLE":
			result += " float64"
		}
		if i != len(output.Field)-1 {
			result += ","
		}
	}
	return result
}
