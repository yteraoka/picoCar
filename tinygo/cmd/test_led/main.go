// test_led は GPIO17 の LED を 5 回点滅させます (lib_LED.py の main() の移植)。
package main

import "github.com/momorara/picoCar/tinygo/lib/led"

func main() {
	led.Blink(5)
	led.End()
}
