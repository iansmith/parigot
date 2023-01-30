package lib

const (
	// QueueNoError means just what it sounds like.  All Ids that are errors represent
	// no error as 0.
	QueueNoError QueueErrorCode = 0
	// QueueInvalidName means that the given queue name is a not a valid
	// identifier.  Identifiers must contain only ascii alphanumeric characters
	// and the symbols ".", ",","_" and "-".  The first letter of a queue name
	// must be an alphabetic character.
	QueueInvalidName = 1
	// Queue internal error means that the queue's implementation (not the values)
	// passed to it) is the problem.  This is roughly a 500 not a 401.
	QueueInternalError = 2
	// QueueNoPayload is an error that means that an attempt was made to create
	// a message a nil payload.  Payloads are mandatory and senders are optional.
	QueueNoPayload = 3
)
