package tcp

type MissingRequiredField struct {}

func (err MissingRequiredField) Error () string {
    return "Some required fields are nil"
}