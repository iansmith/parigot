package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Service struct {
	Name         string
	ProtoPackage string
	Package      string
	ModuleFile   string
	Module       string
	TomlFile     string
}

type Config struct {
	Parigot map[string]Service
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 || flag.NArg() > 1 {
		log.Fatal("you need to pass in one directory to search recursively for .p.toml files")
	}
	result := searchModules(flag.Arg(0))
	for pkg, svc := range result {
		log.Printf("%s:%+v", pkg, svc)
	}
	generateLocator(flag.Arg(0), result)
}

func searchModules(path string) map[string]Service {
	s, err := os.Stat(path)
	if err != nil && os.IsExist(err) {
		log.Fatal("command line argument %s does not exist", path)
	}
	if err != nil {
		log.Fatal("%s: %v", path, err)
	}

	if !s.IsDir() {
		log.Fatal("command line argument %s is not a director", path)
	}

	modFiles := searchAndCollect(path, func(s string) bool { return s == "go.mod" })
	localModNameToFullName := make(map[string]string)
	for _, f := range modFiles {
		mod := parseModFileForModuleLine(f)
		localModNameToFullName[mod] = f
	}
	modToToml := make(map[string][]string)
	for m, n := range localModNameToFullName {
		dir, _ := filepath.Split(n)
		if dir == "" {
			cwd, _ := os.Getwd()
			dir = cwd
		}
		modToToml[m] = searchAndCollect(dir, func(s string) bool { return strings.HasSuffix(s, ".p.toml") })
	}
	services := make(map[string]Service)
	for m, tomlFiles := range modToToml {
		for _, tomlF := range tomlFiles {
			last := strings.LastIndex(tomlF, string(filepath.Separator))
			if last < 0 {
				log.Fatalf("cannot understand path of '%s' (%v)", tomlF, last)
			}
			parseTomlFile(tomlF, m, localModNameToFullName, services)
		}
	}
	return services
}

func searchAndCollect(current string, pred func(fn string) bool) []string {
	collectedFiles := []string{}
	dirent, err := os.ReadDir(current)
	if err != nil {
		log.Fatalf("failed reading directory [%v] %s:%v", current == "", current, err)
	}
	for _, ent := range dirent {
		if !ent.IsDir() && !pred(ent.Name()) {
			continue
		}
		if !ent.IsDir() {
			collectedFiles = append(collectedFiles, filepath.Join(current, ent.Name()))
			continue
		}
		childFiles := searchAndCollect(filepath.Join(current, ent.Name()), pred)
		collectedFiles = append(collectedFiles, childFiles...)
	}
	return collectedFiles
}
func parseModFileForModuleLine(filename string) string {
	fp, err := os.Open(filename)
	if err != nil {
		log.Fatalf("opening %s:%v", filename, err)
	}
	b, err := io.ReadAll(fp)
	if err != nil {
		log.Fatalf("reading %s:%v", filename, err)
	}
	s := string(b)
	scanner := bufio.NewScanner(strings.NewReader(s))
	moduleName := ""
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		l := scanner.Text()
		if strings.Contains(l, "module") {
			smaller := strings.Replace(l, "module", "", 1)
			smallest := strings.TrimSpace(smaller)
			if len(smallest) == 0 {
				log.Fatalf("unable to understand line in module file '%s'", l)
			}
			moduleName = smallest
			break
		}
	}
	if moduleName == "" {
		log.Fatal("Unable to parse module file %s", filename)
	}
	return moduleName
}

func parseTomlFile(tomlPath, pkg string, modToName map[string]string, services map[string]Service) {
	fp, err := os.Open(tomlPath)
	if err != nil {
		log.Fatalf("unable to open %s:%v", tomlPath, err)
	}
	defer fp.Close()
	b, err := io.ReadAll(fp)
	if err != nil {
		log.Fatalf("unable to read %s:%v", tomlPath, err)
	}
	dec := toml.NewDecoder(strings.NewReader(string(b)))
	var conf Config
	_, err = dec.Decode(&conf)
	if err != nil {
		log.Fatalf("unable to understand toml file %s:%v", tomlPath, err)
	}
	//addServicesToModule(modDir, conf.Service)
	for _, svc := range conf.Parigot {
		betterSvc := Service{
			Name:         svc.Name,
			ProtoPackage: svc.ProtoPackage,
			ModuleFile:   modToName[pkg],
			Module:       pkg,
			TomlFile:     tomlPath,
		}
		services[pkg] = betterSvc
	}
}

func generateLocator(searchdir string, services map[string]Service) {
	for pkg, svc := range services {
		fp, err := os.Create(filepath.Join(filepath.Dir(svc.TomlFile), "locator.go"))
		if err != nil {
			log.Fatalf("unable to open locator.go in %s: %v", pkg, err)
		}
		defer fp.Close()
		var buf bytes.Buffer
		buf.WriteString("package " + packageShortName(pkg) + "\n")
		locatorService(pkg, svc, &buf)

		_, err = io.Copy(fp, &buf)
		if err != nil {
			log.Printf("unable to copy locator service text: %v", err)
		}
	}
}

func locatorService(pkg string, svc Service, buf *bytes.Buffer) {
	prefix := pkg
	connect := pkg + "/" + packageShortName(pkg) + "connect"
	magic := "/proto/gen/" //xxxfixme
	importC := fmt.Sprintf("import \"%s%s%s\"", prefix, magic, connect)
	buf.WriteString(importC + "\n")
	buf.WriteString("type " + svc.Name + " vvvconnect" + "." + svc.Name + "Client\n")

	buf.WriteString("func Locate" + svc.Name + "() ")
	buf.WriteString(svc.Name + "{")
	buf.WriteString(
		"}")

}

func packageShortName(pkg string) string {
	last := strings.LastIndex(pkg, string(os.PathSeparator))
	if last == -1 {
		log.Fatalf("unable to understand the package name assocated with %s", pkg)
	}
	derivedPackageName := pkg[last+1:]
	return derivedPackageName
}
