package file

import (
	"context"
	"log"
	"testing"

	apishared "github.com/iansmith/parigot/api/shared"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/file/v1"
)

// if you look at the tests in queue_test you'll see that
// the tests ignore the wrapper functions and use the
// real impleentations directly.

const filePath = apishared.FileServicePathPrefix + "testfile.txt"
const fileContent = "Hello!Parigot!"

var notExistFid = file.NewFileId()

func TestOpenClose(t *testing.T) {
	svc := newFileSvc((context.Background()))
	svc.isTesting = true

	creatAGoodFile(t, svc)
	closeAGoodFile(t, svc)

	// Test case: Open a non-existent file.
	badFid := testFileOpen(t, svc, "/parigot/app/badfile.txt", "open a file with non-existent path",
		true, int32(file.FileErr_NotExistError))
	if !badFid.IsEmptyValue() {
		t.Errorf("Attempted to open a non-existent file.")
	}
	// Test case: Open a file with a bad path.
	badFid = testFileOpen(t, svc, "badfile.txt", "open a file with bad path name",
		true, int32(file.FileErr_InvalidPathError))
	if !badFid.IsEmptyValue() {
		t.Errorf("Unexpectedly opened a file with the a bad path name")
	}
	// Test case: Close a non-existent file.
	testFileClose(t, svc, notExistFid, "close a non-existent file", true, int32(file.FileErr_NotExistError))

	// Test case: Open and then close a file.
	fid := testFileOpen(t, svc, filePath, "open a file", false, int32(file.FileErr_NoError))
	testFileClose(t, svc, fid, "close a file", false, int32(file.FileErr_NoError))

	// Create a good file
	fid = creatAGoodFile(t, svc)
	openAGoodFile(t, svc)

	// Test case: Open a file that is already open.
	testFileOpen(t, svc, filePath, "open an already opened file", true, int32(file.FileErr_AlreadyInUseError))
	testFileClose(t, svc, fid, "close a file", false, int32(file.FileErr_NoError))
	// Test case: Close a file that is already close
	testFileClose(t, svc, fid, "close a already closed file", true, int32(file.FileErr_FileClosedError))
}

func TestCreateClose(t *testing.T) {
	svc := newFileSvc((context.Background()))
	svc.isTesting = true

	// Test case: Create a file with . in the path name
	badFid := testFileCreate(t, svc, "/parigot/app/./file.go", fileContent,
		"bad path name with .", true, int32(file.FileErr_InvalidPathError))
	if !badFid.IsZeroValue() {
		t.Errorf("Unexpectedly created a file with the a bad path name")
	}
	// Test case: Create a file with .. in the path name
	badFid = testFileCreate(t, svc, "/parigot/app/../file.go", fileContent,
		"bad path name with ..", true, int32(file.FileErr_InvalidPathError))
	if !badFid.IsZeroValue() {
		t.Errorf("Unexpectedly created a file with the a bad path name")
	}
	// Test case: Create a file without prefix '/parigot/app/'
	badFid = testFileCreate(t, svc, "dir/file.go", fileContent, "bad path name without right prefix",
		true, int32(file.FileErr_InvalidPathError))
	if !badFid.IsZeroValue() {
		t.Errorf("Unexpectedly created a file with the a bad path name")
	}
	// Test case: Create a file with a prefix close to the right one
	badFid = testFileCreate(t, svc, "/parigot/workspace/file.go", fileContent,
		"bad path name with a prefix close to the right one",
		true, int32(file.FileErr_InvalidPathError))
	if !badFid.IsZeroValue() {
		t.Errorf("Unexpectedly created a file with the a bad path name")
	}

	// Test case: Create a file with a good name
	fid := testFileCreate(t, svc, filePath, fileContent, "good path name",
		false, int32(file.FileErr_NoError))
	// Test case: Create a file that is already in use.
	testFileCreate(t, svc, filePath, fileContent, "create a file in use",
		true, int32(file.FileErr_AlreadyInUseError))
	// Test case: Close a good file
	testFileClose(t, svc, fid, "close a file", false, int32(file.FileErr_NoError))

	// Test case: Create a file that already exists.
	fid2 := testFileCreate(t, svc, filePath, fileContent, "good path name",
		false, int32(file.FileErr_NoError))
	if !fid.Equal(fid2) {
		t.Errorf("Unexpected creation of a new file.")
	}

	// Test case: Close a file twice, expecting an error on the second attempt.
	testFileClose(t, svc, fid, "close a file", false, int32(file.FileErr_NoError))
	testFileClose(t, svc, fid, "close a closed file", true, int32(file.FileErr_FileClosedError))
}

