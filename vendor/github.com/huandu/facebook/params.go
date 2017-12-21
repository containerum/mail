// A facebook graph api client in go.
// https://github.com/huandu/facebook/
//
// Copyright 2012 - 2015, Huan Du
// Licensed under the MIT license
// https://github.com/huandu/facebook/blob/master/LICENSE

package facebook

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
)

const (
	_MIME_FORM_URLENCODED = "application/x-www-form-urlencoded"
	_MIME_FORM_DATA       = "multipart/form-data"
)

var (
	typeOfPointerToBinaryData = reflect.TypeOf(&binaryData{})
	typeOfPointerToBinaryFile = reflect.TypeOf(&binaryFile{})
)

// API params.
//
// For general uses, just use Params as an ordinary map.
//
// For advanced uses, use MakeParams to create Params from any struct.
type Params map[string]interface{}

// Makes a new Params instance by given data.
// Data must be a struct or a map with string keys.
// MakeParams will change all struct field name to lower case name with underscore.
// e.g. "FooBar" will be changed to "foo_bar".
//
// Returns nil if data cannot be used to make a Params instance.
func MakeParams(data interface{}) (params Params) {
	if p, ok := data.(Params); ok {
		return p
	}

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}

			params = nil
		}
	}()

	params = makeParams(reflect.ValueOf(data))
	return
}

func makeParams(value reflect.Value) (params Params) {
	for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	// only map with string keys can be converted to Params
	if value.Kind() == reflect.Map && value.Type().Key().Kind() == reflect.String {
		params = Params{}

		for _, key := range value.MapKeys() {
			params[key.String()] = value.MapIndex(key).Interface()
		}

		return
	}

	if value.Kind() != reflect.Struct {
		return
	}

	params = Params{}
	num := value.NumField()

	for i := 0; i < num; i++ {
		name := camelCaseToUnderScore(value.Type().Field(i).Name)
		field := value.Field(i)

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		switch field.Kind() {
		case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Invalid:
			// these types won't be marshalled in json.
			params = nil
			return

		default:
			params[name] = field.Interface()
		}
	}

	return
}

// Encodes params to query string.
// If map value is not a string, Encode uses json.Marshal() to convert value to string.
//
// Encode will panic if Params contains values that cannot be marshalled to json string.
func (params Params) Encode(writer io.Writer) (mime string, err error) {
	if params == nil || len(params) == 0 {
		mime = _MIME_FORM_URLENCODED
		return
	}

	// check whether params contains any binary data.
	hasBinary := false

	for _, v := range params {
		typ := reflect.TypeOf(v)

		if typ == typeOfPointerToBinaryData || typ == typeOfPointerToBinaryFile {
			hasBinary = true
			break
		}
	}

	if hasBinary {
		return params.encodeMultipartForm(writer)
	}

	return params.encodeFormUrlEncoded(writer)
}

func (params Params) encodeFormUrlEncoded(writer io.Writer) (mime string, err error) {
	var jsonStr []byte
	written := false

	for k, v := range params {
		if written {
			io.WriteString(writer, "&")
		}

		io.WriteString(writer, url.QueryEscape(k))
		io.WriteString(writer, "=")

		if reflect.TypeOf(v).Kind() == reflect.String {
			io.WriteString(writer, url.QueryEscape(reflect.ValueOf(v).String()))
		} else {
			jsonStr, err = json.Marshal(v)

			if err != nil {
				return
			}

			io.WriteString(writer, url.QueryEscape(string(jsonStr)))
		}

		written = true
	}

	mime = _MIME_FORM_URLENCODED
	return
}

func (params Params) encodeMultipartForm(writer io.Writer) (mime string, err error) {
	w := multipart.NewWriter(writer)
	defer func() {
		w.Close()
		mime = w.FormDataContentType()
	}()

	for k, v := range params {
		switch value := v.(type) {
		case *binaryData:
			var dst io.Writer
			filePart := createFormFile(k, value.Filename, value.ContentType)
			dst, err = w.CreatePart(filePart)

			if err != nil {
				return
			}

			_, err = io.Copy(dst, value.Source)

			if err != nil {
				return
			}

		case *binaryFile:
			var dst io.Writer
			var file *os.File
			var path string

			filePart := createFormFile(k, value.Filename, value.ContentType)
			dst, err = w.CreatePart(filePart)

			if err != nil {
				return
			}

			if value.Path == "" {
				path = value.Filename
			} else {
				path = value.Path
			}

			file, err = os.Open(path)

			if err != nil {
				return
			}

			_, err = io.Copy(dst, file)

			if err != nil {
				return
			}

		default:
			var dst io.Writer
			var jsonStr []byte

			dst, err = w.CreateFormField(k)

			if reflect.TypeOf(v).Kind() == reflect.String {
				io.WriteString(dst, reflect.ValueOf(v).String())
			} else {
				jsonStr, err = json.Marshal(v)

				if err != nil {
					return
				}

				_, err = dst.Write(jsonStr)

				if err != nil {
					return
				}
			}
		}
	}

	return
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func createFormFile(fieldName, fileName, contentType string) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			quoteEscaper.Replace(fieldName), quoteEscaper.Replace(fileName)))

	if contentType == "" {
		contentType = mime.TypeByExtension(path.Ext(fileName))

		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	h.Set("Content-Type", contentType)
	return h
}
