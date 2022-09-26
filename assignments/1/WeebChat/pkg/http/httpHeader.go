package http

import (
	"fmt"
	"strings"
	"time"
)

type HttpHeader struct {
	headers map[string]string
}

type Cookie struct {
	Name  string
	Value string

	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string

	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite bool
	Raw      string
	Unparsed []string
}

type RepresentationHeader struct {
	ContentType     string
	ContentEncoding string
	ContentLanguage []string
	ContentLocation string
	ContentLength   int
}

func NewHeader() HttpHeader {
	headers := make(map[string]string)

	httpHeader := HttpHeader{}
	httpHeader.headers = headers

	return httpHeader
}

func getRepresentationHeader(httpHeader *HttpHeader) RepresentationHeader {
	//TODO: implement
	return RepresentationHeader{}
}

func (httpHeader *HttpHeader) AddHeader(key string, value string) {
	httpHeader.headers[key] = value
}

func (httpHeader *HttpHeader) Remove(key string) {
	//TODO: implement
}

func (httpHeader *HttpHeader) Get(key string) string {

	return httpHeader.headers[key]
}

func (httpHeader *HttpHeader) String() string {
	var builder strings.Builder
	for key, value := range httpHeader.headers {
		if key == "Cookie" {
			continue
		}
		builder.WriteString(key)
		builder.WriteString(": ")
		builder.WriteString(value)
		builder.WriteString("\n")
	}

	return builder.String()
}

func (cookie *Cookie) StringSetCookie() string {
	//TODO: Set up cookies more detail
	var builder strings.Builder
	builder.WriteString("Set-Cookie: ")

	builder.WriteString(fmt.Sprintf("%s=%s; ", cookie.Name, cookie.Value))
	if cookie.Path != "" {
		builder.WriteString(fmt.Sprintf("%s=%s; ", "Path", cookie.Path))
	}
	return builder.String()

}
