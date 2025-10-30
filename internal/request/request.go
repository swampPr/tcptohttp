// Package request provides request  INFO:  request parsing
package request

import (
	"bytes"
	"fmt"
	"io"
)

// RequestLine struct  INFO:
type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

// Request struct  INFO:
type Request struct {
	RequestLine RequestLine
	state       parserState
}

func newRequest() *Request {
	return &Request{state: StateInit}
}

var (
	ErrorMalformedRequestLine   = fmt.Errorf("malformed request line")
	ErrorUnsupportedHTTPVersion = fmt.Errorf("unsupported HTTP version")
	ErrorRequestInErrorState    = fmt.Errorf("request in error state")
	Separator                   = []byte("\r\n")
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
	StateErr  parserState = "error"
)

func parseRequestLine(s []byte) (*RequestLine, int, error) {
	idx := bytes.Index(s, Separator)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := s[:idx]
	read := idx + len(Separator)

	parts := bytes.Fields(startLine)
	if len(parts) != 3 {
		return nil, 0, ErrorMalformedRequestLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[1]) != "1.1" || string(httpParts[0]) != "HTTP" {
		return nil, 0, ErrorMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HTTPVersion:   string(httpParts[1]),
	}

	return rl, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateErr:
			return 0, ErrorRequestInErrorState
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.state = StateErr
				return 0, err
			}
			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.state = StateDone
		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateErr
}

// RequestFromReader function  INFO:
func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	// NOTE: buffer could get overrun... a header that exceeds 1k would do that
	// Or the body
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		// TODO: idk what to do here
		if err != nil {
			return nil, err
		}

		bufLen += n

		readN, err := request.parse(buf[:bufLen+n])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN

	}

	return request, nil
}
