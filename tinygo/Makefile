# picoCar TinyGo 版
#
#   make                       全プログラムをビルド (build/*.uf2)
#   make flash-wsl PROG=main   WSL から Windows 経由で書き込み (usbipd 不要)
#   make flash PROG=main       tinygo flash で書き込み (Linux / usbipd で接続した WSL)
#   make monitor               シリアル出力を表示 (Linux / usbipd で接続した WSL)

TARGET ?= pico2-w
PROG   ?= main
TINYGO ?= tinygo

PROGS := $(notdir $(wildcard cmd/*))

.PHONY: all build flash flash-wsl monitor clean

all: $(addprefix build/,$(addsuffix .uf2,$(PROGS)))

build: build/$(PROG).uf2

build/%.uf2: FORCE
	@mkdir -p build
	$(TINYGO) build -target=$(TARGET) -o $@ ./cmd/$*

flash:
	$(TINYGO) flash -target=$(TARGET) ./cmd/$(PROG)

flash-wsl: build/$(PROG).uf2
	scripts/flash-wsl.sh $<

monitor:
	$(TINYGO) monitor -target=$(TARGET)

clean:
	rm -rf build

.PHONY: FORCE
FORCE:
