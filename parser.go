//Copyright 2015 ideazxy, Ivan Korostelyov
//
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package iso8583

import (
	"errors"
	"fmt"
	"reflect"
)

// Parser for ISO 8583 messages
type Parser struct {
	messages  map[string]reflect.Type
	MtiEncode int
}

// Register MTI
func (p *Parser) Register(mti string, tpl interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Critical error:" + fmt.Sprint(r))
		}
	}()

	if len(mti) != 4 {
		return errors.New("MTI must be a 4 digit numeric field")
	}
	v := reflect.ValueOf(tpl)
	// TODO do more check
	if p.messages == nil {
		p.messages = make(map[string]reflect.Type)
	}
	p.messages[mti] = reflect.Indirect(v).Type()

	return nil
}

func decodeMti(raw []byte, encode int) (string, error) {
	mtiLen := 4
	if encode == BCD {
		mtiLen = 2
	}
	if len(raw) < mtiLen {
		return "", errors.New("bad MTI raw data")
	}

	var mti string
	switch encode {
	case ASCII:
		mti = string(raw[:mtiLen])
	case BCD:
		mti = string(bcd2Ascii(raw[:mtiLen]))
	default:
		return "", errors.New("invalid encode type")
	}
	return mti, nil
}

//Parse MTI
func (p *Parser) Parse(raw []byte) (ret *Message, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Critical error:" + fmt.Sprint(r))
			ret = nil
		}
	}()

	mti, err := decodeMti(raw, p.MtiEncode)
	if err != nil {
		return nil, err
	}

	tp, ok := p.messages[mti]
	if !ok {
		return nil, errors.New("no template registered for MTI: " + mti)
	}
	tpl := reflect.New(tp)
	initStruct(tp, tpl)
	msg := NewMessage(mti, tpl.Interface())
	msg.MtiEncode = p.MtiEncode
	return msg, msg.Load(raw)
}

func initStruct(tp reflect.Type, val reflect.Value) {
	for i := 0; i < tp.NumField(); i++ {
		field := reflect.Indirect(val).Field(i)
		fieldType := tp.Field(i)
		switch fieldType.Type.Kind() {
		case reflect.Ptr: // only initialize Ptr fields
			fieldValue := reflect.New(fieldType.Type.Elem())
			field.Set(fieldValue)
		}
	}
}
