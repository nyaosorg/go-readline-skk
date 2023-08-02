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
	// ~/ はパッケージ側で展開されます
	key := keys.CtrlJ
	// key := keys.CtrlBackslash
	if err := skk.SetupTo(key, "~/.skk-jisyo-nyagos", "SKK-JISYO.L"); err != nil {
		return err
	}

	var ed readline.Editor
	text, err := ed.ReadLine(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("TEXT:", text)

	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
}
