// Package mode はモードピン (GPIO18〜21) を読み込みます (lib_mode.py の移植)。
//
// ターミナルブロック (プラグ) が差さっているピンは 0 になります。
package mode

import "machine"

var pins = [4]machine.Pin{
	machine.GP21, // mode 1
	machine.GP20, // mode 2
	machine.GP19, // mode 3
	machine.GP18, // mode 4
}

func init() {
	for _, p := range pins {
		p.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}
}

// Pin はプラグの差さっている mode 番号 (1〜4) を返します。プラグがなければ 0 です。
func Pin() int {
	for i, p := range pins {
		if !p.Get() {
			return i + 1
		}
	}
	return 0
}
