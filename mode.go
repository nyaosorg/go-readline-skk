package skk

import (
	"context"
	"fmt"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
	"os"
)

type Config struct {
	UserJisyoPath    string
	SystemJisyoPaths []string
	CtrlJ            keys.Code
	BindTo           CanBindKey
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
		var err error
		skkMode.userJisyoStamp, err = skkMode.User.load(c.UserJisyoPath)
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
		c.BindTo = readline.GlobalKeyMap
	}
	c.BindTo.BindKey(skkMode.ctrlJ, skkMode)
	if !c.KeepModeOnExit {
		skkMode.setupQuitWithLatinMode(c.BindTo)
	}
	return skkMode, nil
}

func (skkMode *Mode) setupQuitWithLatinMode(X CanBindKey) {
	X.BindKey(keys.Enter, &readline.GoCommand{
		Name: "SKK_ACCEPT_LINE_WITH_LATIN_MODE",
		Func: skkMode.cmdAcceptLineWithLatinMode,
	})
	X.BindKey(keys.CtrlC, &readline.GoCommand{
		Name: "SKK_INTRRUPT_WITH_LATIN_MODE",
		Func: skkMode.cmdIntrruptWithLatinMode,
	})
}

// String returns the name as the command starting SKK
func (M *Mode) String() string {
	return "SKK_MODE"
}

// Call is readline.Command to start SKK henkan mode.
func (M *Mode) Call(ctx context.Context, B *readline.Buffer) readline.Result {
	M.enable(B, hiragana)
	M.displayMode(B, msgHiragana)
	return readline.CONTINUE
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

	stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return M.User.saveAs(filename)
	}
	tmpName := filename + ".TMP"

	if err == nil && stat.ModTime() != M.userJisyoStamp {
		// merge
		other := newJisyo()
		if err = other.Load(filename); err != nil {
			return fmt.Errorf("fail to merge: %w", err)
		}
		for _, h := range M.User.ariHistory {
			if h.val == nil {
				delete(other.ari, h.key)
			} else {
				other.ari[h.key] = h.val
			}
		}
		for _, h := range M.User.nasiHistory {
			if h.val == nil {
				delete(other.nasi, h.key)
			} else {
				other.nasi[h.key] = h.val
			}
		}
		M.User = other
	}
	if err := M.User.saveAs(tmpName); err != nil {
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
