package core

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	HANDSHAKE_FIRST  = 11
	HANDSHAKE_SECOND = 12
	HANDSHAKE_THIRD  = 13
)

func FirstParse(buffer []byte) (*FirstBody, error) {
	data := make([]byte, 0)
	for i := 0; i < len(buffer); i++ {
		if buffer[i] == '\n' {
			data = append(data, buffer[0:i]...)
			break
		}
	}
	first := &FirstBody{}
	fmt.Println(string(data))
	err := json.Unmarshal(data, first)
	if err != nil {
		return nil, errors.New("json error")
	}
	if first.Header.Version != VERSION {
		return nil, errors.New("version error")
	}
	if first.Header.MagicNumber != MAGIC_NUMBER {
		return nil, errors.New("magicnumber error")
	}
	if first.Header.MessageType != MESSAGE_TYPE_INIT {
		return nil, errors.New("type error")
	}
	if first.NeedCap <= 0 {
		return nil, errors.New("cap error")
	}
	return first, nil
}

func ReadJsonByLineFeed(buffer []byte) (*Message, error) {
	data := make([]byte, 0)
	for i := 0; i < len(buffer); i++ {
		if buffer[i] == '\n' {
			data = append(data, buffer[0:i]...)
			break
		}
	}
	message := &Message{}
	err := json.Unmarshal(data, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func CheckHeader(message *Message) bool {
	if message == nil {
		return false
	}
	if message.Header.MagicNumber != MAGIC_NUMBER {
		return false
	} else if message.Header.Version != VERSION {
		return false
	}
	return true
}
