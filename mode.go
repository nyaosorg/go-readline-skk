package skk

import (
	"context"
	"errors"
	"io"
	"os"

	rl "github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

// ErrJisyoNotFound is an error that means dictionary file not found
var ErrJisyoNotFound = errors.New("Jisyo not found")

// New creats an instance with empty dictionaries.
func New() *Mode {
	return &Mode{
		User:       Jisyo{},
		System:     Jisyo{},
		MiniBuffer: MiniBufferOnNextLine{},
	}
}

// Load loads dictionaries and returns new SKK instance.
// A SKK instance is both a container for dictionaries and a command of readline.
func Load(userJisyoFname string, systemJisyoFnames ...string) (*Mode, error) {
	jisyo := New()
	var err error
	if userJisyoFname != "" {
		jisyo.User.Load(userJisyoFname)
	}
	for _, fn := range systemJisyoFnames {
		err = jisyo.System.Load(fn)
		if err == nil {
			return jisyo, nil
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
	return rl.CONTINUE
}

// Setup sets Ctrl-J in readline's global keymap to boot into SKK mode.
// If you want to set the SKK for a specific readline keymap,
// give the return value of the Load function as the second argument of BindKey
func Setup(userJisyoFname string, systemJisyoFnames ...string) error {
	M, err := Load(userJisyoFname, systemJisyoFnames...)
	if err != nil {
		return err
	}
	rl.GlobalKeyMap.BindKey(keys.CtrlJ, M)
	return nil
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
func (M *Mode) SaveUserJisyo(filename string) error {
	filename = expandEnv(filename)
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
