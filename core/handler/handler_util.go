// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/constant"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/valyala/fasthttp"
	"reflect"
	"strconv"
	"strings"
	"syscall"
)

func parseReqWithBody(ret *reflect.Value, ctx *fasthttp.RequestCtx) error {
	var (
		retElem = ret.Elem()
		req     = &ctx.Request
	)
	if cmdArgsBytes := req.Header.Peek("*"); cmdArgsBytes != nil {
		if field := retElem.FieldByName("CmdArgs"); field.IsValid() {
			field.Set(reflect.ValueOf([]string{string(cmdArgsBytes)}))
		}
	}
	if isJsonCall(req) {
		if req.Header.ContentLength() == 0 {
			return nil
		}
		return utils.Decode(req.BodyStream(), ret.Interface())
	} else if isFormCall(req) {
		return parseForm(retElem, ctx)
	}
	return syscall.EINVAL
}

func parseValue(v reflect.Value, ctx *fasthttp.RequestCtx, cate string) (err error) {
	if v.Kind() != reflect.Ptr {
		err = errors.Info(syscall.EINVAL, "formutil.ParseValue: ret.type != pointer")
		return
	}
	v = v.Elem()
	t := v.Type()
	if v.Kind() != reflect.Struct {
		switch v.Kind() {
		case reflect.Slice: // uri like ?a=1&a=1&a=3 => function(a []int) => a [1, 2, 3]
			peeks := ctx.QueryArgs().PeekMulti(t.Name())
			if peeks == nil {
				v.Set(reflect.Zero(t))
				return
			}
			multiPeek := parseMultiPeek(peeks)
			// slice element type
			sliceType := t.Elem()
			slice := reflect.MakeSlice(sliceType, 0, 0)
			for index, peek := range multiPeek {
				err = strconvParseValue(slice.Index(index), peek)
				if err != nil {
					err = errors.Info(err, "formutil.ParseValue: parse slice field -", t.Name(), index).Detail(err)
					return
				}
			}
			v.Set(slice)
		case reflect.Map:
			if ctx.QueryArgs().Len() == 0 {
				v.Set(reflect.Zero(t))
				return
			}
			m := make(map[string]any)
			ctx.QueryArgs().VisitAll(func(key, value []byte) {
				m[string(key)] = string(value)
			})
			v.Set(reflect.ValueOf(m))
		default: // uri like ?a=1 => function(a int) => a 1
			// just handle one args
			var arg string
			ctx.URI().QueryArgs().VisitAll(func(key, value []byte) {
				arg = string(value)
			})
			err = strconvParseValue(v, arg)
		}
		return
	}
	// args is struct
	for i := 0; i < v.NumField(); i++ {
		sf := t.Field(i)
		if sf.Tag == "" { // no tag
			if sf.Anonymous {
				if err = parseValue(v.Field(i).Addr(), ctx, cate); err != nil {
					return
				}
			}
			continue
		}
		formTag := sf.Tag.Get(cate)
		if formTag == "" { // no form tag, skip
			continue
		}
		tag, opts, parseTagErr := parseTag(formTag)
		if parseTagErr != nil {
			err = errors.Info(parseTagErr, "Parse struct field:", sf.Name).Detail(parseTagErr)
			return
		}
		sfv := v.Field(i)
		peeks := ctx.QueryArgs().PeekMulti(tag)
		if opts.fhas {
			if err = setHas(v, sf.Name, peeks == nil); err != nil {
				return
			}
		}
		if peeks == nil {
			if !opts.fdefault {
				sfv.Set(reflect.Zero(sf.Type))
			}
			continue
		}
		fv := parseMultiPeek(peeks)
		if len(fv) == 0 {
			sfv.Set(reflect.Zero(sf.Type))
			continue
		}
		switch sfv.Kind() {
		case reflect.Slice:
			sft := sfv.Type()
			n := len(fv)
			slice := reflect.MakeSlice(sft, n, n)
			for j := 0; j < n; j++ {
				err = strconvParseValue(slice.Index(j), fv[j])
				if err != nil {
					err = errors.Info(err, "formutil.ParseValue: parse slice field -", sf.Name, j).Detail(err)
					return
				}
			}
			sfv.Set(slice)
		default:
			err = strconvParseValue(sfv, fv[0])
			if err != nil {
				err = errors.Info(err, "formutil.ParseValue: parse struct field -", sf.Name).Detail(err)
				return
			}
		}
	}
	return
}

func parseMultiPeek(peeks [][]byte) []string {
	fv := make([]string, len(peeks))
	for index, peek := range peeks {
		fv[index] = string(peek)
	}
	return fv
}

func isJsonCall(req *fasthttp.Request) bool {
	var ct string

	if ctBytes := req.Header.Peek(constant.HeaderContentType); ctBytes == nil {
		return false
	} else {
		ct = string(ctBytes)
	}

	return ct == "application/json" || strings.HasPrefix(ct, "application/json;")
}

