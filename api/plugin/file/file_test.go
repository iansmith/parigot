package file

import (
	"context"
	"log"
	"testing"

	apishared "github.com/iansmith/parigot/api/shared"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/file/v1"
)

const filePath = apishared.FileServicePathPrefix + "testfile.txt"
const fileContent = "Hello!Parigot!"

var notExistFid = file.NewFileId()
var contentBuf = []byte(fileContent)

// there is a bug
// If I use this test as the last test for a file, it does not create a real file on the disk.
// if I put it as the first test in the file, everything works fine

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

	// Test cases 1: Create files with bad path names
	//  	. in the path name 						- "/parigot/app/./file.go"
	// 		.. in the path name 					- "/parigot/app/../file.go"
	// 		without prefix '/parigot/app/' 			- "dir/file.go"
	// 		with a prefix close to the right one 	- "/parigot/workspace/file.go"
	// 		too much parts in the path name 		- "/parigot/app/1/2/.../19/file.go"
	longPath := "/parigot/app/1/2/3/4/5/6/7/8/9/10/11/12/13/14/15/16/17/18/19/file.go"
	for _, currentName := range []string{"/parigot/app/./file.go", "/parigot/app/../file.go",
		"dir/file.go", "/parigot/workspace/file.go", longPath} {

		badFid := testFileCreate(t, svc, currentName, fileContent,
			"Test cases 1 in CreateClose", true, int32(file.FileErr_InvalidPathError))
		if !badFid.IsZeroValue() {
			t.Errorf("Unexpectedly created a file with the a bad path name")
		}
	}

	// Test case 2: Create a file with a good name
	fid := testFileCreate(t, svc, filePath, fileContent,
		"Test cases 2 in CreateClose", false, int32(file.FileErr_NoError))
	// Test case 3: Create a file that is already in the written status
	testFileCreate(t, svc, filePath, fileContent, "Test cases 3 in CreateClose",
		true, int32(file.FileErr_AlreadyInUseError))
	// Test case 4: Close a good file
	testFileClose(t, svc, fid, "Test cases 4 in CreateClose", false, int32(file.FileErr_NoError))

	// Test case 5: Create a file that is already in the read status
	openAGoodFile(t, svc)
	testFileCreate(t, svc, filePath, fileContent, "Test cases 5 in CreateClose",
		true, int32(file.FileErr_AlreadyInUseError))
	closeAGoodFile(t, svc)

	// Test case 6: Create a file that already exists.
	fid2 := testFileCreate(t, svc, filePath, fileContent, "Test cases 6 in CreateClose",
		false, int32(file.FileErr_NoError))
	if !fid.Equal(fid2) {
		t.Errorf("Unexpected creation of a new file.")
	}

	// Test case 7: Close a file twice, expecting an error on the second attempt.
	testFileClose(t, svc, fid, "Test cases 7 in CreateClose", false, int32(file.FileErr_NoError))
	testFileClose(t, svc, fid, "Test cases 7 in CreateClose", true, int32(file.FileErr_FileClosedError))
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

func TestWrite(t *testing.T) {
	svc := newFileSvc((context.Background()))
	svc.isTesting = true

	// Test case 1: Write a file that does not exist
	testFileWrite(t, svc, notExistFid, contentBuf, "Test cases 1 in Write",
		true, int32(file.FileErr_NotExistError))

	fid := creatAGoodFile(t, svc)
	closeAGoodFile(t, svc)

	// Test case 2: Write a closed file
	testFileWrite(t, svc, fid, contentBuf, "Test cases 2 in Write",
		true, int32(file.FileErr_FileClosedError))

	// Test case 3: Write a read file
	openAGoodFile(t, svc)
	testFileWrite(t, svc, fid, contentBuf, "Test cases 3 in Write",
		true, int32(file.FileErr_AlreadyInUseError))

	// Test case 4: Write a file with 0 length buffer
	fid = creatAGoodFile(t, svc)
	testFileWrite(t, svc, fid, []byte{}, "Test cases 4 in Write",
		false, int32(file.FileErr_NoError))

	// Test case 5: Write a file with a good buffer
	testFileWrite(t, svc, fid, contentBuf, "Test cases 5 in Write", false, int32(file.FileErr_NoError))
}

