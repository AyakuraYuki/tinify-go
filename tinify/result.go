package tinify

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type Result struct {
	*Meta
	data []byte
}

func NewResult(header http.Header, data []byte) *Result {
	r := &Result{}
	r.Meta = NewMeta(header)
	r.data = data
	return r
}

func (r *Result) Data() []byte {
	if r == nil {
		return nil
	}
	return r.data
}

func (r *Result) ToBuffer() []byte {
	return r.Data()
}

func (r *Result) ToFile(dst string) (err error) {
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, r.data, os.ModePerm)
	return err
}

func (r *Result) Size() int64 {
	if r == nil {
		return 0
	}
	return r.Meta.Size()
}

func (r *Result) MediaType() string {
	return r.ContentType()
}

func (r *Result) ContentType() string {
	if r == nil {
		return ""
	}
	return r.Meta.MimeType()
}

type ErrorData struct {
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (e *ErrorData) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("[%s] %s", e.Err, e.Message)
}
