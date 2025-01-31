package encoding

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"runtime"
	fuzz "github.com/AdaLogics/go-fuzz-headers"
)

func FuzzEncoding(data []byte) int {
	f := fuzz.NewConsumer(data)
	decType, err := f.GetInt()
	if err != nil {
		return 0
	}
	b1, err := f.GetBytes()
	if err != nil {
		return 0
	}
	b2, err := f.GetBytes()
	if err != nil {
		return 0
	}
	defer func() {
		if r := recover(); r != nil {
		}
		runtime.GC()
	}()
	switch decType%5 {
	case 0:
		e, err := f.GetString()
		if err != nil || len(e) != 32 {
			return 0
		}
		enc := base32.NewEncoding(e)
		d := base32.NewDecoder(enc, bytes.NewReader(b1))
		_, _ = d.Read(b2)
		return 1
	case 1:
		e, err := f.GetString()
		if err != nil || len(e) != 64 {
			return 0
		}
		for i := 0; i < len(e); i++ {
			if e[i] == '\n' || e[i] == '\r' {
				return 0
			}
		}
		enc := base64.NewEncoding(e)
		d := base64.NewDecoder(enc, bytes.NewReader(b1))
		_, _ = d.Read(b2)
		return 1
	case 2:
		d := gob.NewDecoder(bytes.NewReader(b1))
		_ = d.Decode(b2)
		return 1
	case 3:
		d := json.NewDecoder(bytes.NewReader(b1))
		_ = d.Decode(b2)
		return 1
	case 4:
		d := xml.NewDecoder(bytes.NewReader(b1))
		_, _ = d.Token()
		return 1
	}
	return 1
}
