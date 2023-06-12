package main

import (
	"context"
	"testing"

	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/file/v1"
	filemsg "github.com/iansmith/parigot/g/msg/file/v1"
)

// if you look at the tests in queue_test you'll see that
// the tests ignore the wrapper functions and use the
// real impleentations directly.

func TestOpenClose(t *testing.T) {
	svc := newFileSvc((context.Background()))

	// Needs to be rewritten
	testFileCreate(t, svc, "/app/", "good path name", false, 0)

	// que := filemsg.OpenRequest{
	// 	Path: "$$$$.###",
	// }

	// try open a file with a badly formed path name, make sure it fails
	// see filepath.Clean()

	// try opening and closing a file

	// also try opening a file twice, for now should be an error
	// also try closing a file twice

	// when expecting an error, for now you'll have to use
	// the raw numbers from the list of errors
	// apiwasm/file/fileerr.go

	///////////////////////////////////////
	// new a file service
}

func TestCreateClose(t *testing.T) {

}

func testFileCreate(t *testing.T, svc *fileSvcImpl, path string, msg string, errorExpected bool, expectedCode uint16) file.FileId {
	ctx := pcontext.DevNullContext(context.Background())
	t.Helper()

	create := &filemsg.CreateRequest{
		Path: path,
	}
	resp := &filemsg.CreateResponse{}
	err := file.NewFileErrIdFromRaw(svc.create(ctx, create, resp))
	if errorExpected {
		if !err.IsError() {
			t.Errorf("expected error: %s :%s", msg, err.Short())
		}
		if file.FileErrIdCode(expectedCode) != err.ErrorCode() {
			t.Errorf("wrong code : %s, expected %d but got %d", msg, expectedCode, err.ErrorCode())
		}
		return file.ZeroValueFileId()
	}

	// no error expected case
	if err.IsError() {
		t.Errorf("unexpected error: %s :%s", msg, err.Short())
		return file.ZeroValueFileId()
	}

	return file.MustUnmarshalFileId(resp.GetId())
}
