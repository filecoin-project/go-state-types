package exitcode

var names = map[ExitCode]string{
	Ok: "Ok",

	// System errors
	SysErrSenderInvalid:      "SysErrSenderInvalid",
	SysErrSenderStateInvalid: "SysErrSenderStateInvalid",
	SysErrIllegalInstruction: "SysErrIllegalInstruction",
	SysErrInvalidReceiver:    "SysErrInvalidReceiver",
	SysErrInsufficientFunds:  "SysErrInsufficientFunds",
	SysErrOutOfGas:           "SysErrOutOfGas",
	SysErrIllegalExitCode:    "SysErrIllegalExitCode",
	SysErrFatal:              "SysFatal",
	SysErrMissingReturn:      "SysErrMissingReturn",
	SysErrReserved1:          "SysErrReserved1",
	SysErrReserved2:          "SysErrReserved2",
	SysErrReserved3:          "SysErrReserved3",
	SysErrReserved4:          "SysErrReserved4",
	SysErrReserved5:          "SysErrReserved5",
	SysErrReserved6:          "SysErrReserved6",

	// Common errors
	ErrIllegalArgument:   "ErrIllegalArgument",
	ErrNotFound:          "ErrNotFound",
	ErrForbidden:         "ErrForbidden",
	ErrInsufficientFunds: "ErrInsufficientFunds",
	ErrIllegalState:      "ErrIllegalState",
	ErrSerialization:     "ErrSerialization",
	ErrUnhandledMessage:  "ErrUnhandledMessage",
	ErrUnspecified:       "ErrUnspecified",
	ErrAssertionFailed:   "ErrAssertionFailed",
	ErrReadOnly:          "ErrReadOnly",
	ErrNotPayable:        "ErrNotPayable",
}
