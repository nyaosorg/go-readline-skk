//go:build run
// +build run

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hymkor/go-readline-skk"
	"github.com/nyaosorg/go-readline-ny"
)

func main() {
	skk.Setup("SKK-JISYO.L")

	var ed readline.Editor
	text, err := ed.ReadLine(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
	fmt.Println("TEXT:", text)
}
