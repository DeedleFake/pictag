//go:build !goexperiment.jsonv2

package main

import (
	"encoding/json"
	"io"
)

func marshal(w io.Writer, data any) error {
	return json.NewEncoder(w).Encode(data)
}
