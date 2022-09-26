package err

func HandleConnError(err error) {
	if err != nil {
		LogError(err)
	}
}
