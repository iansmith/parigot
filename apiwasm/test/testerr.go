package main

import (
	testg "github.com/iansmith/parigot/g/test/v1"
)

const (
	// TestNoError means just what it sounds like.  All Ids that are errors represent
	// no error as 0.
	TestNoError testg.TestErrIdCode = iota + testg.TestErrIdGuestStart
	// TestServiceNotFound means that the service that was supposed to be under test
	// could not be found.
	TestErrorServiceNotFound
	// TestErrorExec means that the exec itself (not the thing being execed) has failed.
	TestErrorExec
	// TestErrorSendFailed means that the Test code itself could not create the
	// necessary queue entries.
	TestErrorSendFailed
	// TestErrInternal means that the Test code itself (not the code under test) has had
	//a problem.
	TestErrorInternal
	// TestErrRegexpFailed means that the regexp provided by the caller did not compile
	// and is not a valid go regexp.
	TestErrorrRegexpFailed
	// TestErrorMarshaling is used to when we cannot marshal arguments
	// into to a protobuf.
	TestErrorMarshaling
	// TestErrorQueue means that the there was internal error with the queue
	// that is used by the Test service.
	TestErrorQueue
)
