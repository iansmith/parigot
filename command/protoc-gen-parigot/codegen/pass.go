package codegen

import "strings"

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
	outParam := NewCGParameterNoFormal(out.GetCGType())

	if in.IsEmpty() {
		result += empty(method, true, inParam)
	}
	if !in.IsEmpty() {
		if !in.GetCGType().IsBasic() {
			//comp := in.GetCGType().GetCompositeType()
			//if len(comp.GetField()) != 0 {
			//	//for num, f := range comp.GetField() {
			//	//	cgt := NewCGTypeFromBasic(f.GetType().String(), method.GetParent().GetLanguage(),
			//	//		method.GetParent().GetFinder(),
			//	//		method.GetParent().GetParent().GetPackage())
			//	//	cgp := NewCGParameterFromField(f, cgt)
			//	//	result += fn(method, num, cgp)
			//	//	if num != len(comp.GetField())-1 {
			//	//		result += method.GetFormalArgSeparator()
			//	//	}
			//	//}
			//}
			fakeFormal := in.GetCGType().ShortName()[0:1]
			cgp := NewCGParameterFromString(strings.ToLower(fakeFormal), in.GetCGType())
			result += fn(method, 0, cgp)
		}
	}
	if out.IsEmpty() {
		result += empty(method, false, outParam)
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
							cgt := NewCGTypeFromBasic(f.GetType().String(), svc.GetLanguage(), svc.GetFinder(), protoPackage)
							cgp := NewCGParameterFromField(f, cgt)
							fn(gen, svc, method, cgp)
						}
					}
				}
			}
		}
	}
}
