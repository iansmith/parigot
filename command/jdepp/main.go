package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var rootParigotPkg = flag.String("rpkg", "github.com/iansmith/parigot", "name of the package that is the root of parigot")
var generatedParigotPkg = flag.String("ppkg", "github.com/iansmith/parigot/g", "name of the package that indicates a parigot generated file")
var parigotPath = flag.String("ppath", ".", "path to the root directory of parigot source code")
var myPkg = flag.String("m", "", "name of an extra package (usually your code) to look for that indicates generated code")
var buildDir = flag.String("b", "build", "location binaries are built to, relative to the starting dir")
var templateDir = flag.String("t", "command/protoc-gen-parigot/template", "location of the parigot templates")
var lang = flag.String("l", "go", "comma separated list of languages to check for")
var help = flag.Bool("h", false, "this help message")
var makefile = flag.String("k", "Makefile", "location of the Makefile relative to the start of the search")
var verbose = flag.Bool("v", true, "print out lots of information about what jdepp is finding")

const magicLine = "#### Do not remove this line or edit below it.  The rest of this file is computed by jdepp."

type binDep struct {
	name     string
	parigot  []string
	user     []string
	internal []string
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 || *help == true {
		fmt.Printf("[-ppkg parigot generated files package] [-ppath root of parigot source] [-m extra package containing generated code] [path to start search]\n")
		flag.Usage()
		os.Exit(1)
	}
	rootPackage := inputArgToRootPackage(flag.Arg(0))
	log.Printf("%45s:%s", "root package", rootPackage)

	moduleMap := make(map[string]string)
	ignoredModuleMap := make(map[string]string)
	populateModuleMaps(flag.Arg(0), rootPackage, moduleMap, ignoredModuleMap)
	for k, v := range moduleMap {
		log.Printf("%45s:%s", "module map["+k+"]", v)
	}
	for k, v := range moduleMap {
		log.Printf("%45s:%s", "ignored module map["+k+"]", v)
	}
	log.Printf("%45s:%s", "root parigot package", *rootParigotPkg)
	if rootPackage != *rootParigotPkg {
		dummy := make(map[string]string)
		log.Printf("%45s", "repopulating module maps")
		populateModuleMaps(*parigotPath, *rootParigotPkg, moduleMap, dummy)
	}

	programEntryFile := nonIgnoredMainGo(flag.Arg(0), ignoredModuleMap)

	binaryMap := make(map[string]string)
	populateBinaryMap(programEntryFile, binaryMap)

	goSource := make(map[string][]string)
	populateGoSources(binaryMap, goSource)

	//
	// this is the brains of the operation, it figures out the dependencies or panics when it cannot
	//
	allDeps := make(map[string]map[string]struct{})
	for binary, source := range goSource {
		depSet := make(map[string]struct{})
		allDeps[filepath.Join(flag.Arg(0), *buildDir, binary)] = depSet
		for _, src := range source {
			//src is the main.go but implicitly you are dependent on everything in THAT package as well
			dir, _ := filepath.Split(src)
			file, err := findGoFileInDir(dir)
			if err != nil {
				log.Fatalf("%v", err)
			}
			for _, f := range file {
				depSet[f] = struct{}{}
			}
			allImport, generatedParigot, generatedUser, err := findImportInSource(src, *generatedParigotPkg, *myPkg)
			if err != nil {
				log.Fatalf("%v", err)
			}
			myImports := []string{}
			for _, import_ := range allImport {
				if strings.HasPrefix(import_, rootPackage) {
					myImports = append(myImports, import_)
				}
			}
			if len(myImports) > 0 || len(generatedParigot) > 0 || len(generatedUser) > 0 {
				for _, i := range myImports {
					path, dir := findPathForImport(i, moduleMap, rootPackage)
					file, err := findGoFileInDir(filepath.Join(path, dir))
					if err != nil {
						log.Fatalf("%v", err)
					}
					for _, f := range file {
						depSet[f] = struct{}{}
					}
				}
				for _, i := range generatedParigot {
					//templ, err := parigotTemplatesForImport(i)
					//if err != nil {
					//	log.Fatalf("v", err)
					//}
					//for _, t := range templ {
					//	depSet[t] = struct{}{}
					//}
					keyPart := strings.TrimPrefix(i, *generatedParigotPkg)
					for path, mod := range moduleMap {
						path := filepath.Clean(path)
						if strings.HasSuffix(mod, keyPart) {
							// ugh, special cases
							if !strings.HasSuffix(filepath.Clean(path), "/go") {
								panic(fmt.Sprintf("can't understand parigot layout, expected module in %s to end in /go", path))
							}
							dir := strings.TrimSuffix(path, "/go")
							_, err := os.Stat(filepath.Join(dir, "proto"))
							if err != nil && errors.Is(err, os.ErrNotExist) {
								panic(fmt.Sprintf("can't figure out parigot layout, expected module in %s to have sibling /proto", path))
							}
							if err != nil {
								log.Fatalf("%v", err)
							}
							file, err := findNamedFile(filepath.Join(dir, "proto"), ".proto", true)
							if err != nil {
								log.Fatalf("%v", err)
							}
							for _, f := range file {
								depSet[f] = struct{}{}
							}
						}
					}

					depSet[filepath.Join(*parigotPath, "build/protoc-gen-parigot")] = struct{}{}
				}
				for _, i := range generatedUser {
					path, dir := findPathForImport(i, moduleMap, rootPackage)
					if !strings.HasPrefix(dir, "proto/g") {
						panic(fmt.Sprintf("sorry, not smart enough to understand your file layout,expected %s as a prefix on %s",
							"proto/g", dir))
					}
					directoryWithProtos := filepath.Join(path, "proto")
					info, err := os.Stat(directoryWithProtos)
					if errors.Is(err, os.ErrNotExist) {
						panic(fmt.Sprintf("sorry, not smart enough to understand your file layout,expected there to be a protos directory inside %s",
							path))
					}
					if !info.IsDir() {
						panic(fmt.Sprintf("sorry, not smart enough to understand your file layout,expected protos to be a directory inside %s (but its a file)",
							path))
					}
					file, err := findNamedFile(directoryWithProtos, ".proto", true)
					for _, f := range file {
						depSet[f] = struct{}{}
					}
				}
			}
		}
	}
	// recompute the bottom of the makefile
	newDeps := depsToBuffer(allDeps)
	mkpath := filepath.Join(flag.Arg(0), *makefile)
	rewriteMakefile(mkpath, newDeps)
	os.Exit(0)
}

