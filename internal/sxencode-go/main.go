package sxencode

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
)

var (
	// Delimiters for vector (array/slice) literals in S-expression encoding.
	VectorOpen  = "("
	VectorClose = ")"
)

type Encoder struct {
	w                  io.Writer
	OnTypeNotSupported func(reflect.Value) (string, error)
}

func (enc *Encoder) writeByte(b byte) error {
	_, err := enc.w.Write([]byte{b})
	return err
}

func (enc *Encoder) write(b []byte) error {
	_, err := enc.w.Write(b)
	return err
}

func (enc *Encoder) writeString(s string) error {
	_, err := io.WriteString(enc.w, s)
	return err
}

type Sexpressioner interface {
	Sexpression() string
}

var toLispString = strings.NewReplacer(
	`"`, `\"`,
	`\`, `\\`,
)

func (enc *Encoder) tmpMarshal(value reflect.Value) (string, error) {
	var buffer strings.Builder
	tmp := *enc
	tmp.w = &buffer
	if err := tmp.encode(value); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func sxprTags(t *reflect.StructField) []string {
	tag, ok := t.Tag.Lookup("sxpr")
	if !ok {
		return nil
	}
	return strings.Split(tag, ",")
}

type tagInfoT struct {
	name      string
	omitEmpty bool
	noName    bool
	tags      []string
}

func tagInfo(t *reflect.StructField) (r *tagInfoT) {
	r = &tagInfoT{
		name: t.Name,
	}
	tags := sxprTags(t)
	if len(tags) <= 0 {
		return
	}
	if tags[0] != "" {
		r.name = tags[0]
	}
	for _, tag1 := range tags[1:] {
		if tag1 == "omitempty" {
			r.omitEmpty = true
		} else if tag1 == "noname" {
			r.noName = true
		} else {
			r.tags = append(r.tags, tag1)
		}
	}
	return
}

func (enc *Encoder) encode(value reflect.Value) error {
	k := value.Kind()
	if value.CanInterface() {
		if v, ok := value.Interface().(Sexpressioner); ok {
			_, err := io.WriteString(enc.w, v.Sexpression())
			return err
		}
	}
	switch k {
	case reflect.Interface, reflect.Pointer:
		return enc.encode(value.Elem())
	case reflect.Struct:
		if err := enc.writeByte('('); err != nil {
			return err
		}
		types := value.Type()
		fields := reflect.VisibleFields(types)
		for i, t := range fields {
			if !t.IsExported() {
				continue
			}
			fieldValue := value.Field(i)

			tag := tagInfo(&t)

			if tag.omitEmpty && fieldValue.IsZero() {
				continue
			}

			s, err := enc.tmpMarshal(fieldValue)
			if err != nil {
				return err
			}
			if tag.noName {
				_, err = fmt.Fprintf(enc.w, " %s", s)
			} else if s != "" && tag.name != "-" {
				_, err = fmt.Fprintf(enc.w, "(%s %s)", tag.name, s)
			}
			if err != nil {
				return err
			}

		}
		return enc.writeByte(')')
	case reflect.String:
		if err := enc.writeByte('"'); err != nil {
			return err
		}
		if _, err := io.WriteString(enc.w, toLispString.Replace(value.String())); err != nil {
			return err
		}
		return enc.writeByte('"')
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err := fmt.Fprint(enc.w, value.Int())
		return err
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		_, err := fmt.Fprint(enc.w, value.Uint())
		return err
	case reflect.Float32, reflect.Float64:
		_, err := fmt.Fprint(enc.w, value.Float())
		return err
	case reflect.Array, reflect.Slice:
		if err := enc.writeString(VectorOpen); err != nil {
			return err
		}
		if n := value.Len(); n > 0 {
			i := 0
			for {
				s, err := enc.tmpMarshal(value.Index(i))
				if err != nil {
					return err
				}
				if s != "" {
					enc.writeString(s)
				} else {
					enc.writeString("nil")
				}
				if i++; i >= n {
					break
				}
				if err := enc.writeByte(' '); err != nil {
					return err
				}
			}
		}
		return enc.writeString(VectorClose)
	case reflect.Map:
		iter := value.MapRange()
		enc.writeByte('(')
		for iter.Next() {
			key, err := enc.tmpMarshal(iter.Key())
			if err != nil {
				return err
			}
			val, err := enc.tmpMarshal(iter.Value())
			if err != nil {
				return err
			}
			if key != "" && val != "" {
				_, err := fmt.Fprintf(enc.w, "(%s %s)", key, val)
				if err != nil {
					return err
				}
			}
		}
		return enc.writeByte(')')
	case reflect.Bool:
		if value.Bool() {
			return enc.writeByte('t')
		} else {
			return enc.writeString("nil")
		}
	default:
		if enc.OnTypeNotSupported != nil {
			s, err := enc.OnTypeNotSupported(value)
			if err != nil {
				return err
			}
			return enc.writeString(s)
		}
	}
	return nil
}

func (enc *Encoder) Encode(v any) error {
	return enc.encode(reflect.ValueOf(v))
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func Marshal(v any) ([]byte, error) {
	var data bytes.Buffer
	enc := NewEncoder(&data)
	enc.Encode(v)
	return data.Bytes(), nil
}
