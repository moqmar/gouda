package gouda

// AssertNil panics if something is not nil. Useful for error handling.
func AssertNil(something interface{}) {
	if something != nil {
		panic(something)
	}
}

// Unwrap returns the first element if the second one is nil, otherwise panics. Useful for error handling.
func Unwrap(something interface{}, err interface{}) interface{} {
	AssertNil(err)
	return something
}
