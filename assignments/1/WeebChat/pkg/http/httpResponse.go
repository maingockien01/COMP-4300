package http

import (
	"strings"
)

type HttpResponse struct {
	HttpVersion string
	StatusCode  string
	StatusMsg   string
	Headers     HttpHeader
	Body        HttpBody

	Request HttpRequest

	cookies []Cookie
}

func (res *HttpResponse) GetCookies() []Cookie {
	return res.cookies
}

func (res *HttpResponse) GetCookie(name string) (cookie Cookie) {
	cookie = Cookie{}

	for _, value := range res.cookies {
		if value.Name == name {
			cookie = value
		}
	}

	return
}

func (res *HttpResponse) AddCookie(cookie Cookie) {
	res.cookies = append(res.cookies, cookie)
}

func (res *HttpResponse) String() string {
	var builder strings.Builder

	builder.WriteString(res.stringResponseLine())

	builder.WriteString("\n")

	//Headers
	builder.WriteString(res.Headers.String())
	//Set cookies
	for _, cookie := range res.cookies {
		builder.WriteString(cookie.StringSetCookie())
		builder.WriteString("\n")
	}

	//Empty line
	builder.WriteString("")
	builder.WriteString("\n")

	//Body
	builder.WriteString(res.Body.String())

	return builder.String()
}

func (res *HttpResponse) stringResponseLine() string {
	var builder strings.Builder

	builder.WriteString("HTTP/")
	builder.WriteString(res.HttpVersion)
	builder.WriteString(" ")
	builder.WriteString(res.StatusCode)
	builder.WriteString(" ")
	builder.WriteString(res.StatusMsg)

	return builder.String()
}

func (res *HttpResponse) SetBody(raw []byte) {

	res.Body = HttpBody{}
	res.Body.Raw = raw

	// res.Headers.AddHeader("Content-length", string(rune(len(raw))))
}
