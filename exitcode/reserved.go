package exitcode

// The system error codes are reserved for use by the runtime.
// No actor may use one explicitly. Correspondingly, no runtime invocation should abort with an exit
// code outside this list.
// We could move these definitions out of this package and into the runtime spec.
const (
	Ok = ExitCode(0)

	// Indicates that the actor identified as the sender of a message is not valid as a message sender:
	// - not present in the state tree
	// - not an account actor (for top-level messages)
	// - code CID is not found or invalid
	// (not found in the state tree, not an account, has no code).
	SysErrSenderInvalid = ExitCode(1)

	// Indicates that the sender of a message is not in a state to send the message:
	// - invocation out of sequence (mismatched CallSeqNum)
	// - insufficient funds to cover execution
	SysErrSenderStateInvalid = ExitCode(2)

	// Indicates the message receiver trapped (panicked).
	SysErrIllegalInstruction = ExitCode(4)

	// Indicates that the receiver of a message is not valid (and cannot be implicitly created).
	SysErrInvalidReceiver = ExitCode(5)

	// Indicates that a message sender has insufficient balance for the value being sent.
	// Note that this is distinct from SysErrSenderStateInvalid when a top-level sender can't
	// cover value transfer + gas. This code is only expected to come from inter-actor sends.
	SysErrInsufficientFunds = ExitCode(6)

	// Indicates that message execution (including subcalls) used more gas than the specified
	// limit.
	SysErrOutOfGas = ExitCode(7)

	// Indicates that the actor attempted to exit with a reserved exit code.
	SysErrIllegalExitCode = ExitCode(9)

	// Indicates that something unexpected happened in the system. This always indicates a bug.
	SysErrFatal = ExitCode(10)

	// Indicates the actor returned a block handle that doesn't exist.
	SysErrMissingReturn = ExitCode(11)

	// Unused
	SysErrReserved1 = ExitCode(3)
	SysErrReserved2 = ExitCode(8)
	SysErrReserved3 = ExitCode(12)
	SysErrReserved4 = ExitCode(13)
	SysErrReserved5 = ExitCode(14)
	SysErrReserved6 = ExitCode(15)

	// Used by the builtin actors, so we keep these around for historical reasons.

	// DEPRECATED
	SysErrInvalidMethod = ExitCode(3)
	// DEPRECATED
	SysErrForbidden = ExitCode(8)
	// DEPRECATED
	SysErrorIllegalActor = ExitCode(9)
	// DEPRECATED
	SysErrorIllegalArgument = ExitCode(10)
)

// The initial range of exit codes is reserved for system errors.
// Actors may define codes starting with this one.
const FirstActorErrorCode = ExitCode(16)
