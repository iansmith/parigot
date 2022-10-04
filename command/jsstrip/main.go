package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/iansmith/parigot/command/transform"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

const (
	WasiModule    = "wasi_snapshot_preview1"
	ParigotModule = "parigot_abi"
	JSModule      = "env"
	FdWrite       = "fd_write"
	wat2wasm      = "wat2wasm"
	wasm2wat      = "wasm2wat"

	parigotFilename      = "parigot-transformed.wat"
	jsNotImplemented     = "$github.com/iansmith/parigot/abi.JSNotImplemented"
	tinygoNotImplemented = "$github.com/iansmith/parigot/abi.TinyGoNotImplemented"
)

// poisoning? gunshots?
var methodsToKillImport = []string{
	"syscall/js.valueGet",
	"syscall/js.valuePrepareString",
	"syscall/js.valueLoadString",
	"syscall/js.finalizeRef",
	"syscall/js.stringVal",
	"syscall/js.valueSet",
	"syscall/js.valueLength",
	"syscall/js.valueIndex",
	"syscall/js.valueCall",
}

func main() {

	if len(os.Args) != 3 {
		log.Fatalf("unable to understand arguments, should provide exactly two arguments (input wasm file, output wasm file)")
	}

	_, baseArg0 := filepath.Split(os.Args[0])
	tmp, err := os.MkdirTemp("", baseArg0)
	if err != nil {
		log.Fatalf("cannot create temp dir: %v", err)
	}

	// their output is our input
	watVersion, err := convertInputToFormat(os.Args[1], tmp, "", wasm2wat, "wasm to wat")
	if err != nil {
		os.Exit(1)
	}
	modifiedWat := parigotProcessing(watVersion, tmp)
	_, err = convertInputToFormat(modifiedWat, tmp, os.Args[2], wat2wasm, "wat to wasm")
	if err != nil {
		os.Exit(1)
	}
	//os.RemoveAll(tmp)
	log.Printf("%s -> %s\n", os.Args[1], watVersion)
	log.Printf("%s -> %s\n", watVersion, filepath.Join(tmp, parigotFilename))
	log.Printf("%s -> %s\n", filepath.Join(tmp, parigotFilename), os.Args[2])
	os.Exit(0)
}

func parigotProcessing(inputFilename, tmp string) string {
	// Set up the input
	fs, err := antlr.NewFileStream(inputFilename)
	if err != nil {
		log.Fatalf("failed trying to open input file, %v", err)
	}
	// make lexer
	lexer := transform.NewWasmLexer(fs)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := transform.NewWasmParser(stream)

	// Finally parse the expression
	builder := &transform.Builder{}
	antlr.ParseTreeWalkerDefault.Walk(builder, p.Module())
	mod := builder.Module() // only one module right now
	strippingPass(mod)
	patchingPass(mod)
	outName := filepath.Join(tmp, parigotFilename)
	out, err := os.Create(outName)
	if err != nil {
		log.Fatalf("unable to create output file: %v", err)
	}
	out.WriteString(mod.IndentedString(0))
	out.Close()
	return outName
}

func patchingPass(mod *transform.Module) {
	result := []transform.TopLevel{}
	// walk all the toplevels
	for _, tl := range mod.Code {
		switch tl.TopLevelType() {
		case transform.FuncDefT:
			tl = processFuncPass2(tl.(*transform.FuncDef))
		default:
			break
		}
		if tl != nil {
			result = append(result, tl)
		}
	}

	mod.Code = result
}

func strippingPass(mod *transform.Module) {
	result := []transform.TopLevel{}
	// walk all the toplevels
	for _, tl := range mod.Code {
		switch tl.TopLevelType() {
		case transform.ImportDefT:
			tl = processImport(tl.(*transform.ImportDef))
		case transform.ExportDefT:
			tl = processExport(tl.(*transform.ExportDef))
		case transform.FuncDefT:
			tl = processFuncPass1(tl.(*transform.FuncDef))
		default:
			break
		}
		if tl != nil {
			result = append(result, tl)
		}
	}
	mod.Code = result
}

