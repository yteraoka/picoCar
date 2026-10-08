// faraway は BOOTSEL スイッチを押すと、その場で向きを変えながら距離を測り、
// 一番遠い方向を向くプログラムです (lib_servo.py の main() の移植)。
//
// mode ピンで 1 ステップの角度が変わります: なし 20度, 1: 30度, 2: 45度, 3: 60度, 4: 90度
package main

import (
	"time"

	"github.com/momorara/picoCar/tinygo/lib/bootsw"
	"github.com/momorara/picoCar/tinygo/lib/mode"
	"github.com/momorara/picoCar/tinygo/lib/picoled"
	"github.com/momorara/picoCar/tinygo/lib/servo"
)

var steps = [5]int{20, 30, 45, 60, 90}

func main() {
	for {
		m := mode.Pin()
		time.Sleep(500 * time.Millisecond)
		bootsw.Wait() // BOOTSEL スイッチが押されたら実行
		picoled.Blink(2)
		servo.Faraway(steps[m], 500*time.Millisecond)
		time.Sleep(5 * time.Second)
	}
}
