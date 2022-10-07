package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Service struct {
	Name    string
	Package string
}
type Config struct {
	Parigot map[string]Service
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("you need to pass in one or more directories to searchModules for .p.toml files")
	}
	for i := 0; i < flag.NArg(); i++ {
		searchModules(flag.Arg(i))
	}
}

func searchModules(path string) {
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
	modToName := make(map[string]string)
	for _, f := range modFiles {
		mod := parseModFileForModuleLine(f)
		modToName[mod] = f
	}
	modToToml := make(map[string][]string)
	for m, n := range modToName {
		dir, _ := filepath.Split(n)
		modToToml[m] = searchAndCollect(dir, func(s string) bool { return strings.HasSuffix(s, ".p.toml") })
	}
	for m, tomlFiles := range modToToml {
		for _, tomlF := range tomlFiles {
			parts := strings.Split(tomlF, string(filepath.Separator))
			if len(parts) < 2 {
				log.Fatalf("cannot understand split path of '%s' (%+v)", tomlF, parts)
			}
			tomlPath := filepath.Join(filepath.Dir(modToName[m]), filepath.Join(parts[1:]...))
			parseTomlFile(tomlPath, m)
		}
	}
}

func searchAndCollect(current string, pred func(fn string) bool) []string {
	collectedFiles := []string{}
	dirent, err := os.ReadDir(current)
	if err != nil {
		log.Fatalf("failed reading directory %s:%v", current, err)
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

func parseTomlFile(tomlPath, modDir string) {
	log.Printf("parseTomlFile: %s,%s", tomlPath, modDir)
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
	log.Printf("%s->%+v", tomlPath, conf.Parigot)
}
