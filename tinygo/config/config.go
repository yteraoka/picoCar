// Package config は picoCar の設定値です (config.py の移植)。
//
// サーボの個体差に合わせてここの値を調整し、再ビルドして書き込んでください。
// デューティ値は MicroPython の duty_u16() と同じ 0〜65535 のスケールです。
// 調整方法は 説明資料/前進時の直線性の調整.txt と 説明資料/左右回転の調整.txt を参照。
package config

import "time"

// Duty は車輪ごとの 前進, 後退, 停止 のデューティ値です。
type Duty struct {
	Forward uint16
	Back    uint16
	Stop    uint16
}

// RightDuty は右車輪の 前進, 後退, 停止 のデューティ値です。
var RightDuty = Duty{Forward: 4420, Back: 5250, Stop: 4900}

// LeftDuty は左車輪の 前進, 後退, 停止 のデューティ値です。
var LeftDuty = Duty{Forward: 5312, Back: 4450, Stop: 4900}

// AdjustmentRotateTime は回転調整用の時間です。mode 1〜4 に対応します。
var AdjustmentRotateTime = [4]time.Duration{
	1300 * time.Millisecond,
	1400 * time.Millisecond,
	1500 * time.Millisecond,
	1700 * time.Millisecond,
}

// RightRotateTime は右回転で 180 度回転する時間です。
var RightRotateTime = 1600 * time.Millisecond

// LeftRotateTime は左回転で 180 度回転する時間です。
var LeftRotateTime = 1650 * time.Millisecond