func depsToBuffer(allDeps map[string]map[string]struct{}) *bytes.Buffer {
	var result bytes.Buffer
	for bin, deps := range allDeps {
		if len(deps) == 0 {
			continue
		}
		result.WriteString(fmt.Sprintf("### jdepp computed dependencies for binary: %s\n", bin))
		result.WriteString(fmt.Sprintf("%s: \\\n", bin))
		ct := 0
		for f := range deps {
			result.WriteString(fmt.Sprintf("\t%s", f))
			if ct != len(deps)-1 {
				result.WriteString(fmt.Sprintf(" \\"))
			}
			ct++
			result.WriteString(fmt.Sprintf("\n"))
		}
		result.WriteString("\n")
	}
	return &result
}

func rewriteMakefile(mkpath string, newDeps *bytes.Buffer) {
	fp, err := os.Open(mkpath)
	if err != nil {
		log.Fatalf("looking for Makefile: %v", err)
	}
	var mBuffer bytes.Buffer
	mScanner := bufio.NewScanner(fp)
	found := false
	for mScanner.Scan() {
		line := mScanner.Text()
		mBuffer.WriteString(line + "\n")
		if line == magicLine {
			found = true
			break
		}
	}
	if !found {
		log.Printf("unable to find magic line in makefile (%s)", mkpath)
		log.Fatalf("looking for:'%s'", magicLine)
	}
	fp.Close()
	err = os.Rename(mkpath, mkpath+".orig")
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("jdepp: renamed %s to %s", mkpath, mkpath+".orig")
	fp, err = os.Create(mkpath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	countSaved, err := io.Copy(fp, &mBuffer)
	if err != nil {
		log.Fatalf("copying previous buffer: %v", err)
	}
	countNew, err := io.Copy(fp, newDeps)
	if err != nil {
		log.Fatalf("copying new deps: %v", err)
	}
	log.Printf("created %s with a size of %d bytes", mkpath, countSaved+countNew)
}

func parigotTemplatesForImport(path string) ([]string, error) {
	pkg := *generatedParigotPkg
	candidate := filepath.Join(pkg, "parigot/")
	path = strings.TrimPrefix(path, candidate)
	languages := strings.Split(*lang, ",")
	result := []string{}
	for _, l := range languages {
		t := *templateDir
		templDir := filepath.Join(t, l)
		dir, err := os.Stat(templDir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				panic(fmt.Sprintf("unable to find template directory for language %s, should be %s", l, templDir))
			}
			return nil, fmt.Errorf("parigotTemplatesForImport: %v", err)
		}
		if !dir.IsDir() {
			panic(fmt.Sprintf("for language %s, %s should be a directory, but it is not", l, templDir))
		}
		ent, err := os.ReadDir(templDir)
		if err != nil {
			return nil, fmt.Errorf("parigotTemplatesForImport: %v", err)
		}
		for _, f := range ent {
			result = append(result, filepath.Join(templDir, f.Name()))
		}
	}
	return result, nil
}

