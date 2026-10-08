// test_ir は赤外線を受信している間 GPIO17 の LED を点灯させます (lib_iR.py の main() の移植)。
package main

import (
	"time"

	"github.com/momorara/picoCar/tinygo/lib/ir"
	"github.com/momorara/picoCar/tinygo/lib/led"
)

func main() {
	for {
		sense := ir.Read()
		if sense == 0 {
			led.On()
		} else {
			led.Off()
		}
		println(sense)
		time.Sleep(10 * time.Millisecond)
	}
}
