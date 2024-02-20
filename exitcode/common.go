package exitcode

// Common error codes that may be shared by different actors.
// Actors may also define their own codes, including redefining these values.

const (
	// ErrIllegalArgument indicates that a method parameter is invalid.
	ErrIllegalArgument = FirstActorErrorCode + iota
	// ErrNotFound indicates that a requested resource does not exist.
	ErrNotFound
	// ErrForbidden indicates that an action is disallowed.
	ErrForbidden
	// ErrInsufficientFunds indicates that a balance of funds is insufficient.
	ErrInsufficientFunds
	// ErrIllegalState indicates that an actor's internal state is invalid.
	ErrIllegalState
	// ErrSerialization indicates a de/serialization failure within actor code.
	ErrSerialization
	// ErrUnhandledMessage indicates that the actor cannot handle this message.
	ErrUnhandledMessage
	// ErrUnspecified indicates that the actor failed with an unspecified error.
	ErrUnspecified
	// ErrAssertionFailed indicates that the actor failed a user-level assertion
	ErrAssertionFailed
	// ErrReadOnly indicates that the actor cannot perform the requested operation
	// in read-only mode.
	ErrReadOnly
	/// ErrNotPayable indicates the method cannot handle a transfer of value.
	ErrNotPayable

	// Common error codes stop here.  If you define a common error code above
	// this value it will have conflicting interpretations
	FirstActorSpecificExitCode = ExitCode(32)
)
