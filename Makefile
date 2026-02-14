ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
    WHICH=where.exe
    DEL=del
    NUL=nul
else
    SET=export
    WHICH=which
    DEL=rm
    NUL=/dev/null
endif

ifndef GO
    SUPPORTGO=go1.20.14
    GO:=$(shell $(WHICH) $(SUPPORTGO) 2>$(NUL) || echo go)
endif

all:
	$(GO) fmt ./...
	$(GO) build

setup: SKK-JISYO.L SKK-JISYO.emoji

SKK-JISYO.L :
	curl -O https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.L

SKK-JISYO.emoji :
	curl -O https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.emoji
