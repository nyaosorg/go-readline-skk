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

func loadWithConfigString(config string) (skkMode *Mode, closer func() error, errs []error) {
	skkMode = New()
	closer = func() error { return nil }
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
					closer = func() error {
						return skkMode.SaveUserJisyo(value)
					}
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
	return skkMode, closer, errs
}

func (o *onDemandLoad) Call(ctx context.Context, B *readline.Buffer) readline.Result {
	config, ok := o.Func()
	if !ok {
		readline.GlobalKeyMap.BindKey(o.Key, nil)
		return readline.CONTINUE
	}
	skkMode, closer, errs := loadWithConfigString(config)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(B.Out, "\n%s", e.Error())
		}
		B.RepaintAll()
		return readline.CONTINUE
	}
	o.closer = closer
	readline.GlobalKeyMap.BindKey(o.Key, skkMode)
	readline.GlobalKeyMap.BindKey(keys.Enter, &readline.GoCommand{
		Name: "SKK_ACCEPT_LINE_WITH_LATIN_MODE",
		Func: skkMode.cmdAcceptLineWithLatinMode,
	})
	readline.GlobalKeyMap.BindKey(keys.CtrlC, &readline.GoCommand{
		Name: "SKK_INTRRUPT_WITH_LATIN_MODE",
		Func: skkMode.cmdIntrruptWithLatinMode,
	})
	skkMode.Call(ctx, B)
	return readline.CONTINUE
}

func SetupOnDemand(key keys.Code, f func() (string, bool)) func() error {
	o := &onDemandLoad{
		Key:  key,
		Func: f,
	}
	readline.GlobalKeyMap.BindKey(key, o)
	return o.Close
}
