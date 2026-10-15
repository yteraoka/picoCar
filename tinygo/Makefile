# picoCar TinyGo 版
#
#   make                       全プログラムをビルド (build/*.uf2)
#   make flash-wsl PROG=main   WSL から Windows 経由で書き込み (usbipd 不要)
#   make flash PROG=main       tinygo flash で書き込み (Linux / usbipd で接続した WSL)
#   make monitor               シリアル出力を表示 (Linux / usbipd で接続した WSL)
#
# リモコン版 (cmd/remote) の無線 LAN 設定は wifi.mk (git 管理外) に書くか、
# make の引数で指定します:  make flash-wsl PROG=remote WIFI_SSID=xxxx WIFI_PASSWORD=yyyy

TARGET ?= pico2-w
PROG   ?= main
TINYGO ?= tinygo

PROGS := $(notdir $(wildcard cmd/*))

-include wifi.mk
WIFI_SSID     ?=
WIFI_PASSWORD ?=
WIFI_IP       ?=
LDFLAGS_remote = -X 'main.ssid=$(WIFI_SSID)' -X 'main.password=$(WIFI_PASSWORD)' -X 'main.ipAddr=$(WIFI_IP)'

.PHONY: all build flash flash-wsl monitor clean

all: $(addprefix build/,$(addsuffix .uf2,$(PROGS)))

build: build/$(PROG).uf2

build/%.uf2: FORCE
	@mkdir -p build
	$(TINYGO) build -target=$(TARGET) -ldflags="$(LDFLAGS_$*)" -o $@ ./cmd/$*

flash:
	$(TINYGO) flash -target=$(TARGET) -ldflags="$(LDFLAGS_$(PROG))" ./cmd/$(PROG)

flash-wsl: build/$(PROG).uf2
	scripts/flash-wsl.sh $<

monitor:
	$(TINYGO) monitor -target=$(TARGET)

clean:
	rm -rf build

.PHONY: FORCE
FORCE:
