//go:build goexperiment.jsonv2

package main

import (
	"encoding/json/v2"
	"io"
)

func marshal(w io.Writer, data any) error {
	return json.MarshalWrite(w, data)
}
