// Package ir は GPIO16 の赤外線センサーを読み取ります (lib_iR.py の移植)。
//
// リモコンの信号はパルス状の点滅なので、短い信号は感知できない場合があります。
// 長めの信号を出すリモコンボタンを確認してください。
package ir

import "machine"

const pin = machine.GP16

func init() {
	pin.Configure(machine.PinConfig{Mode: machine.PinInput})
}

// Read はセンサーのピンの値を返します。赤外線を感知すると 0 になります。
func Read() int {
	if pin.Get() {
		return 1
	}
	return 0
}
