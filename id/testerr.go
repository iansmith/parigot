package id

const (
	// TestNoError means just what it sounds like.  All Ids that are errors represent
	// no error as 0.
	TestNoError TestErrorCode = 0
	// TestServiceNotFound means that the service that was supposed to be under test
	// could not be found.
	TestErrorServiceNotFound TestErrorCode = 1
	// TestErrorExec means that the exec itself (not the thing being execed) has failed.
	TestErrorExec TestErrorCode = 2
	// TestErrorSendFailed means that the Test code itself could not create the
	// necessary queue entries.
	TestErrorSendFailed TestErrorCode = 3
	// TestErrInternal means that the Test code itself (not the code under test) has had
	//a problem.
	TestErrorInternal TestErrorCode = 4
	// TestErrRegexpFailed means that the regexp provided by the caller did not compile
	// and is not a valid go regexp.
	TestErrorrRegexpFailed TestErrorCode = 5
)
