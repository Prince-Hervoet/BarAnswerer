package util

import (
	"fmt"
	"time"
)

func Info(msg string) {
	fmt.Println("[INFO] " + time.Now().GoString() + ": " + msg)
}