func isFormCall(req *fasthttp.Request) bool {
	var ct string
	if ctBytes := req.Header.Peek(constant.HeaderContentType); ctBytes == nil {
		return false
	} else {
		ct = string(ctBytes)
	}
	return ct == "application/x-www-form-urlencoded" || strings.HasPrefix(ct, "application/x-www-form-urlencoded;") || ct == "multipart/form-data" || strings.HasPrefix(ct, "multipart/form-data;")
}

func parseForm(retElem reflect.Value, ctx *fasthttp.RequestCtx) (err error) {
	if retElem.Kind() != reflect.Ptr {
		err = errors.Info(syscall.EINVAL, "formutil.ParseValue: ret.type != pointer")
		return
	}
	v := retElem.Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		sf := t.Field(i)
		formTag := sf.Tag.Get("form")
		if formTag == "" {
			continue
		}
		sfv := v.Field(i)
		formValue := ctx.FormValue(formTag)
		if len(formValue) == 0 {
			continue
		}
		if err = strconvParseValue(sfv, string(formValue)); err != nil {
			return errors.Info(err, "formutil.ParseValue: parse form field -", sf.Name).Detail(err)
		}
	}
	return nil
}

func prefixOf(name string) (prefix string, ok bool) {

	if !constant.Is(constant.UPPER, rune(name[0])) {
		return
	}
	for i := 1; i < len(name); i++ {
		if !constant.Is(constant.LOWER, rune(name[i])) {
			return name[:i], true
		}
	}
	return
}

// --------------------------------------------------------------------

func setHas(v reflect.Value, name string, has bool) (err error) {

	sfHas := v.FieldByName("Has" + name)
	if sfHas.Kind() != reflect.Bool {
		err = errors.New("Struct filed `Has" + name + "` not found or not bool")
		return
	}
	sfHas.SetBool(has)
	return
}

type tagParseOpts struct {
	fhas     bool
	fdefault bool
}

func parseTag(tag1 string) (tag string, opts tagParseOpts, err error) {

	if tag1 == "" {
		err = errors.New("Struct field has no tag")
		return
	}

	parts := strings.Split(tag1, ",")
	tag = parts[0]
	for i := 1; i < len(parts); i++ {
		switch parts[i] {
		case "has":
			opts.fhas = true
		case "default":
			opts.fdefault = true
		case "omitempty":
		default:
			err = errors.New("Unknown tag option: " + parts[i])
			return
		}
	}
	return
}

// --------------------------------------------------------------------

// --------------------------------------------------------------------

func strconvParse(ret interface{}, str string) (err error) {

	v := reflect.ValueOf(ret)
	if v.Kind() != reflect.Ptr {
		return syscall.EINVAL
	}

	return strconvParseValue(v.Elem(), str)
}

func strconvParseValue(v reflect.Value, str string) (err error) {

	var iv int64
	var uv uint64
	var fv float64

retry:
	switch v.Kind() {
	case reflect.String:
		v.SetString(str)
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			v.SetBytes([]byte(str))
		} else {
			return syscall.EINVAL
		}
	case reflect.Int:
		iv, err = strconv.ParseInt(str, 10, 0)
		v.SetInt(iv)
	case reflect.Uint:
		uv, err = strconv.ParseUint(str, 10, 0)
		v.SetUint(uv)
	case reflect.Int64:
		iv, err = strconv.ParseInt(str, 10, 64)
		v.SetInt(iv)
	case reflect.Uint32:
		uv, err = strconv.ParseUint(str, 10, 32)
		v.SetUint(uv)
	case reflect.Int32:
		iv, err = strconv.ParseInt(str, 10, 32)
		v.SetInt(iv)
	case reflect.Uint16:
		uv, err = strconv.ParseUint(str, 10, 16)
		v.SetUint(uv)
	case reflect.Int16:
		iv, err = strconv.ParseInt(str, 10, 16)
		v.SetInt(iv)
	case reflect.Uint64:
		uv, err = strconv.ParseUint(str, 10, 64)
		v.SetUint(uv)
	case reflect.Ptr:
		elem := reflect.New(v.Type().Elem())
		v.Set(elem)
		v = elem.Elem()
		goto retry
	case reflect.Struct:
		method := v.Addr().MethodByName("ParseValue") // ParseValue(str string) error
		if method.IsValid() {
			out := method.Call([]reflect.Value{reflect.ValueOf(str)})
			ret := out[0].Interface()
			if ret != nil {
				return ret.(error)
			}
			return nil
		}
		return syscall.EINVAL
	case reflect.Uint8:
		uv, err = strconv.ParseUint(str, 10, 8)
		v.SetUint(uv)
	case reflect.Int8:
		iv, err = strconv.ParseInt(str, 10, 8)
		v.SetInt(iv)
	case reflect.Uintptr:
		uv, err = strconv.ParseUint(str, 10, 64)
		v.SetUint(uv)
	case reflect.Float64:
		fv, err = strconv.ParseFloat(str, 64)
		v.SetFloat(fv)
	case reflect.Float32:
		fv, err = strconv.ParseFloat(str, 32)
		v.SetFloat(fv)
	case reflect.Bool:
		var bv bool
		bv, err = strconv.ParseBool(str)
		v.SetBool(bv)
	default:
		return syscall.EINVAL
	}
	return
}

// --------------------------------------------------------------------
