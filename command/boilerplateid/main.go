package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
)

//go:embed template/*
var templateFS embed.FS
var singletonPrimary bool
var singletonError bool
var idNoImport bool

func init() {
	flag.BoolVar(&singletonPrimary, "p", false, "indicates to generate only a primary id")
	flag.BoolVar(&singletonError, "e", false, "indicates to generate only an error id")
	flag.BoolVar(&idNoImport, "i", false, "indicates if an import for the id package is needed")
}

func main() {
	flag.Parse()

	if !singletonError && !singletonPrimary {
		generatePrimaryAndError()
		return
	}
	if singletonError && singletonPrimary {
		log.Fatal("the single primary and the single error options are exclusive")
		return
	}
	if singletonError {
		generateErrorOnly()
		return
	}
	if singletonPrimary {
		generatePrimaryOnly()
		return
	}
}
func generateErrorOnly() {
	generateOne("usage: <package name> <errorId name> <errorId Letter> <errorId Shortname>",
		"template/erroronly.tmpl", "errorName", "errorLetter", "errorShortName",
		false, true)
}
func generatePrimaryOnly() {
	generateOne("usage: <package name> <primaryId name> <primaryId letter> <primaryId shortName>",
		"template/primaryonly.tmpl", "primaryName", "primaryLetter", "primaryShortName",
		true, false)
}

func generateOne(usage, path, name, letter, short string, genPrimary, genError bool) {
	if flag.NArg() != 4 {
		log.Fatalf(usage)
	}
	if len(stripLeadingBackslash(flag.Arg(2))) != 1 {
		log.Fatalf("supplied primary letter '%s' is more than one character", flag.Arg(3))
	}
	dot := make(map[string]interface{})
	dot["pkg"] = stripLeadingBackslash(flag.Arg(0))
	dot[name] = stripLeadingBackslash(flag.Arg(1))
	dot[letter] = stringToByteHex(stripLeadingBackslash(flag.Arg(2)))
	dot[short] = stripLeadingBackslash(flag.Arg(3))
	executeTemplate(path, dot, genPrimary, genError)

}
func generatePrimaryAndError() {
	if flag.NArg() != 7 {
		log.Fatalf("usage: <package name> <primary id name> <error id name> <primaryLetter> <primaryShortname> <errorLetter> <errorShortname>")
	}
	if len(stripLeadingBackslash(flag.Arg(3))) != 1 {
		log.Fatalf("supplied primary letter '%s' is more than one character", flag.Arg(3))
	}
	if len(stripLeadingBackslash(flag.Arg(5))) != 1 {
		log.Fatalf("supplied error letter '%s' is more than one character", flag.Arg(5))
	}

	dot := make(map[string]interface{})
	dot["pkg"] = stripLeadingBackslash(flag.Arg(0))
	dot["primaryName"] = stripLeadingBackslash(flag.Arg(1))
	dot["errorName"] = stripLeadingBackslash(flag.Arg(2))
	dot["primaryLetter"] = stringToByteHex(stripLeadingBackslash(flag.Arg(3)))
	dot["primaryShortName"] = stripLeadingBackslash(flag.Arg(4))
	dot["errorLetter"] = stringToByteHex(stripLeadingBackslash(flag.Arg(5)))
	dot["errorShortName"] = stripLeadingBackslash(flag.Arg(6))
	executeTemplate("template/idanderr.tmpl", dot, true, true)
}

func executeTemplate(name string, dot map[string]interface{}, genPrimary, genError bool) {
	dot["idNoImport"] = idNoImport

	if idNoImport {
		dot["idPkg"] = ""
		if genPrimary {
			dot["primaryCast"] = fmt.Sprintf("IdRoot[def%s](f)", dot["primaryName"].(string))
		}
		if genError {
			dot["errorCast"] = fmt.Sprintf("IdRoot[def%s](f)", dot["errorName"].(string))
		}
	} else {
		dot["idPkg"] = "id."
		if genPrimary {
			dot["primaryCast"] = fmt.Sprintf("id.IdRoot[def%s](f)", dot["primaryName"].(string))
		}
		if genError {
			dot["errorCast"] = fmt.Sprintf("id.IdRoot[def%s](f)", dot["errorName"].(string))
		}
	}
	templ, err := template.ParseFS(templateFS, name)
	if err != nil {
		log.Fatal(err)
	}
	if err := templ.Execute(os.Stdout, dot); err != nil {
		log.Fatal(err)
	}

}
func stripLeadingBackslash(s string) string {
	if s[0:1] == "\\" {
		return s[1:]
	}
	return s
}

// assumes s has already been checked to be 1 byte long
func stringToByteHex(s string) string {
	b := []byte(s)[0]
	return fmt.Sprintf("0x%02x", b)
}
