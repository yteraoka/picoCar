// Package cds は GPIO26 (ADC0) の Cds セルで明るさを測ります (lib_Cds.py の移植)。
//
// キャリブレーション方法: cmd/test_cds を書き込み、表示される元データの
// 最大の明るさ・最大の暗さの値を Max, Min に設定してください。
package cds

import (
	"fmt"
	"machine"
)

// 明暗の範囲
const (
	Max = 60000
	Min = 1500
)

var adc = machine.ADC{Pin: machine.ADC0}

func init() {
	machine.InitADC()
	adc.Configure(machine.ADCConfig{})
}

// Read は明るさを 0〜1000 で返します。verbose が true なら読み取り値を表示します。
func Read(verbose bool) int {
	org := int(adc.Get())
	v := min(max(org, Min), Max)
	lv := (v - Min) * 1000 / (Max - Min)
	if verbose {
		fmt.Println("Analog1 Value: ", org, v, lv)
	}
	return lv
}
