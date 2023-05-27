package protocol

type InitProtocol struct {
	MagicNumber int8
	Version     int8
	NeedSize    int32
}

func (here *InitProtocol) ToByteArray() {

}