func populateBinaryMap(entryPoint []string, binaryMap map[string]string) {
	for _, j := range entryPoint {
		parent, _ := filepath.Split(j)
		if parent == "" {
			parent, _ = os.Getwd()
		}
		_, name := filepath.Split(parent)
		binaryMap[name] = parent
	}
}
func goModToPackage(filename string, target string) (string, string, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return "", "", fmt.Errorf("goModToPackage: %v", err)
	}
	buf, err := io.ReadAll(fp)
	if err != nil {
		return "", "", err
	}
	modBuffer := bufio.NewScanner(bytes.NewBuffer(buf))
	for modBuffer.Scan() {
		line := modBuffer.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module") {
			line = strings.TrimPrefix(line, "module")
			line = strings.TrimSpace(line)
			if target != "" && !strings.HasPrefix(line, target) {
				log.Printf("\tjdepp:go.mod for module %s is not part of package %s", filename, target)
				return "", filename, nil
			}
			return strings.TrimSpace(line), "", nil
		}
	}
	if modBuffer.Err() != nil {
		return "", "", modBuffer.Err()
	}
	panic(fmt.Sprintf("unable to find module declaration in %s", filename))
}
func findGoModFile(start string) ([]string, error) {
	return findNamedFile(start, "go.mod", false)
}
func findMainGoFile(start string) ([]string, error) {
	return findNamedFile(start, "main.go", false)
}
func findNamedFile(start string, target string, suffix bool) ([]string, error) {
	if strings.HasSuffix(start, ".git") || strings.Index(start, "tmp/parse") != -1 { // can be large
		return []string{}, nil
	}
	result := []string{}
	content, err := os.ReadDir(start)
	if err != nil {
		return nil, err
	}
	for _, c := range content {
		if c.IsDir() {
			additional, err := findNamedFile(filepath.Join(start, c.Name()), target, suffix)
			if err != nil {
				return nil, err
			}
			result = append(result, additional...)
		}
		if !suffix {
			if c.Name() == target {
				result = append(result, filepath.Join(start, c.Name()))
			}
		} else {
			if strings.HasSuffix(c.Name(), target) {
				result = append(result, filepath.Join(start, c.Name()))
			}
		}
	}
	return result, nil
}

func matchDirUpward(file string, moduleMap map[string]string) (string, error) {
	for k, v := range moduleMap {
		if strings.HasPrefix(file, k) {
			log.Printf("\tjdepp: ignoring binary based on %s because of %s", file, v)
			return k, nil
		}
	}
	return "", nil
}

