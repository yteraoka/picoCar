// test_cds は Cds の値を 1 秒ごとに 60 回表示します (lib_Cds.py の main() の移植)。
// 表示される元データを使って lib/cds の Max, Min をキャリブレーションできます。
package main

import (
	"time"

	"github.com/momorara/picoCar/tinygo/lib/cds"
)

func main() {
	first := cds.Read(true)
	for range 60 {
		println(first, cds.Read(true))
		time.Sleep(time.Second)
	}
}
