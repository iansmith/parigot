package main

import "github.com/iansmith/parigot/g/queue/v1"

const (
	// QueueInvalidName means that the given queue name is a not a valid
	// identifier.  Identifiers must contain only ascii alphanumeric characters
	// and the symbols ".", ",","_" and "-".  The first letter of a queue name
	// must be an alphabetic character.
	QueueInvalidName queue.QueueErrIdCode = iota + queue.QueueErrIdGuestStart
	// Queue internal error means that the queue's implementation (not the values)
	// passed to it) is the problem.  This is roughly a 500 not a 401.
	QueueInternalError
	// QueueNoPayload is an error that means that an attempt was made to create
	// a message a nil payload.  Payloads are mandatory and senders are optional.
	QueueNoPayload
	// QueueNotFound means that the Queue name requested could not be found.
	// This the queue equivalent of 404.
	QueueNotFound
)
