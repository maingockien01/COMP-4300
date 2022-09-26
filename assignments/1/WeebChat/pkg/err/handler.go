package err

import (
	"fmt"
)

func HandleFatalError(err error) {
	if err != nil {
		panic(err)
	}
}

func LogError(err error) {
	//TODO: log error to file
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
	}
}
