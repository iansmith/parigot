package file

import (
	"context"
	"log"
	"strings"
	"testing"

	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/file/v1"
)

// if you look at the tests in queue_test you'll see that
// the tests ignore the wrapper functions and use the
// real impleentations directly.

const filePath = "/parigot/app/file.go"
const fileContent = "Hello!Parigot!"

var notExistFid = file.NewFileId()

func TestOpenClose(t *testing.T) {
	svc := newFileSvc((context.Background()))
	creatAGoodFile(t, svc)

	// try opening a file does not exist
	badFid := testFileOpen(t, svc, "/parigot/app/badfile.txt", "path doesn't exist",
		true, int32(file.FileErr_NotExistError))
	if !badFid.IsEmptyValue() {
		t.Errorf("trying to open a file which does not exist")
	}
	// try open a file with a badly formed path name
	badFid = testFileOpen(t, svc, "badfile.txt", "bad path name",
		true, int32(file.FileErr_InvalidPathError))
	if !badFid.IsEmptyValue() {
		t.Errorf("accidentally opened a file with the a bad path name")
	}
	// try close a non-existing file
	testFileClose(t, svc, notExistFid, "close a non-exisiting file", true, int32(file.FileErr_NotExistError))

	// try opening and closing a file
	fid := testFileOpen(t, svc, filePath, "open a file", false, int32(file.FileErr_NoError))
	testFileClose(t, svc, fid, "close a file", false, int32(file.FileErr_NoError))

	// also try opening a file twice, for now should be an error
	fid = creatAGoodFile(t, svc)
	openAGoodFile(t, svc)
	testFileOpen(t, svc, filePath, "open a open file", true, int32(file.FileErr_AlreadyInUseError))
	// also try closing a file twice, there should be an error in the second time
	testFileClose(t, svc, fid, "close a file", false, int32(file.FileErr_NoError))
	testFileClose(t, svc, fid, "close a file that does not exist", true, int32(file.FileErr_NotExistError))
}

func TestCreateClose(t *testing.T) {
	svc := newFileSvc((context.Background()))

	// create a file with . in the path name
	badFid := testFileCreate(t, svc, "/parigot/app/./file.go", fileContent,
		"bad path name with .", true, int32(file.FileErr_InvalidPathError))
	if !badFid.IsZeroValue() {
		t.Errorf("accidentally created a file with the a bad path name")
	}
	// create a file with .. in the path name
	badFid = testFileCreate(t, svc, "/parigot/app/../file.go", fileContent,
		"bad path name with ..", true, int32(file.FileErr_InvalidPathError))
	if !badFid.IsZeroValue() {
		t.Errorf("accidentally created a file with the a bad path name")
	}
	// create a file without prefix '/parigot/app/'
	badFid = testFileCreate(t, svc, "dir/file.go", fileContent, "bad path name without right prefix",
		true, int32(file.FileErr_InvalidPathError))
	if !badFid.IsZeroValue() {
		t.Errorf("accidentally created a file with the a bad path name")
	}
	// create a file with a prefix close to the right one
	badFid = testFileCreate(t, svc, "/parigot/workspace/file.go", fileContent, "bad path name with a prefix close to the right one",
		true, int32(file.FileErr_InvalidPathError))
	if !badFid.IsZeroValue() {
		t.Errorf("accidentally created a file with the a bad path name")
	}

	// create a file with a good name
	fid := testFileCreate(t, svc, filePath, fileContent, "good path name", false, int32(file.FileErr_NoError))
	// create a file already exist
	fid2 := testFileCreate(t, svc, filePath, fileContent, "good path name", false, int32(file.FileErr_NoError))
	if !fid.Equal(fid2) {
		t.Errorf("unexpected that a new file was created")
	}

	// close a file twice, the seconde time there should have an error
	testFileClose(t, svc, fid, "close a file", false, int32(file.FileErr_NoError))
	testFileClose(t, svc, fid, "close a closed file", true, int32(file.FileErr_NotExistError))

	// create a file with the same path
	fid2 = creatAGoodFile(t, svc)
	if fid.Equal(fid2) {
		t.Errorf("unexpected the created file has the same ID as the deleted file")
	}
}

func TestRead(t *testing.T) {
	svc := newFileSvc((context.Background()))

	// Read a file that does not exist
	testFileRead(t, svc, notExistFid, make([]byte, 2), "read a non_existent file",
		true, int32(file.FileErr_NotExistError))

	fid := creatAGoodFile(t, svc)

	// read a closed file
	testFileRead(t, svc, fid, make([]byte, 2), "read a closed file",
		true, int32(file.FileErr_FileClosedError))

	openAGoodFile(t, svc)
	// read a file content "Hello!Parigot!" twice.
	// the 1st reading should be "Hello!"
	// 2nd should be "Parigot!"
	_, readBuf := testFileRead(t, svc, fid, make([]byte, 6), "read a file",
		false, int32(file.FileErr_NoError))
	if string(readBuf) != fileContent[:6] {
		t.Errorf("unexpected that read was not as expected")
	}
	_, readBuf = testFileRead(t, svc, fid, make([]byte, 8), "read a file",
		false, int32(file.FileErr_NoError))
	if string(readBuf) != fileContent[6:] {
		t.Errorf("unexpected that read was not as expected")
	}
}