func TestDelete(t *testing.T) {
	svc := newFileSvc((context.Background()))
	svc.isTesting = true

	// Test case 1: Delete a file that does not exist
	testFileDelete(t, svc, filePath, "Test case 1 in TestDelete", true, int32(file.FileErr_NotExistError))

	// Test case 2: Delete a file that is already in the written status
	creatAGoodFile(t, svc)
	testFileDelete(t, svc, filePath, "Test case 2 in TestDelete", true, int32(file.FileErr_AlreadyInUseError))
	closeAGoodFile(t, svc)

	// Test case 3: Delete a file that is already in the read status
	openAGoodFile(t, svc)
	testFileDelete(t, svc, filePath, "Test case 3 in TestDelete", true, int32(file.FileErr_AlreadyInUseError))

	// Test case 4: Delete a file that is in the closed status
	closeAGoodFile(t, svc)
	testFileDelete(t, svc, filePath, "Test case 4 in TestDelete", false, int32(file.FileErr_NoError))

	// Test case 5: Delete a file that is already deleted
	testFileDelete(t, svc, filePath, "Test case 5 in TestDelete", true, int32(file.FileErr_NotExistError))
}

func TestRealFiles(t *testing.T) {
	svc := newFileSvc((context.Background()))

	// Happy path
	fid := testFileCreate(t, svc, filePath, fileContent, "create a real good file",
		false, int32(file.FileErr_NoError))
	testFileWrite(t, svc, fid, contentBuf, "write a real file", false, int32(file.FileErr_NoError))
	testFileClose(t, svc, fid, "close a file", false, int32(file.FileErr_NoError))

	testFileOpen(t, svc, filePath, "read a file", false, int32(file.FileErr_NoError))
	testFileRead(t, svc, fid, make([]byte, 6), "read a real file", false, int32(file.FileErr_NoError))
	testFileClose(t, svc, fid, "close a file", false, int32(file.FileErr_NoError))
	testFileDelete(t, svc, filePath, "delete a file", false, int32(file.FileErr_NoError))
}

//
// helpers
//

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
		return file.FileIdEmptyValue(), []byte{}
	}
	// If an error was not expected but one occurred.
	if errCode != int32(file.FileErr_NoError) {
		t.Fatalf("Unexpected error occurred while reading a file: %s. Error code: %d", msg, errCode)
	}

	return file.UnmarshalFileId(resp.GetId()), buf
}

func testFileDelete(t *testing.T, svc *fileSvcImpl, fpath string, msg string,
	errExpected bool, expectedErrCode int32) {

	ctx := pcontext.DevNullContext(context.Background())
	t.Helper()

	req := &file.DeleteRequest{
		Path: fpath,
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

func testFileWrite(t *testing.T, svc *fileSvcImpl, fid file.FileId, buf []byte,
	msg string, errExpected bool, expectedErrCode int32) file.FileId {

	ctx := pcontext.DevNullContext(context.Background())
	t.Helper()

	req := &file.WriteRequest{
		Id:  fid.Marshal(),
		Buf: buf,
	}
	resp := &file.WriteResponse{}
	errCode := svc.write(ctx, req, resp)
	if errExpected {
		if errCode != expectedErrCode {
			t.Fatalf("Unexpected error code while writing a file: %s. Expected %d, but got %d",
				msg, expectedErrCode, errCode)
		}
		return file.FileIdEmptyValue()
	}
	// If an error was not expected but one occurred.
	if errCode != int32(file.FileErr_NoError) {
		t.Fatalf("Unexpected error occurred while writing a file: %s. Error code: %d", msg, errCode)
	}

	return file.UnmarshalFileId(resp.GetId())
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
