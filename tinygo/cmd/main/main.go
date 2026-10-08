// main は picoCar の基本プログラムです (main.py / main_01.py の移植)。
//
// 電源を入れ BOOTSEL スイッチが押されたら、その場で向きを変えて距離を測り、
// 一番遠い方向に向いて前進します。
//
//  1. 距離が 15cm 以内になったら停止
//  2. その場で向きを変えて距離を測り、一番遠い方向に向いて前進
//
// 走行中に赤外線を感知したら停止し、その場で距離の遠い方に回転します。
package main

import (
	"time"

	"github.com/momorara/picoCar/tinygo/lib/bootsw"
	"github.com/momorara/picoCar/tinygo/lib/ir"
	"github.com/momorara/picoCar/tinygo/lib/led"
	"github.com/momorara/picoCar/tinygo/lib/servo"
	"github.com/momorara/picoCar/tinygo/lib/sr04"
)

func main() {
	servo.Stop()
	time.Sleep(500 * time.Millisecond)
	servo.Stop()
	bootsw.Wait() // BOOTSEL スイッチが押されたら実行

	for {
		led.Blink(2)
		// 回転して遠い方向を向く
		println("lib_servo.faraway")
		servo.Faraway(45, 200*time.Millisecond)
		time.Sleep(2 * time.Second)
		// 前方に直進
		println("lib_servo.run")
		servo.Run()

		// 前方の距離が 15cm 未満になるまで前進を続ける
		println("lib_SR04.read")
		for sr04.Read() > 15 {
			// 赤外線を感知したら前進を止める
			println("lib_iR.read")
			if irDetected() {
				break
			}
		}

		println("lib_servo.stop")
		servo.Stop()
	}
}

// irDetected は 1ms 空けて 2 回続けて赤外線を感知したら true を返します。
func irDetected() bool {
	if ir.Read() != 0 {
		return false
	}
	println("赤外線感知")
	time.Sleep(time.Millisecond)
	if ir.Read() != 0 {
		return false
	}
	println("赤外線感知")
	return true
}