func testFileCreate(t *testing.T, svc *fileSvcImpl, fpath string, content string, msg string,
	errExpected bool, expectedErrCode int32) file.FileId {

	ctx := pcontext.DevNullContext(context.Background())
	t.Helper()

	req := &file.CreateRequest{
		Path:    fpath,
		Content: content,
	}
	resp := &file.CreateResponse{}
	errCode := svc.create(ctx, req, resp)
	if errExpected {
		if errCode == int32(file.FileErr_NoError) {
			log.Fatalf("expected error in creating a file: %s :%d", msg, errCode)
		}
		if expectedErrCode != errCode {
			log.Fatalf("wrong error code in creating a file: %s, expected %d but got %d",
				msg, expectedErrCode, errCode)
		}
		return file.FileIdZeroValue()
	}

	// no error expected case
	if errCode != int32(file.FileErr_NoError) {
		log.Fatalf("unexpected error in creating a file: %s :%d", msg, errCode)
	}

	return file.UnmarshalFileId(resp.GetId())
}

func testFileClose(t *testing.T, svc *fileSvcImpl, fid file.FileId, msg string,
	errExpected bool, expectedErrCode int32) {

	ctx := pcontext.DevNullContext(context.Background())
	t.Helper()

	req := &file.CloseRequest{
		Id: fid.Marshal(),
	}
	resp := &file.CloseResponse{}

	errCode := svc.close(ctx, req, resp)
	if errExpected {
		if errCode == int32(file.FileErr_NoError) {
			log.Fatalf("expected error in closing a file (%s): %s: %d", fid, msg, errCode)
		}
		if errCode != expectedErrCode {
			log.Fatalf("wrong error code in closing a file (%s): %s expected %d but got %d",
				fid, msg, expectedErrCode, errCode)
		}
		return
	}

	// no error expected case
	if errCode != int32(file.FileErr_NoError) {
		log.Fatalf("unexpected error in closing a file (%s): %s :%d", fid, msg, errCode)
	}
	candidate := file.UnmarshalFileId(resp.GetId())
	if !fid.Equal(candidate) {
		log.Fatalf("created and closed file id don't match")
	}
}

func testFileOpen(t *testing.T, svc *fileSvcImpl, fpath string, msg string,
	errExpected bool, expectedErrCode int32) file.FileId {

	ctx := pcontext.DevNullContext(context.Background())
	t.Helper()

	req := &file.OpenRequest{
		Path: fpath,
	}
	resp := &file.OpenResponse{}
	errCode := svc.open(ctx, req, resp)
	if errExpected {
		if errCode == int32(file.FileErr_NoError) {
			log.Fatalf("expected error in opening a file: %s :%d", msg, errCode)
		}
		if errCode != expectedErrCode {
			log.Fatalf("wrong error code in opening a file: %s, expected %d but got %d",
				msg, expectedErrCode, errCode)
		}
		return file.FileIdEmptyValue()
	}

	// no error expected case
	if errCode != int32(file.FileErr_NoError) {
		log.Fatalf("unexpected error in opening a file: %s :%d", msg, errCode)
	}

	return file.UnmarshalFileId(resp.GetId())
}

func testFileRead(t *testing.T, svc *fileSvcImpl, fid file.FileId, buf []byte,
	msg string, errExpected bool, expectedErrCode int32) (file.FileId, []byte) {

	ctx := pcontext.DevNullContext(context.Background())
	t.Helper()

	req := &file.ReadRequest{
		Id:  fid.Marshal(),
		Buf: buf,
	}
	resp := &file.ReadResponse{}
	errCode := svc.read(ctx, req, resp)
	if errExpected {
		if errCode == int32(file.FileErr_NoError) {
			log.Fatalf("expected error in reading a file: %s :%d", msg, errCode)
		}
		if errCode != expectedErrCode {
			log.Fatalf("wrong error code in reading a file: %s, expected %d but got %d",
				msg, expectedErrCode, errCode)
		}
		return file.FileIdEmptyValue(), make([]byte, 0)
	}
	// no error expected case
	if errCode != int32(file.FileErr_NoError) {
		log.Fatalf("unexpected error in reading a file: %s :%d", msg, errCode)
	}

	return file.UnmarshalFileId(resp.GetId()), buf
}

func creatAGoodFile(t *testing.T, svc *fileSvcImpl) file.FileId {
	currentTime := pcontext.CurrentTime(svc.ctx)
	fid := file.NewFileId()

	newFileInfo := fileInfo{
		id:             fid,
		path:           filePath,
		length:         len(fileContent),
		content:        fileContent,
		status:         Fs_Close,
		createDate:     currentTime,
		lastAccessTime: currentTime,
	}

	fileDataCache := *svc.fileDataCache
	fileDataCache[fid] = &newFileInfo

	fpathTofid := *svc.fpathTofid
	fpathTofid[filePath] = fid

	return fid
}

func openAGoodFile(t *testing.T, svc *fileSvcImpl) {
	fileDataCache := *svc.fileDataCache
	fpathTofid := *svc.fpathTofid
	fid := fpathTofid[filePath]
	myFileInfo := fileDataCache[fid]

	myFileInfo.lastAccessTime = pcontext.CurrentTime(svc.ctx)
	myFileInfo.status = Fs_Open
	myFileInfo.reader = strings.NewReader(myFileInfo.content)
}
