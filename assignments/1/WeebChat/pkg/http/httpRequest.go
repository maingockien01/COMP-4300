package http

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type HttpRequest struct {
	Method       string
	URL          string
	HttpVersion  string
	VersionMajor int
	VersionMinor int
	Headers      HttpHeader
	Body         HttpBody

	RemoteAddr string
	RequestURI string

	cookies []Cookie
}

//Throws: HttpRequestInvalidError
func ParseRequest(message string) (req HttpRequest, err error) {
	err = HttpRequestInvalidError{}
	req = HttpRequest{}
	//Get scanner to read line by line
	scanner := bufio.NewScanner(strings.NewReader(message))
	// scanner.Split(bufio.ScanLines)
	fmt.Println(message)
	//Parse the request line
	scanner.Scan()
	requestLine := scanner.Text()

	req.Method, req.URL, req.HttpVersion, req.VersionMajor, req.VersionMinor, err = parseRequestLine(requestLine)

	if err != nil {
		return
	}

	//Parse headers
	req.Headers = NewHeader()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		//If line is empty which is end of header section
		if line == "" {
			break
		}

		//Parse header line
		key, value, parseErr := parseHeaderLine(line)

		if parseErr != nil {
			return
		}
		req.Headers.AddHeader(key, value)
		//Parse cookie if header is cookie
		if key == "Cookie" {
			cookies, parseErr := parseCookie(value)

			if parseErr != nil {
				return
			}

			req.cookies = append(req.cookies, cookies...)
		}
	}

	//Parse body
	//TODO
	bodyRaw := ""
	for scanner.Scan() {
		line := scanner.Text()
		bodyRaw += strings.TrimSpace(line)
	}
	req.Body.Raw = []byte(strings.TrimSpace(bodyRaw))

	err = nil
	return
}

//Throws: HttpRequestInvalidError
func parseRequestLine(requestLine string) (method string, url string, version string, majorVer int, minorVer int, err error) {
	//Sanitize parameters
	requestLine = strings.TrimSpace(requestLine)
	//Split request line into tokens
	tokens := strings.Split(requestLine, " ")

	//Assume all request line must have 3 fileds: Method URI HTTPVersion
	if len(tokens) != 3 {
		err = HttpRequestInvalidError{}
		return
	}

	method = tokens[0]

	url = tokens[1]

	httpVersion := tokens[2]

	if !strings.HasPrefix(httpVersion, "HTTP") {
		err = HttpRequestInvalidError{}
		return
	}

	httpVersionTokens := strings.Split(httpVersion, "/")

	if len(httpVersionTokens) != 2 {
		err = HttpRequestInvalidError{}
		return
	}

	version = httpVersionTokens[1]

	versionTokens := strings.Split(version, ".")

	if len(versionTokens) != 2 {
		err = HttpRequestInvalidError{}
		return
	}

	majorVer, parseErr := strconv.Atoi(versionTokens[0])

	if parseErr != nil {
		err = HttpRequestInvalidError{}
		return
	}

	minorVer, parseErr = strconv.Atoi(versionTokens[1])

	if parseErr != nil {
		err = HttpRequestInvalidError{}
		return
	}

	return

}

//Throws HttpRequestInvalidError
func parseHeaderLine(headerLine string) (key string, value string, err error) {
	//Init return value
	err = HttpRequestInvalidError{}

	headerPair := strings.Split(headerLine, ": ")

	if len(headerPair) != 2 {
		return
	}

	key = strings.TrimSpace(headerPair[0])
	value = strings.TrimSpace(headerPair[1])

	err = nil

	return

}

//Throws HttpRequestInvalidError
func parseCookie(cookiesRaw string) (cookies []Cookie, err error) {
	err = HttpRequestInvalidError{}

	cookies = make([]Cookie, 0)

	cookiesTokens := strings.Split(cookiesRaw, ";")

	for _, cookieRaw := range cookiesTokens {
		//Sanitize cookie raw
		cookieRaw = strings.TrimSpace(cookieRaw)
		//Split to key - value pair
		cookieTokens := strings.Split(cookieRaw, "=")

		if len(cookieTokens) != 2 {
			return
		}

		cookie := Cookie{}

		cookie.Name = strings.TrimSpace(cookieTokens[0])
		cookie.Value = strings.TrimSpace(cookieTokens[1])

		cookies = append(cookies, cookie)
	}

	err = nil
	return
}

func (req *HttpRequest) String() string {
	//TODO
	return ""
}

func (req *HttpRequest) GetHeader(field string) (value string) {
	value = req.Headers.Get(field)
	return
}

func (req *HttpRequest) GetCookies() []Cookie {
	return req.cookies
}

func (req *HttpRequest) GetCookie(name string) (cookie Cookie) {

	for _, value := range req.cookies {
		if value.Name == name {
			cookie = value
			return
		}
	}

	cookie = Cookie{Name: "NOT_EXIST"}
	return
}

func (req *HttpRequest) AddCookie(cookie Cookie) {
	req.cookies = append(req.cookies, cookie)
}
