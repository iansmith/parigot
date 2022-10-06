package transform

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type TypeDescriptor interface {
	Name() string
	NoResult() bool
	NumResult() int
	Result(int) string
	AllResult() []string
	NoParam() bool
	NumParam() int
	Param(int) string
	AllParam() []string
	HasStringParam() bool
	APIWrapper() string
	NumTrueParams() int
}

type tdesc struct {
	name      string
	param     []string
	result    []string
	trueParam []string
}

// NewTypeDescriptor returns a Typedescriptor for a type that has the sequence
// of parameters and return result types provided.  Either of these may be nil.
func NewTypeDescriptor(n string, p []string, r []string, t []string) TypeDescriptor {
	return &tdesc{
		name:      n,
		param:     p,
		result:    r,
		trueParam: t,
	}
}
func (t *tdesc) Name() string             { return t.name }
func (t *tdesc) NoResult() bool           { return len(t.result) == 0 }
func (t *tdesc) NumResult() int           { return len(t.result) }
func (t *tdesc) Result(i int) string      { return t.result[i] }
func (t *tdesc) AllResult() []string      { return t.result }
func (t *tdesc) NoParam() bool            { return len(t.param) == 0 }
func (t *tdesc) NumParam() int            { return len(t.param) }
func (t *tdesc) Param(i int) string       { return t.param[i] }
func (t *tdesc) AllParam() []string       { return t.param }
func (t *tdesc) IsStringParam(i int) bool { return t.trueParam[i] == "string" }
func (t *tdesc) NumTrueParams() int       { return len(t.trueParam) }
func (t *tdesc) HasStringParam() bool {
	for i := 0; i < len(t.trueParam); i++ {
		if t.trueParam[i] == "string" {
			return true
		}
	}
	return false
}

// Number for type descriptor returns the int that corresponds to that descriptor
// in this module.  If necessary, it will create a new type number and attach it to
// the module so this promise can always be true.
func (m *Module) NumberForTypeDescriptor(descriptor TypeDescriptor) int {
	last := -1
	for i, tl := range m.Code {
		if tl.TopLevelType() != TypeDefT {
			continue
		}
		last = i
		typeDef := tl.(*TypeDef)
		if typeDef.compare(descriptor) == true {
			return i
		}
	}
	if last == -1 {
		panic("no typedefs found in this module")
	}
	var ret *ResultDef
	var param *ParamDef
	if !descriptor.NoParam() {
		param = &ParamDef{
			Type: &TypeNameSeq{
				Name: descriptor.AllParam(),
			},
		}
	}
	if !descriptor.NoResult() {
		ret = &ResultDef{
			Type: &TypeNameSeq{
				Name: descriptor.AllResult(),
			},
		}
	}
	newDef := &TypeDef{
		Annotation: last,
		Func: &FuncSpec{
			Param:  param,
			Result: ret,
		},
	}
	m.Code = append(m.Code[:last], append([]TopLevel{newDef}, m.Code[last:]...)...)
	return newDef.Annotation
}

func (t *TypeDef) compare(descriptor TypeDescriptor) bool {
	if descriptor.NoResult() && t.Func.Result != nil {
		return false
	}
	if !descriptor.NoResult() && t.Func.Result == nil {
		return false
	}
	if descriptor.NoParam() && t.Func.Param != nil {
		return false
	}
	if !descriptor.NoParam() && t.Func.Param == nil {
		return false
	}
	// we know now that we can just test descriptor for having either of these
	// since we would have already failed otherwise
	if !descriptor.NoParam() {
		if descriptor.NumParam() != len(t.Func.Param.Type.Name) {
			return false
		}
		for i := 0; i < descriptor.NumParam(); i++ {
			if descriptor.Param(i) != t.Func.Param.Type.Name[i] {
				return false
			}
		}
	}
	if !descriptor.NoResult() {
		if descriptor.NumResult() != len(t.Func.Result.Type.Name) {
			return false
		}
		for i := 0; i < descriptor.NumResult(); i++ {
			if descriptor.Result(i) != t.Func.Param.Type.Name[i] {
				return false
			}
		}
	}
	return true
}

