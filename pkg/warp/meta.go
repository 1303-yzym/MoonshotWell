package warp

import (
	"errors"
	"reflect"
	"strings"
)

type IMeta interface {
	meta() struct{}
}

type Meta struct{}

var _ IMeta = (*Meta)(nil)

func (m Meta) meta() struct{} { return struct{}{} }

type Nil struct{}

var _ IMeta = (*Nil)(nil)

func (m Nil) meta() struct{} { return struct{}{} }

type MetaInfo struct {
	Path        string
	Method      string
	Comment     string
	Middlewares []string
}

func getTag(st *reflect.StructTag, key string) (string, error) {
	value, ok := st.Lookup(key)
	if !ok {
		return "", errors.New("tag " + key + " is not defined")
	}

	if value == "" {
		return "", errors.New("tag " + key + " value is empty")
	}

	return value, nil
}

func parseMeta[Q IMeta](q Q) (m MetaInfo, err error) {
	tf := reflect.TypeOf(q)
	if tf.Kind() == reflect.Ptr {
		tf = tf.Elem()
	}

	if tf.NumField() == 0 {
		return m, errors.New(tf.Name() + " has no fields")
	}

	field := tf.Field(0)

	structTag := field.Tag
	if m.Path, err = getTag(&structTag, "path"); err != nil {
		return
	}

	method, err := getTag(&structTag, "method")
	if err != nil {
		return
	}

	m.Method = strings.ToUpper(method)

	if m.Comment, err = getTag(&structTag, "comment"); err != nil {
		return
	}

	return
}
