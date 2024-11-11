package errors

import "fmt"

type ObjectNotFound struct {
	id string
}

func ObjectNotFoundException(id string) error {
	return &ObjectNotFound{
		id: id,
	}
}

func (e *ObjectNotFound) Error() string {
	return fmt.Sprintf("Object with id: %s is not found", e.id)
}

type IncorrectDateTimeFormat struct {}

func IncorrectDateTimeFormatException() error {
	return &IncorrectDateTimeFormat{}
}

func (e *IncorrectDateTimeFormat) Error() string {
	return "Incorrect Date/Time format used in 'expired_at' header, Please use YYYY-MM-DD hh:mm:ss format"
}

type InvalidRequestBody struct {}

func InvalidRequestBodyException() error {
	return &InvalidRequestBody{}
}

func (e *InvalidRequestBody) Error() string {
	return "Invalid request body. Body should be in valid JSON format"
}