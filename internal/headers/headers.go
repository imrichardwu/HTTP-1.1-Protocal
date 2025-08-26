package headers

import (
	"bytes"
	"fmt"
	"strings"
)

func isToken(str []byte) bool {
	found := false
	for _, ch := range str {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			found = true
		} else if ch == '!' || ch == '#' || ch == '$' || ch == '%' || ch == '&' || ch == '\'' || ch == '*' || ch == '+' || ch == '-' || ch == '.' || ch == '^' || ch == '_' || ch == '`' || ch == '|' || ch == '~' {
			found = true
		} else {
			return false
		}
	}
	if !found {
		return found
	}

	return true
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("melformed header")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])
	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("melformed field name")
	}

	return string(name), string(value), nil
}

type Headers struct {
	headers map[string]string
}

var rn = []byte("\r\n")

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) (string, bool) {
	str, ok := h.headers[strings.ToLower(name)]
	return str, ok
}

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)
	if v, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s, %s", v, value)
	} else {
		h.headers[name] = value
	}
}

func (h *Headers) Replace(name, value string) {
	name = strings.ToLower(name)
	h.headers[name] = value
}

func (h *Headers) ForEach(cb func(name, value string)) {
	for names, value := range h.headers {
		cb(names, value)
	}
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break
		}

		// Empty Header
		if idx == 0 {
			done = true
			read += len(rn)
			break
		}

		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}

		if !isToken([]byte(name)) {
			return 0, false, fmt.Errorf("malformed header name")
		}

		// Skip past the "\r\n" and continue to look for more fileline
		read += idx + len(rn)

		h.Set(name, value)

	}
	return read, done, nil
}
