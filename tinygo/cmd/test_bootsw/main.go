// test_bootsw は BOOTSEL スイッチが押されるまで基板上 LED を点滅させます
// (lib_bootSW.py の main() の移植)。
package main

import (
	"time"

	"github.com/momorara/picoCar/tinygo/lib/bootsw"
)

func main() {
	for {
		bootsw.Wait()
		println("BOOTSEL が押されました")
		time.Sleep(500 * time.Millisecond)
	}
}
