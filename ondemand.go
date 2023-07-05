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

func (o *onDemandLoad) Call(ctx context.Context, B *readline.Buffer) readline.Result {
	env, ok := o.Func()
	if !ok {
		readline.GlobalKeyMap.BindKey(o.Key, nil)
		return readline.CONTINUE
	}
	skkMode := New()
	ok = true
	failed := false
	succeeded := false
	for ok {
		var token string
		token, env, ok = strings.Cut(env, ";")

		key, value, hasEqual := strings.Cut(token, "=")
		var err error
		if hasEqual {
			if strings.EqualFold(key, "user") {
				err = skkMode.User.Load(value)
				if os.IsNotExist(err) {
					err = nil
				}
				if err == nil {
					o.closer = func() error {
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
			fmt.Fprintf(B.Out, "\n%s", err.Error())
			failed = true
		} else {
			succeeded = true
		}
	}
	if failed {
		B.RepaintAll()
	}
	if succeeded {
		readline.GlobalKeyMap.BindKey(o.Key, skkMode)
		return skkMode.Call(ctx, B)
	}
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
