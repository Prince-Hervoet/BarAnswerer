package protocol

const (
	// init magic number
	MAGIC_NUMBER = 101
	TYPE_REQUEST = 1
	TYPE_ACK     = 1 << 1
	TYPE_INIT    = 1 << 2
	TYPE_INFORM  = 1 << 3
	TYPE_QUIT    = 1 << 4
)
