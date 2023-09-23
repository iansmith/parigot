package file

import (
	"context"
	"log"
	"os"
	"path/filepath"

	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/g/file/v1"
)

const (
	filePath    = apishared.FileServicePathPrefix + "testfile.txt"
	fileContent = "Hello!Parigot!"
)

func creatAGoodFile(svc *fileSvcImpl, fpath string, fcontent string) file.FileId {
	// return svc.createANewFile(fpath, fcontent)
	fid, err := svc.createANewFile(context.Background(), fpath, fcontent)
	if err != nil {
		log.Fatal("Failed to create a file: ", err)
	}
	return fid
}

func openAGoodFile(ctx context.Context, svc *fileSvcImpl) {
	fid := (*svc.fpathTofid)[filePath]
	myFileInfo := (*svc.fileDataCache)[fid]

	_, myFileInfo.modTime = currentTimeHost(ctx)
	myFileInfo.status = Fs_Read

	var err error
	myFileInfo.rdClose, err = openHookForStrings(myFileInfo.content)
	if err != nil {
		log.Fatal("Failed to open a file: ", err)
	}
}

func closeAGoodFile(svc *fileSvcImpl) {
	fid := (*svc.fpathTofid)[filePath]
	(*svc.fileDataCache)[fid].status = Fs_Close
}

func createDirOnHost(dirPath string) {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		log.Fatal("Failed to create directory: ", err)
	}
}

// Creates 3 test files in the specified directory
// path/bar.txt, path/t1/unreadable.txt, path/t1/t2/foo.txt
// path/t1/unreadable.txt is unreadable
func createTestFilesOnHost(path string, content string) {

	f1 := filepath.Join(path, "bar.txt")
	f2 := filepath.Join(path, "t1/unreadable.txt")
	f3 := filepath.Join(path, "t1/t2/foo.txt")

	filelogger.Info("createTestFilesOnHost", "path", path, "content", content)
	for _, f := range []string{f1, f2, f3} {
		createDirOnHost(filepath.Dir(f))

		// Create the file
		f, err := os.Create(f)
		if err != nil {
			log.Fatal("Failed to create file: ", err)
		}
		defer f.Close()

		_, err = f.WriteString(content)
		if err != nil {
			log.Fatal("Failed to write to file: ", err)
		}

	}

	// Make unreadable.txt unreadable
	fm := os.FileMode(0222)
	err := os.Chmod(f2, fm)
	if err != nil {
		log.Fatal("Failed to make file unreadable: ", err)
	}
}

func delTestDirOnHost(path string) error { return os.RemoveAll(path) }
