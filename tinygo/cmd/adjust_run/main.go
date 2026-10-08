// adjust_run は前進時の直進性を調整するプログラムです (adjust_run.py の移植)。
//
// BOOTSEL スイッチを押すと 8 秒間走るので、その直進性を見て config.go の
// RightDuty / LeftDuty の Forward (前進) の値を調整します。
// 左にズレた場合は左サーボの回転数が右より少なく、右にズレた場合は多いためです。
//
//	mode 1: 前進
//	mode 2: 後退
//
// 詳しくは 説明資料/前進時の直線性の調整.txt を参照。
package main

import (
	"time"

	"github.com/momorara/picoCar/tinygo/config"
	"github.com/momorara/picoCar/tinygo/lib/bootsw"
	"github.com/momorara/picoCar/tinygo/lib/led"
	"github.com/momorara/picoCar/tinygo/lib/mode"
	"github.com/momorara/picoCar/tinygo/lib/picoled"
	"github.com/momorara/picoCar/tinygo/lib/servo"
)

func main() {
	r, l := config.RightDuty, config.LeftDuty
	println(r.Forward, r.Back, r.Stop)
	println(l.Forward, l.Back, l.Stop)

	// mode の数だけ点滅
	m := mode.Pin()
	led.Blink(m)
	for {
		time.Sleep(500 * time.Millisecond)
		bootsw.Wait() // BOOTSEL スイッチが押されたら実行
		picoled.Blink(2)

		servo.Stop()
		switch m {
		case 1:
			println("前進")
			servo.SetRight(r.Forward)
			servo.SetLeft(l.Forward)
			time.Sleep(8 * time.Second)
		case 2:
			println("後退")
			servo.SetRight(r.Back)
			servo.SetLeft(l.Back)
			time.Sleep(8 * time.Second)
		}
		servo.Stop()

		// mode の数だけ点滅
		m = mode.Pin()
		led.Blink(m)
	}
}
