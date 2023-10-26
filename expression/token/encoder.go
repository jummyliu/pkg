package token

import (
	"fmt"
	"strings"

	"github.com/jummyliu/pkg/number"
)

type Encodable interface {
	Encode(any) string
	Decode(string) any
}

type emptyEncoder struct{}

func (encoder emptyEncoder) Encode(val any) string {
	return val.(string)
}
func (encoder emptyEncoder) Decode(val string) any {
	return val
}

type stringEncoder struct{}

func (encoder stringEncoder) Encode(val any) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(val.(string), "'", "\\'"))
}
func (encoder stringEncoder) Decode(val string) any {
	if len(val) < 2 {
		return val
	}
	return val[1 : len(val)-1]
}

type numberEncoder struct{}

func (encoder numberEncoder) Encode(val any) string {
	return fmt.Sprintf("%f", val.(float64))
}
func (encoder numberEncoder) Decode(val string) any {
	return number.ParseFloat[float64](val)
}

type boolEncoder struct{}

func (encoder boolEncoder) Encode(val any) string {
	return fmt.Sprintf("%v", val)
}
func (encoder boolEncoder) Decode(val string) any {
	return val == "true"
}

var (
	EmptyEncoder  emptyEncoder
	StringEncoder stringEncoder
	NumberEncoder numberEncoder
	BoolEncoder   boolEncoder
)
