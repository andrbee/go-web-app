package utils

import "errors"

var (
	INTERNAL_SERVER_ERROR = errors.New("Internal Server Error")
	SYSTEM_TRY_OPERATION_LATER = errors.New("Sorry, try later")
)
