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
	closer := skk.SetupOnDemand(keys.CtrlJ, func() (string, bool) {
		return os.LookupEnv("NYAGOSKK")
	})
	defer closer()

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
