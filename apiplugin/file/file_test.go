package file

import "testing"

// if you look at the tests in queue_test you'll see that
// the tests ignore the wrapper functions and use the
// real impleentations directly.

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

}
