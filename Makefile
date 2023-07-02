all:
	go fmt
	go build

setup: SKK-JISYO.L SKK-JISYO.emoji

SKK-JISYO.L :
	curl -O https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.L

SKK-JISYO.emoji :
	curl -O https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.emoji
