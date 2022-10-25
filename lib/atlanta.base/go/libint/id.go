package libint

// Id is a type representing a global identifier in parigot.  They are composed of
// a character ('x') and a number.  In production, that number is 112 bits of
// randomness or a small integer.  The small integer case in production is for
// error ids, to indicate a call has failed in an "expected" way.
// In development, the numbers are always small integers, so printing them out
// is easier.  In all cases, the character that is the highest order byte indicates
// the type of thing the id represents. Ids of different types with the same number
// are not equal.
type Id interface {
	// Short returns a short string for debugging, like [s-6a29].  The number is
	// the last 2 bytes of the full id number.
	Short() string
	// String returns a long string that uniquely identifies this id.  This is
	// usualy something like [r-xx-xxxxxxxx-xxxxxxxx-xxxx-xxxx] where all the x's
	// are hex digits.  Note that Short() is the equivalent of the first and last
	// five characters of string.  If the number is a small integer, the leading
	// zeros are omitted.
	String() string
	// IsError returns true if this is an error type id and there is an error.  It returns
	// false if this is an error type id and there is no error (0 value).  If
	// this is not an error type id, it panics.
	IsError() bool
	// Type returns the name of the type of id, like "service" or "locate error"
	Type() string
	// Equal returns true if the two ids are of the same type and have the same number.
	Equal(Id) bool
}
