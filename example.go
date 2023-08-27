//go:build run
// +build run

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-readline-skk"
)

func mains() error {
	var ed readline.Editor

	// Windows でエスケープシーケンスを有効にする
	closer := colorable.EnableColorsStdout(nil)
	defer closer()
	ed.Writer = colorable.NewColorableStdout()

	// ~/ はパッケージ側で展開されます
	skkMode, err := skk.Config{
		UserJisyoPath:    "~/.go-skk-jisyo",
		SystemJisyoPaths: []string{"SKK-JISYO.L", "SKK-JISYO.emoji"},
		CtrlJ:            keys.CtrlJ,
		KeepModeOnExit:   false,
		BindTo:           &ed,
	}.Setup()

	if err != nil {
		return err
	}
	defer func() {
		err := skkMode.SaveUserJisyo()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}()

	for {
		text, err := ed.ReadLine(context.Background())
		if err != nil {
			return err
		}
		fmt.Println("TEXT:", text)
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
}
