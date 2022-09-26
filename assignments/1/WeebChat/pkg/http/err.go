package http

type HttpRequestInvalidError struct {}

func (err HttpRequestInvalidError) Error () string {
    return "Request is not in HTTP format"
}