func FuncToDescriptor(fn interface{}) TypeDescriptor {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic("FuncToDescriptor can only be called on functions")
	}
	if t.IsVariadic() {
		panic("unable to process variable functions")
	}

	var in []string
	var out []string
	var trueParam []string
	if t.NumIn() > 0 {
		for i := 0; i < t.NumIn(); i++ {
			p := t.In(i)
			switch p.Kind() {
			case reflect.Int, reflect.Int64:
				in = append(in, "i64")
				trueParam = append(trueParam, "i64")
			case reflect.Int32:
				in = append(in, "i32")
				trueParam = append(trueParam, "i32")
			case reflect.Float32:
				in = append(in, "f32")
				trueParam = append(trueParam, "f32")
			case reflect.Float64:
				in = append(in, "f64")
				trueParam = append(trueParam, "f64")
			case reflect.String: //ptr + len
				trueParam = append(trueParam, "string")
				in = append(in, "i32", "i32")
			default:
				panic("unable to handle type in conversion " + p.String())
			}
		}
	}
	if t.NumOut() > 0 {
		for i := 0; i < t.NumOut(); i++ {
			r := t.Out(i)
			switch r.Kind() {
			case reflect.Int, reflect.Int64:
				out = append(out, "i64")
			case reflect.Int32:
				out = append(out, "i32")
			case reflect.Float32:
				out = append(out, "f32")
			case reflect.Float64:
				out = append(out, "f32")
			case reflect.String: //ptr + len
				out = append(out, "i32", "i32")
			default:
				panic("unable to handle type in conversion " + r.String())
			}
		}
	}
	return &tdesc{
		name:      getFunctionName(fn),
		param:     in,
		result:    out,
		trueParam: trueParam,
	}
}

//result["parigot_abi.OutputStringStub"] = wasmtime.WrapFunc(store, func(p1 int32, p2 int32) {
//	abi_impl.OutputString(strConvert(memPtr, p1, p2))
//})

func (t *tdesc) APIWrapper() string {
	var buf bytes.Buffer
	parts := strings.Split(t.Name(), "/")
	shortName := parts[len(parts)-1]
	buf.WriteString("\tlinkage[\"")
	buf.WriteString(fmt.Sprintf("%s\"] =", abiImplToAbiForLinkage(shortName)))
	buf.WriteString("wasmtime.WrapFunc(store," + "func(")
	for i := 0; i < t.NumParam(); i++ {
		buf.WriteString(fmt.Sprintf("p%d %s",
			i, expandTypenameForGo(t.Param(i))))
		if i != t.NumParam()-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	if t.NumResult() > 1 {
		panic("cant generate code for number of results other than 0 or 1")
	}
	if t.NumResult() == 1 {
		buf.WriteString(expandTypenameForGo(t.Result(0)))
	}
	buf.WriteString("{" + newline(2))
	if !t.NoResult() {
		buf.WriteString("result:=")
	}
	buf.WriteString(shortName + "(")
	count := 0
	for i := 0; i < t.NumTrueParams(); i++ {
		if t.IsStringParam(i) {
			buf.WriteString(fmt.Sprintf("strConvert(memPtr,"+
				"int32(uintptr(p%d)+memPtr),"+
				"int32(uintptr(p%d)+memPtr+uintptr(p%d)))",
				count, count+1, count))
			count += 2
		} else {
			buf.WriteString(fmt.Sprintf("p%d", count))
			count++
		}
		if i != t.NumTrueParams()-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	if t.NumResult() > 0 {
		buf.WriteString(newline(2) + " return result")
	}
	buf.WriteString(newline(1) + "})")
	buf.WriteString(newline(0))
	return buf.String()
}

func getFunctionName(fn interface{}) string {
	//strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()), ".")
	//return strs[len(strs)-1]
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}
func newline(numTabs int) string {
	t := "\n"
	for i := 0; i < numTabs; i++ {
		t += "\t"
	}
	return t
}

func expandTypenameForGo(in string) string {
	switch in {
	case "i32":
		return "int32"
	case "i64":
		return "int64"
	case "f32":
		return "float32"
	case "f64":
		return "float64"
	}
	panic("unknown type from type dscriptor:" + in)
}

func abiImplToAbiForLinkage(s string) string {
	return strings.Replace(s, "abi_impl", "parigot_abi", 1)
}
