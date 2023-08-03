package skk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	rl "github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

// ErrJisyoNotFound is an error that means dictionary file not found
var ErrJisyoNotFound = errors.New("Jisyo not found")

// New creats an instance with empty dictionaries.
func New() *Mode {
	return &Mode{
		User:       newJisyo(),
		System:     newJisyo(),
		MiniBuffer: MiniBufferOnNextLine{},
		ctrlJ:      keys.CtrlJ,
	}
}

// Load loads dictionaries and returns new SKK instance.
// A SKK instance is both a container for dictionaries and a command of readline.
func Load(userJisyoFname string, systemJisyoFnames ...string) (*Mode, error) {
	M := New()
	var err error
	if userJisyoFname != "" {
		M.User.Load(userJisyoFname)
		M.userJisyoPath = userJisyoFname
	}
	for _, fn := range systemJisyoFnames {
		err = M.System.Load(fn)
		if err == nil {
			return M, nil
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return nil, ErrJisyoNotFound
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

// Setup sets k in readline's global keymap to boot into SKK mode.
// If you want to set the SKK for a specific readline keymap,
// give the return value of the Load function as the second argument of BindKey
func SetupTo(k keys.Code, userJisyoFname string, systemJisyoFnames ...string) (func() error, error) {
	M, err := Load(userJisyoFname, systemJisyoFnames...)
	if err != nil {
		return func() error { return nil }, err
	}
	M.ctrlJ = k
	rl.GlobalKeyMap.BindKey(k, M)
	return M.SaveUserJisyo, nil
}

// Setup sets Ctrl-J in readline's global keymap to boot into SKK mode.
// If you want to set the SKK for a specific readline keymap,
// give the return value of the Load function as the second argument of BindKey
func Setup(userJisyoFname string, systemJisyoFnames ...string) (func() error, error) {
	return SetupTo(keys.CtrlJ, userJisyoFname, systemJisyoFnames...)
}

// WriteTo outputs the user dictionary to w.
// Please note that the character code is UTF8.
func (M *Mode) WriteTo(w io.Writer) (n int64, err error) {
	return M.User.WriteTo(w)
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
	if _, err := M.User.WriteToEucJp(fd); err != nil {
		return err
	}
	if err := fd.Close(); err != nil {
		return err
	}
	if err := os.Rename(filename, filename+".BAK"); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.Rename(tmpName, filename)
}

func (M *Mode) ConfigWithString(config string) (errs []error) {
	for ok := true; ok; {
		var token string
		token, config, ok = strings.Cut(config, ";")

		key, value, hasEqual := strings.Cut(token, "=")
		var err error
		if hasEqual {
			if strings.EqualFold(key, "user") {
				err = M.User.Load(value)
				if os.IsNotExist(err) {
					err = nil
				}
				if err == nil {
					M.userJisyoPath = value
				}
			} else {
				err = fmt.Errorf("SKK-ERROR: unknown option: %s", key)
			}
		} else {
			err = M.System.Load(token)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (M *Mode) enableUntilExit(ctx context.Context, key keys.Code, B *rl.Buffer) rl.Result {
	rl.GlobalKeyMap.BindKey(key, M)
	rl.GlobalKeyMap.BindKey(keys.Enter, &rl.GoCommand{
		Name: "SKK_ACCEPT_LINE_WITH_LATIN_MODE",
		Func: M.cmdAcceptLineWithLatinMode,
	})
	rl.GlobalKeyMap.BindKey(keys.CtrlC, &rl.GoCommand{
		Name: "SKK_INTRRUPT_WITH_LATIN_MODE",
		Func: M.cmdIntrruptWithLatinMode,
	})
	return M.Call(ctx, B)
}
