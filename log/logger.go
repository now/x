package log

// Logger with a simple interface for logging entries and meta-data.
type Logger interface {
	// Entry consisting of a message and zero or more fields.
	//
	// Errors if an entry can’t be added to the log.
	Entry(string, ...Field) error

	// Named is a new Logger with an additional name component.
	//
	// The receiver is returned if the given string is empty.  Otherwise, the new
	// Logger’s name is the given string, if the receiver’s name is empty, or the
	// receiver’s name, a period (U+0046), and the given string, otherwise.
	//
	// This is used to differentiate loggers in a tree-like fashion and a
	// non-empty name will be included in any following log entries.
	Named(string) Logger

	// With fields for additional context in a new Logger.
	//
	// The receiver is returned if the given slice of fields is empty.  Otherwise,
	// the the receiver’s fields and the given fields will be included with any
	// following log entries added to the new Logger.
	With(...Field) Logger
}