func TestRead(t *testing.T) {
	svc := newFileSvc((context.Background()))
	svc.isTesting = true

	// Test case: Read a file that does not exist
	testFileRead(t, svc, notExistFid, make([]byte, 2), "read a non-existent file",
		true, int32(file.FileErr_NotExistError))

	fid := creatAGoodFile(t, svc)
	closeAGoodFile(t, svc)

	// Test case: Read a closed file
	testFileRead(t, svc, fid, make([]byte, 2), "read a closed file",
		true, int32(file.FileErr_FileClosedError))

	openAGoodFile(t, svc)

	// Test case: Read a file with 0 length buffer
	testFileRead(t, svc, fid, make([]byte, 0), "read a file with 0 length buffer",
		false, int32(file.FileErr_NoError))

	// Test case: Read a file content "Hello!Parigot!" twice.
	// 			  First reading should be "Hello!" Second should be "Parigot!"
	_, readBuf := testFileRead(t, svc, fid, make([]byte, 6), "read a file",
		false, int32(file.FileErr_NoError))
	if string(readBuf) != fileContent[:6] {
		t.Errorf("unexpected result: read was not as expected")
	}
	_, readBuf = testFileRead(t, svc, fid, make([]byte, 8), "read a file",
		false, int32(file.FileErr_NoError))
	if string(readBuf) != fileContent[6:] {
		t.Errorf("unexpected result: read was not as expected")
	}

	// Test case: Read a file to the end
	testFileRead(t, svc, fid, make([]byte, 2), "read a file to the end",
		false, int32(file.FileErr_NoError))

	// Test case: Read with a larger buffer than the maximum allowed buffer size
	testFileRead(t, svc, fid, make([]byte, apishared.FileServiceMaxBufSize+1),
		"read with large buffer", true, int32(file.FileErr_LargeBufError))
}

// need to write more decent tests
func TestRealFiles(t *testing.T) {
	svc := newFileSvc((context.Background()))

	// create a file
	fid := testFileCreate(t, svc, filePath, fileContent, "create a real good file",
		false, int32(file.FileErr_NoError))
	testFileClose(t, svc, fid, "close a file", false, int32(file.FileErr_NoError))

	testFileOpen(t, svc, filePath, "read a file", false, int32(file.FileErr_NoError))
	testFileRead(t, svc, fid, make([]byte, 6), "read a real file", false, int32(file.FileErr_NoError))
	testFileClose(t, svc, fid, "close a file", false, int32(file.FileErr_NoError))
	testFileDelete(t, svc, fid, "delete a file", false, int32(file.FileErr_NoError))
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
		if errCode != expectedErrCode {
			t.Fatalf("Unexpected error code while creating a file: %s. Expected %d, but got %d",
				msg, expectedErrCode, errCode)
		}
		return file.FileIdZeroValue()
	}

	// If an error was not expected but one occurred.
	if errCode != int32(file.FileErr_NoError) {
		t.Fatalf("Unexpected error occurred while creating a file: %s. Error code: %d", msg, errCode)
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
		if errCode != expectedErrCode {
			t.Fatalf("Unexpected error code while closing a file: %s. Expected %d, but got %d",
				msg, expectedErrCode, errCode)
		}
		return
	}

	// If an error was not expected but one occurred.
	if errCode != int32(file.FileErr_NoError) {
		t.Fatalf("Unexpected error occurred while closing a file: %s. Error code: %d", msg, errCode)
	}
	// just be careful
	candidate := file.UnmarshalFileId(resp.GetId())
	if !fid.Equal(candidate) {
		log.Fatalf("Created and closed file id don't match")
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
		if errCode != expectedErrCode {
			t.Fatalf("Unexpected error code while openning a file: %s. Expected %d, but got %d",
				msg, expectedErrCode, errCode)
		}
		return file.FileIdEmptyValue()
	}

	// If an error was not expected but one occurred.
	if errCode != int32(file.FileErr_NoError) {
		t.Fatalf("Unexpected error occurred while openning a file: %s. Error code: %d", msg, errCode)
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
		if errCode != expectedErrCode {
			t.Fatalf("Unexpected error code while reading a file: %s. Expected %d, but got %d",
				msg, expectedErrCode, errCode)
		}
		return file.FileIdEmptyValue(), make([]byte, 0)
	}
	// If an error was not expected but one occurred.
	if errCode != int32(file.FileErr_NoError) {
		t.Fatalf("Unexpected error occurred while reading a file: %s. Error code: %d", msg, errCode)
	}

	return file.UnmarshalFileId(resp.GetId()), buf
}

func testFileDelete(t *testing.T, svc *fileSvcImpl, fid file.FileId, msg string,
	errExpected bool, expectedErrCode int32) {

	ctx := pcontext.DevNullContext(context.Background())
	t.Helper()

	req := &file.DeleteRequest{
		Id: fid.Marshal(),
	}
	resp := &file.DeleteResponse{}
	errCode := svc.delete(ctx, req, resp)
	if errExpected {
		if errCode != expectedErrCode {
			t.Fatalf("Unexpected error code while deleting a file: %s. Expected %d, but got %d",
				msg, expectedErrCode, errCode)
		}
		return
	}

	// If an error was not expected but one occurred.
	if errCode != int32(file.FileErr_NoError) {
		t.Fatalf("Unexpected error occurred while deleting a file: %s. Error code: %d", msg, errCode)
	}
}

func creatAGoodFile(t *testing.T, svc *fileSvcImpl) file.FileId {
	return svc.createANewFile(filePath, fileContent)
}

func openAGoodFile(t *testing.T, svc *fileSvcImpl) {
	fid := (*svc.fpathTofid)[filePath]
	myFileInfo := (*svc.fileDataCache)[fid]

	myFileInfo.lastAccessTime = pcontext.CurrentTime(svc.ctx)
	myFileInfo.status = Fs_Read
	myFileInfo.rdClose = openHookForStrings(myFileInfo.content)
}

func closeAGoodFile(t *testing.T, svc *fileSvcImpl) {
	fid := (*svc.fpathTofid)[filePath]
	(*svc.fileDataCache)[fid].status = Fs_Close
}
