go-readline-skk
================

このパッケージは Go言語製のコマンドライン向けの一行入力パッケージ [go-readline-ny] に [SKK] ライクな「かな漢字変換機能」を実現するアドオンです。

![./demo.gif](./demo.gif)

辞書などを全てパラメータで指定する場合
--------------------------------------

```example.go
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
```

環境変数などの設定文字列で辞書などを指定する場合
------------------------------------------------

```example2.go
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
    closer := colorable.EnableColorsStdout(nil)
    defer closer()
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
```

[go-readline-ny]: https://github.com/nyaosorg/go-readline-ny
[SKK]: https://ja.wikipedia.org/wiki/SKK
