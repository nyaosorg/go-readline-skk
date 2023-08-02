package skk

import (
	"context"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

type onDemandLoad struct {
	Key    keys.Code
	Func   func(*Mode) bool
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
	skkMode := New()
	ok := o.Func(skkMode)
	if !ok {
		readline.GlobalKeyMap.BindKey(o.Key, nil)
		return readline.CONTINUE
	}
	o.closer = func() error { return skkMode.SaveUserJisyo() }
	return skkMode.enableUntilExit(ctx, o.Key, B)
}

func SetupOnDemand(key keys.Code, f func(*Mode) bool) func() error {
	o := &onDemandLoad{
		Key:  key,
		Func: f,
	}
	readline.GlobalKeyMap.BindKey(key, o)
	return o.Close
}
