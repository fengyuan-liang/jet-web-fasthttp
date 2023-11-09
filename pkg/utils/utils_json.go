package utils

import (
	"bytes"
	"encoding/gob"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	jsoniter "github.com/json-iterator/go"
	"io"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(in interface{}) (str string, err error) {
	var (
		buf []byte
	)
	buf, err = json.Marshal(in)
	if err != nil {
		return
	}
	str = string(buf)
	return
}

func ObjToByte(in interface{}) (buf []byte, err error) {
	buf, err = json.Marshal(in)
	return
}

// MustObjToByte not json
func MustObjToByte(in interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(in); err != nil {
		xlog.Errorf("MustObjToByte failed, err:%v", err.Error())
		return nil
	}
	return buf.Bytes()
}

func ByteToObj(buf []byte, out interface{}) (err error) {
	dc := json.NewDecoder(bytes.NewReader(buf))
	dc.UseNumber()
	return dc.Decode(out)
}

func Unmarshal(in string, out interface{}) error {
	//return json.Unmarshal([]byte(in), out)
	dc := json.NewDecoder(strings.NewReader(in))
	dc.UseNumber()
	return dc.Decode(out)
}

func ObjToMap(in interface{}) map[string]interface{} {
	var (
		maps map[string]interface{}
		buf  []byte
		err  error
	)
	if buf, err = json.Marshal(in); err != nil {
		//fmt.Println(err)
	} else {
		d := json.NewDecoder(bytes.NewReader(buf))
		d.UseNumber()
		if err = d.Decode(&maps); err != nil {
			//fmt.Println(err)
		} else {
			for k, v := range maps {
				maps[k] = v
			}
		}
	}
	return maps
}

func MapToObj(maps map[string]interface{}, out interface{}) error {
	buf, err := json.Marshal(maps)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, out)
	if err != nil {
		return err
	}

	return nil
}

func ObjToJsonStr(obj interface{}) (str string) {
	str = ""
	var err error
	if obj == nil {
		xlog.Error("obj is nil")
		return
	}
	str, err = Marshal(obj)
	if err != nil {
		xlog.Error(err.Error())
		return
	}
	return
}

func Decode(reader io.Reader, obj interface{}) error {
	return json.NewDecoder(reader).Decode(obj)
}
