package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	var prefix string
	flag.StringVar(&prefix, "prefix", "", "prefix of the package to be patched")
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("you must provide a single directory to begin the patching")
	}
	info, err := os.Stat(flag.Arg(0))
	if err != nil {
		log.Fatalf("failed trying to state the parameter '%s': %v", flag.Arg(0), err)
	}
	if !info.IsDir() {
		log.Fatalf("expected '%s' to be a directory", flag.Arg(0))
	}
	if prefix == "" {
		log.Fatalf("you supply the -prefix parameter with the package prefix")
	}
	start := path.Clean(flag.Arg(0))
	dir := filepath.Dir(start)
	entry, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed trying to read the parent dircetory of '%s': %v", flag.Arg(0), err)
	}
	var found fs.DirEntry
	for _, e := range entry {
		if e.Name() == flag.Arg(0) {
			found = e
			break
		}
	}
	if found == nil || !found.IsDir() {
		log.Fatalf("failed trying to find '%s' in the parent dircetory", flag.Arg(9))
	}
	traverse(found, "", prefix)

}
func traverse(entry fs.DirEntry, pathPrefix, pkgPrefix string) {
	path := filepath.Join(pathPrefix, entry.Name())
	if !entry.IsDir() {
		if strings.HasSuffix(entry.Name(), "pb.go") {
			patch(entry.Name(), pathPrefix, pkgPrefix)
		}
		return
	}

	child, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("failed trying to read children of '%s': %v", path, err)
	}
	for _, c := range child {
		traverse(c, path, pkgPrefix)
	}
}

func patch(name, pathPrefix, packagePrefix string) {
	fullpath := filepath.Join(pathPrefix, name)
	fullPkg := packagePrefix + "/" + pathPrefix
	patchLine(fullpath, fullPkg)
}

const badLine = "reflect.TypeOf(x{}).PkgPath()"

func patchLine(path, value string) {
	unchanged := true
	newFile := path + ".new"
	out, err := os.Create(newFile)
	if err != nil {
		log.Fatalf("unable to create '%s' for patching package: %v", path+".new", err)
	}
	fp, err := os.Open(path)
	if err != nil {
		log.Fatalf("unable to open '%s' to patch reflect line: %v", path, err)
	}
	wr := bufio.NewWriter(out)
	rd := bufio.NewScanner(fp)
	for rd.Scan() {
		rawLine := rd.Text()
		line := strings.TrimSpace(rawLine)
		withQ := fmt.Sprintf("\"%s\"", value)
		if line == "GoPackagePath: reflect.TypeOf(x{}).PkgPath()," {
			newline := strings.Replace(line, badLine, withQ, 1)
			//log.Printf("%s-> %s", path, newline)
			wr.WriteString("\t\t\t" + newline + "\n")
			unchanged = false
			continue
		}
		if line == "reflect \"reflect\"" {
			continue // we delete this line, but don't change the unchanged flag
		}
		wr.WriteString(rawLine + "\n")
	}
	wr.Flush()

	fp.Close()
	out.Close()
	if unchanged {
		if err := os.Remove(newFile); err != nil {
			log.Fatalf("unable to remove '%s': %v", newFile, err)
		}
	} else {
		log.Printf("patched: %s", path)
		if err := os.Rename(newFile, path); err != nil {
			log.Fatalf("unable to rename '%s' to '%s': %v", newFile, path, err)
		}
	}
}
