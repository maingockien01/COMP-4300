package http

import (
    "strings"
)

type HttpBody struct {
    Raw             []byte
    Text            string              //For text type
    Application     []byte              //For application type
    Media           []byte
    Image           []byte
    //Other types but dont need for now

    form            map[string]string

    multiforms      []map[string]string
}

var parsers = map[string]func(string)(HttpBody, error) {
    "application/x-www-form-urlencoded" : ParseApplicationFormUrlEncoded,
}

func ParseHttpBody (bodyString string, representHeaders RepresentationHeader) (body HttpBody, err error) {
    parser := parsers[representHeaders.ContentType]

    if parser != nil {
        return parser(bodyString)
    }
    
    return body, err
}

func ParseApplicationFormUrlEncoded (bodyString string) (body HttpBody, err error) {
    body.Application = []byte(bodyString)
    
    pairs := strings.Split (bodyString, "&")

    body.form = make(map[string]string)

    for _, pair := range pairs {
        keyValue := strings.Split(pair, "=")
        if len(keyValue) != 2 {
            continue
        }
        key := keyValue[0]
        value := keyValue[1]

        body.form[key] = value
    }

    return body, err

}

func (body HttpBody) GetForm (key string) (value string) {
    if body.form == nil {
        return ""
    }

    return body.form[key]
}

func ParseMultiForm (bodyString string) (body HttpBody, err error) {
    //TODO: to be implemented
    return body, err
}

//TODO: I dont know if its working now
func (body *HttpBody) String () string {
    var builder strings.Builder

    builder.WriteString(string(body.Raw))
    return builder.String()
}
