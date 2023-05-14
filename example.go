//go:build run
// +build run

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hymkor/go-readline-skk"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

func mains() error {
	customJisyo := ".skk-jisyo"
	if home, err := os.UserHomeDir(); err == nil {
		customJisyo = filepath.Join(home, customJisyo)
	}
	if err := skk.Setup(customJisyo, "SKK-JISYO.L"); err != nil {
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
