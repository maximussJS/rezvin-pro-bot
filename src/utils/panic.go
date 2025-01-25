package utils

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicIfNotContextError(err error) {
	if err != nil && !IsContextError(err) {
		panic(err)
	}
}

func PanicIfNotRecordNotFound(err error) {
	if err != nil && !IsRecordNotFoundError(err) {
		panic(err)
	}
}
