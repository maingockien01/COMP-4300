package http

type MissingRequiredField struct {}

func (err MissingRequiredField) Error () string {
    return "Some required fields are nil"
}

type HttpRequestInvalidError struct {}

func (err HttpRequestInvalidError) Error () string {
    return "Request is not in HTTP format"
}

