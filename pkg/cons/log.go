package cons

import (
	"fmt"
	"log"
	"runtime"
)

func Log(str ...any) {
	log.Println(str...)
	if runtime.GOOS == "windows" {
		fmt.Print("\r")
	}

}

func NewLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
