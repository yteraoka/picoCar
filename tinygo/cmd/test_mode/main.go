// test_mode は mode ピンの値を 1 秒ごとに表示します (lib_mode.py の main() の移植)。
package main

import (
	"time"

	"github.com/momorara/picoCar/tinygo/lib/mode"
)

func main() {
	for {
		println("mode =", mode.Pin())
		time.Sleep(time.Second)
	}
}
