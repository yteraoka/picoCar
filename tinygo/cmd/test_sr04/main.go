// test_sr04 は 0.5 秒ごとに距離を表示します (lib_SR04.py の main() の移植)。
package main

import (
	"time"

	"github.com/momorara/picoCar/tinygo/lib/sr04"
)

func main() {
	for {
		println("距離 =", sr04.Read(), "cm")
		time.Sleep(500 * time.Millisecond)
	}
}
