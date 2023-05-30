package core

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
