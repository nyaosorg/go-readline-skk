//go:build run

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-skk"
)

func mains() error {
	var ed readline.Editor

	// Windows でエスケープシーケンスを有効にする
	//closer := colorable.EnableColorsStdout(nil)
	//defer closer()
	ed.Writer = colorable.NewColorableStdout()

	// for example:
	//   set "GOREADLINESKK=~/Share/Etc/SKK-JISYO.*;user=~/.go-skk-jisyo"
	//   rem ~/ はパッケージ側で展開されます
	if env := os.Getenv("GOREADLINESKK"); env != "" {
		skkMode, err := skk.Config{
			BindTo: &ed,
		}.SetupWithString(env)

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		} else {
			ed.Coloring = &skk.Coloring{}
			defer func() {
				err := skkMode.SaveUserJisyo()
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				}
			}()
		}
	}
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