func inputArgToRootPackage(arg string) string {
	source := filepath.Join(arg, "go.mod")
	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("unable to determine root go package because directory given as first argument (%s) has no go.mod",
			source)
	}
	rootPackage, ignored, err := goModToPackage(source, "")
	if err != nil {
		log.Fatalf("%v", err)
	}
	if ignored != "" {
		panic("unable to comprehend how this happened ")
	}
	return rootPackage
}
func populateModuleMaps(arg, rootPackage string, moduleMap map[string]string, ignoredModuleMap map[string]string) {
	goDotMod, err := findGoModFile(arg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, f := range goDotMod {
		pkg, ignore, err := goModToPackage(f, rootPackage)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if ignore != "" {
			parent, _ := filepath.Split(f)
			if parent == "" {
				parent = "."
			}
			ignoredModuleMap[parent] = ignore
			continue
		}
		if pkg == "" {
			continue
		}
		parent, _ := filepath.Split(f)
		if parent == "" {
			parent = "."
		}
		moduleMap[parent] = pkg
	}
}
func nonIgnoredMainGo(arg string, ignoredModuleMap map[string]string) []string {
	entryFile, err := findMainGoFile(arg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	notIgnored := []string{}
	for _, m := range entryFile {
		parentGoMod, err := matchDirUpward(m, ignoredModuleMap)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if parentGoMod != "" {
			continue
		}
		notIgnored = append(notIgnored, m)
	}
	return notIgnored
}

func findGoSource(rootDir string) []string {
	result, err := findNamedFile(rootDir, ".go", true)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return result
}

func populateGoSources(binaryMap map[string]string, goSources map[string][]string) {
	for k, v := range binaryMap {
		filtered := []string{}
		result := findGoSource(v)
	outer:
		for _, candidate := range result {
			parts := strings.Split(candidate, "/")
			for _, p := range parts {
				if p == "g" {
					log.Printf("\tjdepp:ignoring generated file %s", candidate)
					continue outer
				}
			}
			filtered = append(filtered, candidate)
		}
		goSources[k] = filtered
	}
}

func findImportInSource(file string, arg string, myPkg string) ([]string, []string, []string, error) {
	fp, err := os.Open(file)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("findImportInSource: %v", err)
	}
	all, err := io.ReadAll(fp)
	if err != nil {
		return nil, nil, nil, err
	}
	sourceScanner := bufio.NewScanner(bytes.NewBuffer(all))
	inside := false
	result := []string{}
	generated := []string{}
	user := []string{}
	for sourceScanner.Scan() {
		line := sourceScanner.Text()
		line = strings.TrimSpace(line)
		if inside {
			if line == "" {
				continue
			}
			if line == ")" {
				inside = false
				continue
			}
			imp := splitOffPossibleImportRename(line, line, file)
			switch {
			case strings.HasPrefix(imp, arg):
				generated = append(generated, imp)
			case myPkg != "" && strings.HasPrefix(imp, myPkg):
				user = append(user, imp)
			default:
				result = append(result, imp)
			}
			continue
		}
		// normal case, we are not inside import(...) block
		if strings.HasPrefix(line, "import ") { //note space
			// two cases here
			short := strings.TrimSpace(strings.TrimPrefix(line, "import ")) //note space
			if short == "(" {
				inside = true
				continue
			}
			imp := splitOffPossibleImportRename(short, line, file)
			switch {
			case strings.HasPrefix(imp, arg):
				generated = append(generated, imp)
			case myPkg != "" && strings.HasPrefix(imp, myPkg):
				user = append(user, imp)
			default:
				result = append(result, imp)
			}
		}
	}
	if sourceScanner.Err() != nil {
		return nil, nil, nil, sourceScanner.Err()
	}
	return result, generated, user, nil
}

func splitOffPossibleImportRename(short string, line string, file string) string {
	if len(short) < 2 {
		panic(fmt.Sprintf("short is too short in %s at line %s", file, line))
	}
	parts := strings.Split(short, " ")
	if len(parts) > 2 {
		panic(fmt.Sprintf("unable to understand import line [2] '%s' in %s", line, file))
	}
	target := parts[0]
	if len(parts) == 2 {
		target = parts[1]
	}
	if target[0:1] != "\"" || target[len(target)-1:] != "\"" {
		panic(fmt.Sprintf("unable to understand import line [3] '%s' in %s", line, file))
	}
	return target[1 : len(target)-1]
}

func findPathForImport(pkg string, moduleMap map[string]string, rootPkg string) (string, string) {
	for path, candidate := range moduleMap {
		path := filepath.Clean(path)
		candidate = filepath.Clean(candidate)
		if strings.HasPrefix(pkg, candidate) && candidate != rootPkg {
			result := filepath.Clean(strings.TrimPrefix(pkg, candidate))
			return path, result
		}
	}
	// now consider root package
	found := false
	pathToRootPkg := ""
	for key, candidate := range moduleMap {
		if candidate == rootPkg {
			found = true
			pathToRootPkg = key
			break
		}
	}
	if !found {
		panic(fmt.Sprintf("unable to find root package %s in module map", rootPkg))
	}
	if strings.HasPrefix(pkg, rootPkg) {
		rootPkg = filepath.Clean(rootPkg)
		remaining := strings.TrimPrefix(pkg, rootPkg)
		if remaining[0:1] == "/" {
			remaining = remaining[1:]
		}
		return pathToRootPkg, remaining
	}
	panic(fmt.Sprintf("unable to find the package needed to understand import of %s", pkg))
}

func findGoFileInDir(dir string) ([]string, error) {
	content, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("findGoFileInDir:%v", err)
	}
	result := []string{}
	for _, f := range content {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".go") {
			result = append(result, filepath.Join(dir, f.Name()))
		}
	}
	return result, nil
}
