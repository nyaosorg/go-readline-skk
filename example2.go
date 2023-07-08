//go:build run
// +build run

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hymkor/go-readline-skk"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

func mains() error {
	skkMode, err := skk.Load("~/.skk-jisyo-nyagos", "SKK-JISYO.L")
	if err != nil {
		return err
	}
	skkMode.System.Load("SKK-JISYO.emoji")

	var ed readline.Editor
	ed.BindKey(keys.CtrlJ, skkMode)
	text, err := ed.ReadLine(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("TEXT:", text)

	skkMode.SaveUserJisyo("~/.skk-jisyo-nyagos")
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
}