func parigotMangleMethod(category, name string) string {
	return fmt.Sprintf("%s_%s")
}

func processImport(importDef *transform.ImportDef) transform.TopLevel {
	if importDef.ImportedAs == WasiModule {
		if importDef.ImportedAs == FdWrite {
			importDef.ModuleName = ParigotModule
			importDef.ImportedAs = parigotMangleMethod("wasi_emulation", "fd_write")
			return importDef
		}
		log.Fatalf("parigot emulation of %s not implemented yet", importDef.ImportedAs)
	}
	if importDef.ModuleName == JSModule {
		for _, s := range methodsToKillImport {
			if importDef.ImportedAs == s {
				return nil // drops the import
			}
		}
		log.Fatalf("to understand what to do with import of %s from module env (js runtime?)", importDef.ImportedAs)
	}
	return importDef
}

func processExport(export *transform.ExportDef) transform.TopLevel {
	return export
}

func processFuncPass2(fn *transform.FuncDef) transform.TopLevel {
	jsErrorCall := &transform.CallOp{
		Arg: jsNotImplemented,
	}
	tgErrorCall := &transform.CallOp{
		Arg: tinygoNotImplemented,
	}

	switch fn.Name {
	case "$runtime.printitf":
		fn.Code = []transform.Stmt{tgErrorCall}
	case "$_syscall/js.Value_.String",
		"$_syscall/js.Type_String",
		"$_syscall/js.Value_.Get",
		"$syscall/js.handleEvent":
		fn.Code = []transform.Stmt{jsErrorCall}
	default:
		break
	}
	return fn
}

func processFuncPass1(fn *transform.FuncDef) transform.TopLevel {
	switch fn.Name {
	case "$runtime.printitf",
		"$_syscall/js.Value_.String",
		"$_syscall/js.Type_String",
		"$_syscall/js.Value_.Get":
		return fn
	default:
		break
	}
	if strings.HasPrefix(fn.Name, "$_syscall/js.") ||
		strings.HasPrefix(fn.Name, "$syscall/js.") ||
		strings.HasPrefix(fn.Name, "$_struct_syscall/js") ||
		strings.HasPrefix(fn.Name, "$_*struct_syscall/js") ||
		strings.HasPrefix(fn.Name, "$_*syscall/js") {
		return nil
	}
	return fn
}

func convertInputToFormat(filename, tmp, outFile, program, path string) (string, error) {
	var outputName string

	if outFile != "" {
		outputName = outFile
	} else {
		// maybe has more than 1 component
		_, basename := filepath.Split(filename)
		if outFile != "" {
			panic("unable to understand input params to convertInputToFormat")
		}
		outputName = basename
		if len(basename) > 5 && strings.HasSuffix(basename, ".wasm") {
			outputName = basename[0:len(basename)-5] + ".wat"
		} else {
			outputName = outputName + ".wat"
		}
		outputName = filepath.Join(tmp, outputName)
	}
	out, err := os.Create(outputName)
	if err != nil {
		log.Printf("converting input file ("+path+") failed, cannot create temp file: %v", err)
		return "", err
	}
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("converting input file ("+path+") failed, input file does not exist: %v", err)
	}
	cmd := exec.Command(program, filename)
	cmd.Stdout = out
	// stderr file
	errFile := filepath.Join(tmp, "wat2wasm-errors")
	errFp, err := os.Create(errFile)
	if err != nil {
		log.Fatalf("converting input file ("+path+") failed, cannot create temporary error file: %v", err)
	}
	cmd.Stderr = errFp
	err = cmd.Run()
	if err != nil {
		os.Remove(outputName) // so as not to confuse make
		log.Printf("conversion of %s failed, errors of %s are in :%s", path, program, errFile)
		return "", err
	}
	out.Close()
	return outputName, nil
}
