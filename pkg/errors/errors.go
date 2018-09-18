package errors

import (
	"bytes"
	"fmt"
)

// Error defines a standard application error.
type Error struct {
	// Error classification for the application.
	Kind Kind

	// Human-readable message.
	Message string

	// Logical operation and nested error.
	Op  Op
	Err error
}

// Error returns the string representation of the error message.
func (e *Error) Error() string {
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	if e.Op != "" {
		fmt.Fprintf(&buf, "%s: ", e.Op)
	}

	// If wrapping an error, print its Error() message.
	// Otherwise print the error code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Kind != Other {
			fmt.Fprintf(&buf, "<%s> ", e.Kind)
		}
		buf.WriteString(e.Message)
	}
	return buf.String()
}

type Op string

// Kind defines the kind or class of an error.
type Kind uint8

// Various kinds of errors
const (
	Other        Kind = iota // Unclassified error
	Internal                 // Internal error
	Conflict                 // Conflict when an entity already exists
	Invalid                  // Invalid input, validation error etc
	NotFound                 // Entity does not exist
	Unauthorized             // Unauthorized to perform an action
)

func (k Kind) String() string {
	switch k {
	case Other:
		return "unclassified error"
	case Internal:
		return "internal error"
	case Invalid:
		return "invalid input"
	case NotFound:
		return "entity not found"
	case Unauthorized:
		return "unauthorized"
	default:
		return "unknown error kind"
	}
}

// Is reports whether err is an *Error of the given Kind.
// If err is nil then Is returns false.
func Is(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Other {
		return e.Kind == kind
	}
	if e.Err != nil {
		return Is(kind, e.Err)
	}
	return false
}

// E is a helper function which constructs an `*Error`
// You can pass it Op, Kind, error (Err) or string (Message) in any order and it'll construct it.
func E(args ...interface{}) error {
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Kind:
			e.Kind = arg
		case Op:
			e.Op = arg
		case error:
			e.Err = arg
		case string:
			e.Message = arg
		}
	}
	return e
}
