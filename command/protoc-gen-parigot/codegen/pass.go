package codegen

import (
	"log"
	"strings"
)

func FuncParamPass(method *WasmMethod,
	fn func(method *WasmMethod, num int, parameter *CGParameter) string,
	empty func(method *WasmMethod, isInput bool, param *CGParameter) string) string {

	result := ""
	in := method.GetInputParam()
	out := method.GetOutputParam()

	// check the types to make sure they are not there
	if in.GetCGType() == nil {
		in.cgType = GetCGTypeForInputParam(in)
	}
	if out.GetCGType() == nil {
		out.cgType = GetCGTypeForOutputParam(out)
	}
	inParam := NewCGParameterNoFormal(in.GetCGType())

	if in.IsEmpty() {
		result += empty(method, true, inParam)
	}
	if !in.IsEmpty() {
		if !in.GetCGType().IsBasic() {
			if method.PullParameters() {
				result += walkParametersPulled(method, fn)
			} else {
				fakeFormal := in.GetCGType().ShortName()[0:1]
				cgp := NewCGParameterFromString(strings.ToLower(fakeFormal), in.GetCGType())
				result += fn(method, 0, cgp)
			}
		}
	}
	// output is handled by OutTypeWalk
	return result
}

func ExpandReturnInfoForOutput(out *OutputParam, m *WasmMethod, protoPkg string) *CGParameter {
	t := out.GetCGType()
	if t == nil {
		t = NewCGTypeFromOutput(out, m, protoPkg)
		out.SetCGType(t)
	}
	if t.IsEmpty() {
		return nil
	}
	if t.IsBasic() {
		log.Fatalf("unable to pull parameters from a basic type on input:%s ", t.ShortName())
	}
	comp := t.GetCompositeType()
	if len(comp.GetField()) == 0 {
		return nil
	}
	if len(comp.GetField()) > 1 {
		log.Fatalf("unable to pull parameters up for a return value with more than 1 field:%s ", t.ShortName())
	}
	return NewCGParameterFromField(comp.GetField()[0], m, protoPkg)
}

func ExpandParamInfoForInput(in *InputParam, m *WasmMethod, protoPkg string) []*CGParameter {
	t := in.GetCGType()
	if in.IsEmpty() {
		return nil
	}
	if t.IsBasic() {
		log.Fatalf("unable to pull parameters from a basic type on input:%s ", t.ShortName())
	}
	comp := t.GetCompositeType()
	field := comp.GetField()
	if len(field) == 0 {
		return nil
	}
	result := make([]*CGParameter, len(field))
	for i, f := range field {
		result[i] = NewCGParameterFromField(f, m, protoPkg)
	}
	return result
}

func MethodsPass(gen *GenInfo,
	fn func(info *GenInfo, svc *WasmService, method *WasmMethod, parameter *CGParameter) string,
	empty func(info *GenInfo, svc *WasmService, method *WasmMethod, isInput bool, param *CGParameter) string) {
	for _, svc := range gen.Service() {
		protoPackage := svc.GetParent().GetPackage()
		for _, method := range svc.GetWasmMethod() {
			in := method.GetInputParam()
			out := method.GetOutputParam()

			// check the types to make sure they are not there
			if in.GetCGType() == nil {
				in.cgType = GetCGTypeForInputParam(in)
			}
			if out.GetCGType() == nil {
				out.cgType = GetCGTypeForOutputParam(out)
			}
			inParam := NewCGParameterNoFormal(in.GetCGType())
			outParam := NewCGParameterNoFormal(out.GetCGType())

			if !in.IsEmpty() {
				empty(gen, svc, method, true, inParam)
			}
			if !out.IsEmpty() {
				empty(gen, svc, method, false, outParam)
			}
			if !in.IsEmpty() {
				if !in.GetCGType().IsBasic() {
					comp := in.GetCGType().GetCompositeType()
					if len(comp.GetField()) != 0 {
						for _, f := range comp.GetField() {
							cgp := NewCGParameterFromField(f, method, protoPackage)
							fn(gen, svc, method, cgp)
						}
					}
				}
			}
		}
	}
}

func walkParametersPulled(m *WasmMethod,
	fn func(method *WasmMethod, num int, parameter *CGParameter) string) string {
	in := m.GetCGInput()
	protoPkg := m.GetParent().GetProtoPackage()
	paramList := ExpandParamInfoForInput(in, m, protoPkg)
	result := ""
	for i, cgp := range paramList {
		result += fn(m, i, cgp)
		if i != len(paramList)-1 {
			result += m.GetLanguage().GetFormalArgSeparator()
		}
	}
	return result
}

func walkOutputPulledParam(m *WasmMethod,
	fn func(protoPkg string, method *WasmMethod, parameter *CGParameter) string) string {
	out := m.GetCGOutput()
	if out.GetCGType().IsEmpty() {
		// this is a semi-error but we tolerate it... you could have specified
		// alwaysPullUp on the service and then if you had any empty functions
		// you'd get a bunch of errors
		return ""
	}
	t := out.GetCGType()
	if t.IsBasic() {
		log.Fatalf("unable to pull up parameters of %s because it is not a composite "+
			"object", t.String("" /*doesnt matter on basic*/))
	}
	composite := t.GetCompositeType()
	if len(composite.GetField()) == 0 {
		return ""
	}
	if len(t.GetCompositeType().GetField()) > 1 {
		log.Fatalf("cant pull up parameters from output type %s, it has more than 1 value",
			t.String(""))
	}
	protoPkg := m.GetParent().GetProtoPackage()
	f := t.GetCompositeType().GetField()[0]
	childType := NewCGTypeFromField(f, m, protoPkg)
	cgParam := NewCGParameterNoFormal(childType)
	return fn(m.GetProtoPackage(), m, cgParam)
}
