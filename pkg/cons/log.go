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
