// test_picoled は基板上の LED を 5 回点滅させます (lib_LED_pico.py の main() の移植)。
package main

import "github.com/momorara/picoCar/tinygo/lib/picoled"

func main() {
	picoled.Blink(5)
	picoled.End()
}
