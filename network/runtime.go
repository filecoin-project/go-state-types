package network

import "io"

// These interfaces are intended to match those from whyrusleeping/cbor-gen, such that code generated from that
// system is automatically usable here (but not mandatory).
type CBORMarshaler interface {
	MarshalCBOR(w io.Writer) error
}

type CBORUnmarshaler interface {
	UnmarshalCBOR(r io.Reader) error
}

type CBORer interface {
	CBORMarshaler
	CBORUnmarshaler
}

// Specifies importance of message, LogLevel numbering is consistent with the uber-go/zap package.
type LogLevel int

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DEBUG LogLevel = iota - 1
	// InfoLevel is the default logging priority.
	INFO
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WARN
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ERROR
)
