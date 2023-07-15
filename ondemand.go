package skk

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

type onDemandLoad struct {
	Key    keys.Code
	Func   func() (string, bool)
	closer func() error
}

func (o *onDemandLoad) Close() error {
	if o.closer != nil {
		return o.closer()
	}
	return nil
}

func (o *onDemandLoad) String() string {
	return "SKK_MODE_ONDEMAND_SETUP"
}

func loadWithConfigString(config string) (skkMode *Mode, errs []error) {
	skkMode = New()
	for ok := true; ok; {
		var token string
		token, config, ok = strings.Cut(config, ";")

		key, value, hasEqual := strings.Cut(token, "=")
		var err error
		if hasEqual {
			if strings.EqualFold(key, "user") {
				err = skkMode.User.Load(value)
				if os.IsNotExist(err) {
					err = nil
				}
				if err == nil {
					skkMode.userJisyoPath = value
				}
			} else {
				err = fmt.Errorf("SKK-ERROR: unknown option: %s", key)
			}
		} else {
			err = skkMode.System.Load(token)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return skkMode, errs
}

func (M *Mode) enableUntilExit(ctx context.Context, key keys.Code, B *readline.Buffer) readline.Result {
	readline.GlobalKeyMap.BindKey(key, M)
	readline.GlobalKeyMap.BindKey(keys.Enter, &readline.GoCommand{
		Name: "SKK_ACCEPT_LINE_WITH_LATIN_MODE",
		Func: M.cmdAcceptLineWithLatinMode,
	})
	readline.GlobalKeyMap.BindKey(keys.CtrlC, &readline.GoCommand{
		Name: "SKK_INTRRUPT_WITH_LATIN_MODE",
		Func: M.cmdIntrruptWithLatinMode,
	})
	return M.Call(ctx, B)
}

func (o *onDemandLoad) Call(ctx context.Context, B *readline.Buffer) readline.Result {
	config, ok := o.Func()
	if !ok {
		readline.GlobalKeyMap.BindKey(o.Key, nil)
		return readline.CONTINUE
	}
	skkMode, errs := loadWithConfigString(config)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(B.Out, "\n%s", e.Error())
		}
		B.RepaintAll()
		readline.GlobalKeyMap.BindKey(o.Key, nil)
		return readline.CONTINUE
	}
	o.closer = func() error { return skkMode.SaveUserJisyo() }
	return skkMode.enableUntilExit(ctx, o.Key, B)
}

func SetupOnDemand(key keys.Code, f func() (string, bool)) func() error {
	o := &onDemandLoad{
		Key:  key,
		Func: f,
	}
	readline.GlobalKeyMap.BindKey(key, o)
	return o.Close
}
