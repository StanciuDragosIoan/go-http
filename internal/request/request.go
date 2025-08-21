package request

import (
	"bytes"
	"fmt"
	"io"

	"boot.theprimagen.tv/cmd/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	state       parserState
}

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *RequestLine) ValidHTTP() bool {
	return r.HttpVersion == "HTTP/1.1"
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

var ERROR_MALFORMED_START_LINE = fmt.Errorf("error malformed request-line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var ERRORR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")
var SEPARATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	//if we haven t parsed a startline
	if idx == -1 {
		return nil, 0, nil //no error we haven t done anything
	}

	//parse line
	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_START_LINE
	}

	httpParts := bytes.Split(parts[2], []byte("/"))

	// Debug: print the httpParts slice
	fmt.Println("httpParts:", httpParts)

	fmt.Println("HTTP major:", httpParts[0])
	fmt.Println("HTTP minor:", httpParts[1])

	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ERROR_MALFORMED_START_LINE
	}
	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return rl, read, nil
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer:
	for {
		currentData := data[read:]
		switch r.state {
		case StateError:
			return 0, ERRORR_REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.state = StateHeaders

		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)

			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			read += n
			if done {
				r.state = StateDone
			}
		case StateDone:
			break outer
		default:
			panic("somehow we have programmed poorly")
		}
	}

	return read, nil

}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func (r *Request) error() bool {
	return r.state == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	//note buffer could get overrun .. a header that exceeds 1k
	//or the body
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		//todo: what to do here?
		if err != nil {
			return nil, err
		}

		bufLen += n

		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
