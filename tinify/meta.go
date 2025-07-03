package tinify

import (
	"net/http"
	"strconv"
)

type Meta struct {
	header http.Header
}

func NewMeta(header http.Header) *Meta {
	return &Meta{header: header}
}

func (m *Meta) Width() int64 {
	w := m.header["Image-Width"]
	if len(w) == 0 {
		return 0
	}
	width, _ := strconv.Atoi(w[0])
	return int64(width)
}

func (m *Meta) Height() int64 {
	h := m.header["Image-Height"]
	if len(h) == 0 {
		return 0
	}
	height, _ := strconv.Atoi(h[0])
	return int64(height)
}

func (m *Meta) MimeType() string {
	ct := m.header["Content-Type"]
	if len(ct) == 0 {
		return ""
	}
	return ct[0]
}

func (m *Meta) Size() int64 {
	cl := m.header["Content-Length"]
	if len(cl) == 0 {
		return 0
	}
	size, _ := strconv.Atoi(cl[0])
	return int64(size)
}

func (m *Meta) CompressionCount() int64 {
	cc := m.header["Compression-Count"]
	if len(cc) == 0 {
		return 0
	}
	count, _ := strconv.Atoi(cc[0])
	return int64(count)
}

func (m *Meta) Location() string {
	vs := m.header["Location"]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}
