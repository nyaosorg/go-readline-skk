package sxencode

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"strings"
)

type Decoder struct {
	r                  io.RuneScanner
	OnTypeNotSupported func(any, reflect.Value) error
}

func Unmarshal(data []byte, v any) error {
	dec := NewDecoder(bytes.NewReader(data))
	return dec.Decode(v)
}

func NewDecoder(r io.Reader) *Decoder {
	rs, ok := r.(io.RuneScanner)
	if !ok {
		rs = bufio.NewReader(r)
	}
	return &Decoder{r: rs}
}

func (D *Decoder) Decode(v any) error {
	sxpr, err := parser1.Read(D.r)
	if err != nil {
		return err
	}
	return D.decode(sxpr, reflect.ValueOf(v))
}

func keyValuePair(value any) *consT {
	cons, ok := value.(*consT)
	if !ok {
		return nil
	}
	if cdr, ok := cons.Cdr.(*consT); ok {
		return &consT{Car: cons.Car, Cdr: cdr.Car}
	}
	return cons
}

func sxpr2list(sxpr any) (pairs []*consT, nopairs []any) {
	for sxpr != nil {
		cons, ok := sxpr.(*consT)
		if !ok {
			return
		}
		if v := keyValuePair(cons.Car); v != nil {
			pairs = append(pairs, v)
		} else {
			nopairs = append(nopairs, cons.Car)
		}
		sxpr = cons.Cdr
	}
	return
}

func (D *Decoder) decode(sxpr any, value reflect.Value) error {
	switch value.Kind() {
	case reflect.Pointer:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		return D.decode(sxpr, value.Elem())
	case reflect.Interface:
		return D.decode(sxpr, value.Elem())
	case reflect.Struct:
		pairs, nopairs := sxpr2list(sxpr)
		if len(pairs) <= 0 && len(nopairs) <= 0 {
			return nil
		}
		types := value.Type()
		fields := reflect.VisibleFields(types)
		for _, sxpr1 := range pairs {
			var key string
			if v, ok := sxpr1.Car.(symbolT); ok {
				key = v.Value
			} else if v, ok := sxpr1.Car.(string); ok {
				key = v
			} else {
				continue
			}
			for i, field1 := range fields {
				tag := tagInfo(&field1)
				if !tag.noName && strings.EqualFold(key, tag.name) {
					D.decode(sxpr1.Cdr, value.Field(i))
					break
				}
			}
		}
		used := 0
		for _, sxpr1 := range nopairs {
			for i := used; i < len(fields); i++ {
				tag := tagInfo(&fields[i])
				if tag.noName {
					D.decode(sxpr1, value.Field(i))
					used = i + 1
					break
				}
			}
		}
	case reflect.String:
		if v, ok := sxpr.(string); ok {
			value.SetString(v)
		} else if v, ok := sxpr.(symbolT); ok {
			value.SetString(v.Value)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, ok := sxpr.(int64); ok {
			value.SetInt(v)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if v, ok := sxpr.(uint64); ok {
			value.SetUint(uint64(v))
		}
	case reflect.Float32, reflect.Float64:
		if v, ok := sxpr.(float64); ok {
			value.SetFloat(v)
		}
	case reflect.Array, reflect.Slice:
		source, ok := sxpr.([]any)
		if !ok {
			for sxpr != nil {
				cons, ok := sxpr.(*consT)
				if !ok {
					source = append(source, sxpr)
					break
				}
				source = append(source, cons.Car)
				sxpr = cons.Cdr
			}
			if source == nil {
				break
			}
		}
		elemType := value.Type().Elem()
		slice := reflect.MakeSlice(value.Type(), 0, len(source))
		for _, s := range source {
			newElem := reflect.New(elemType).Elem()
			if err := D.decode(s, newElem); err != nil {
				return err
			}
			slice = reflect.Append(slice, newElem)
		}
		value.Set(slice)
	case reflect.Map:
		sxprs, _ := sxpr2list(sxpr)
		if len(sxprs) <= 0 {
			return nil
		}
		if value.IsNil() {
			map1 := reflect.MakeMap(value.Type())
			value.Set(map1)
		}
		eleType := value.Type().Elem()
		keyType := value.Type().Key()
		for _, sxpr1 := range sxprs {
			var key any = sxpr1.Car
			if sym, ok := sxpr1.Car.(symbolT); ok {
				key = sym.Value
			}
			keyValue := reflect.New(keyType).Elem()
			if err := D.decode(key, keyValue); err != nil {
				return err
			}
			val := sxpr1.Cdr
			eleValue := reflect.New(eleType).Elem()
			if err := D.decode(val, eleValue); err != nil {
				return err
			}
			value.SetMapIndex(keyValue, eleValue)
		}
	case reflect.Bool:
		if v, ok := sxpr.(bool); ok {
			value.SetBool(v)
		}
	default:
		if D.OnTypeNotSupported != nil {
			err := D.OnTypeNotSupported(sxpr, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
