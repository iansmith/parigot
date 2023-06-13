package main

import (
	"context"
	"log"
	"testing"

	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/file/v1"
)

// if you look at the tests in queue_test you'll see that
// the tests ignore the wrapper functions and use the
// real impleentations directly.

const fileContent = "Hello! Parigot!"

func TestOpenClose(t *testing.T) {
	// try open a file with a badly formed path name, make sure it fails
	// see filepath.Clean()

	// try opening and closing a file

	// also try opening a file twice, for now should be an error
	// also try closing a file twice

	// when expecting an error, for now you'll have to use
	// the raw numbers from the list of errors
	// apiwasm/file/fileerr.go
}

func TestCreateClose(t *testing.T) {
	svc := newFileSvc((context.Background()))

	// create a file with duplicate "/" in the path name
	testFileCreate(t, svc, "/parigot/app///file.go", fileContent, "bad path name", true, int32(file.FileErr_InvalidPathError))
	// create a file with . in the path name
	testFileCreate(t, svc, "/parigot/app/./file.go", fileContent, "bad path name", true, int32(file.FileErr_InvalidPathError))
	// create a file with .. in the path name
	testFileCreate(t, svc, "/parigot/app/../file.go", fileContent, "bad path name", true, int32(file.FileErr_InvalidPathError))
	// create a file without prefix '/parigot/app/'
	testFileCreate(t, svc, "file.go", fileContent, "bad path name", true, int32(file.FileErr_InvalidPathError))

	// create a file with a good name
	fid := creatAGoodFile(t, svc)
	// create a file already exist
	fid2 := creatAGoodFile(t, svc)
	if !fid.Equal(fid2) {
		t.Errorf("unexpected that the file was not appended")
	}

	// close a file twice, the seconde time there should have an error
	testFileClose(t, svc, fid, "close a file", false, int32(file.FileErr_NoError))
	testFileClose(t, svc, fid, "close a file", true, int32(file.FileErr_NotExistError))

	// create a file with the same path
	fid2 = creatAGoodFile(t, svc)
	if fid.Equal(fid2) {
		t.Errorf("unexpected that second creation of a deleted file gives same id")
	}
}

func testFileCreate(t *testing.T, svc *fileSvcImpl, path string, content string, msg string,
	expectedErr bool, expectedErrCode int32) file.FileId {

	ctx := pcontext.DevNullContext(context.Background())
	t.Helper()

	openReq := &file.CreateRequest{
		Path:    path,
		Content: content,
	}
	openResp := &file.CreateResponse{}
	errCode := svc.create(ctx, openReq, openResp)
	if expectedErr {
		if errCode == int32(file.FileErr_NoError) {
			log.Fatalf("expected error from creating a file: %s :%d", msg, errCode)
		}
		if expectedErrCode != errCode {
			log.Fatalf("wrong error code from creating a file: %s, expected %d but got %d",
				msg, expectedErrCode, errCode)
		}
		return file.FileIdZeroValue()
	}

	// no error expected case
	if errCode != int32(file.FileErr_NoError) {
		log.Fatalf("unexpected error from creating a file: %s :%d", msg, errCode)
	}

	return file.UnmarshalFileId(openResp.GetId())
}

func testFileClose(t *testing.T, svc *fileSvcImpl, fid file.FileId, msg string,
	expectedErr bool, expectedErrCode int32) {

	ctx := pcontext.DevNullContext(context.Background())
	t.Helper()

	req := &file.CloseRequest{}
	resp := &file.CloseResponse{}
	req.Id = fid.Marshal()
	errCode := svc.close(ctx, req, resp)
	if expectedErr {
		if errCode == int32(file.FileErr_NoError) {
			log.Fatalf("expected error from closing a file (%s): %s: %d", fid, msg, errCode)
		}
		if errCode != expectedErrCode {
			log.Fatalf("wrong error code from closing a file (%s): %s expected %d but got %d",
				fid, msg, expectedErrCode, errCode)
		}
		return
	}

	// no error expected case
	if errCode != int32(file.FileErr_NoError) {
		log.Fatalf("unexpected error from closing a file (%s): %s :%d", fid, msg, errCode)
	}
	candidate := file.UnmarshalFileId(resp.GetId())
	if !fid.Equal(candidate) {
		log.Fatalf("created and closed file id don't match")
	}
	return
}

func creatAGoodFile(t *testing.T, svc *fileSvcImpl) file.FileId {
	return testFileCreate(t, svc, "/parigot/app/file.go", fileContent, "good path name",
		false, int32(file.FileErr_NoError))
}
