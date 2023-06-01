package util

const (
	// init magic number
	MAGIC_NUMBER = 101
	VERSION      = 1

	// memory
	MEM_FLAG_WORKING = 1
	MEM_FLAG_COMMON  = 0
	MEM_HEADER_SIZE  = 9

	CLIENT      = 0
	SERVER      = 1
	DUPLEX      = 2
	EVENTS_SIZE = 1024

	DEFAULT_FILE_DIR         = "/tmp/share/"
	SHARE_MEMORY_HEADER_SIZE = 17

	RPLY       = 2
	WORKING    = 1
	NO_WORKING = 0

	//message type
	SHACK_HAND_MESSAGE = 0
	WAVE_HAND_MESSAGE  = 1
	NOTICE_MESSAGE     = 2
)
