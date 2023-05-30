package core

const (
	VERSION           = 1
	MAGIC_NUMBER      = int16(178)
	MESSAGE_TYPE_INIT = 1
	MESSAGE_TYPE_DATA = 2
	MESSAGE_TYPE_QUIT = 3
)

type Address struct {
	Ip   string
	Port int
}

type Message struct {
	Header  *MessageHeader
	Payload any
}

type MessageHeader struct {
	MagicNumber int16
	Version     int8
	MessageType int8
}

type FirstBody struct {
	Header  *MessageHeader
	NeedCap int32
	Key     int64
}

type SecondBody struct {
	Header   *MessageHeader
	FilePath string
	Cap      int32
}
