package skk

import (
	"context"
	"os"

	rl "github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

type Config struct {
	UserJisyoPath    string
	SystemJisyoPaths []string
	CtrlJ            keys.Code
	BindTo           canBindKey
	KeepModeOnExit   bool
}

func (c Config) Setup() (skkMode *Mode, err error) {
	skkMode = &Mode{
		User:       newJisyo(),
		System:     newJisyo(),
		MiniBuffer: MiniBufferOnNextLine{},
	}

	if c.CtrlJ != "" {
		skkMode.ctrlJ = c.CtrlJ
	} else {
		skkMode.ctrlJ = keys.CtrlJ
	}
	if c.UserJisyoPath != "" {
		err := skkMode.User.Load(c.UserJisyoPath)
		if err != nil {
			return nil, err
		}
		skkMode.userJisyoPath = c.UserJisyoPath
	}
	for _, fn := range c.SystemJisyoPaths {
		err = skkMode.System.Load(fn)
		if err != nil {
			return nil, err
		}
	}
	if c.BindTo == nil {
		c.BindTo = rl.GlobalKeyMap
	}
	c.BindTo.BindKey(c.CtrlJ, skkMode)
	if !c.KeepModeOnExit {
		c.BindTo.BindKey(keys.Enter, &rl.GoCommand{
			Name: "SKK_ACCEPT_LINE_WITH_LATIN_MODE",
			Func: skkMode.cmdAcceptLineWithLatinMode,
		})
		c.BindTo.BindKey(keys.CtrlC, &rl.GoCommand{
			Name: "SKK_INTRRUPT_WITH_LATIN_MODE",
			Func: skkMode.cmdIntrruptWithLatinMode,
		})
	}
	return skkMode, nil
}

// String returns the name as the command starting SKK
func (M *Mode) String() string {
	return "SKK_MODE"
}

// Call is readline.Command to start SKK henkan mode.
func (M *Mode) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	M.enable(B, hiragana)
	M.message(B, msgHiragana)
	return rl.CONTINUE
}

// SaveUserJisyo saves the user dictionary as filename.
// The file is first created with the name filename+".TMP",
// and replaced with the file of filename after closing.
// The original file is renamed to filename + ".BAK".
func (M *Mode) SaveUserJisyo() error {
	if M.userJisyoPath == "" {
		return nil
	}
	filename := expandEnv(M.userJisyoPath)
	tmpName := filename + ".TMP"
	fd, err := os.Create(tmpName)
	if err != nil {
		return err
	}
	if _, err := M.User.WriteToUtf8(fd); err != nil {
		return err
	}
	if err := fd.Close(); err != nil {
		return err
	}
	backup := filename + ".BAK"
	if err := os.Remove(backup); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Rename(filename, backup); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.Rename(tmpName, filename)
}
