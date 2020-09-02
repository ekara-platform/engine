package model

import (
	"encoding/json"
	"reflect"
)

type (
	//ErrorType type used to represent the type of a validation error
	ErrorType int

	//ValidationErrors represents a list of all error resulting of the construction
	//or the validation of an environment
	ValidationErrors struct {
		Errors []ValidationError
	}

	//ValidationError represents an error created during the construction of the
	// validation or an environment
	ValidationError struct {
		// ErrorType represents the type of the error
		ErrorType ErrorType
		// Location represents the place, within the descriptor, where the error occurred
		Location DescriptorLocation
		// Message represents a human readable message telling what need to be
		// fixed into the descriptor to get rid of this error
		Message string
	}

	validatable interface {
		validate(e Environment, loc DescriptorLocation) ValidationErrors
	}
)

// Error returns the message resulting of the concatenation of all included ValidationError(s)
func (ve ValidationErrors) Error() string {
	s := "Validation errors or warnings have occurred:\n"
	for _, err := range ve.Errors {
		s = s + "\t" + err.ErrorType.String() + ": " + err.Message + " @" + err.Location.Path + "\n\tin: " + err.Location.Descriptor + "\n\t"
	}
	return s
}

// JSonContent returns the serialized content of all validations
// errors as JSON
func (ve ValidationErrors) JSonContent() (b []byte, e error) {
	b, e = json.MarshalIndent(ve.Errors, "", "    ")
	return
}

const (
	//Warning allows to mark validation error as Warning
	Warning ErrorType = 0
	//Error allows to mark validation error as Error
	Error ErrorType = 1
)

// String return the name of the given ErrorType
func (r ErrorType) String() string {
	names := [...]string{
		"Warning",
		"Error"}
	if r < Warning || r > Error {
		return "Unknown"
	}
	return names[r]
}

func (ve *ValidationErrors) merge(other ValidationErrors) {
	ve.Errors = append(ve.Errors, other.Errors...)
}

func (ve *ValidationErrors) append(t ErrorType, e string, l DescriptorLocation) {
	ve.Errors = append(ve.Errors, ValidationError{
		Location:  l,
		Message:   e,
		ErrorType: t,
	})
}

func (ve *ValidationErrors) contains(ty ErrorType, m string, path string) bool {
	for _, v := range ve.Errors {
		if v.ErrorType.String() == ty.String() && v.Message == m && v.Location.Path == path {
			return true
		}
	}
	return false
}

func (ve *ValidationErrors) locate(m string) []ValidationError {
	result := make([]ValidationError, 0)
	for _, v := range ve.Errors {
		if v.Message == m {
			result = append(result, v)
		}
	}
	return result
}

func (ve *ValidationErrors) addError(err error, location DescriptorLocation) {
	ve.append(Error, err.Error(), location)
}

func (ve *ValidationErrors) addWarning(message string, location DescriptorLocation) {
	ve.append(Warning, message, location)
}

// HasErrors returns true if the ValidationErrors contains at least one error
func (ve ValidationErrors) HasErrors() bool {
	for _, v := range ve.Errors {
		if v.ErrorType == Error {
			return true
		}
	}
	return false
}

// HasWarnings returns true if the ValidationErrors contains at least one warning
func (ve ValidationErrors) HasWarnings() bool {
	for _, v := range ve.Errors {
		if v.ErrorType == Warning {
			return true
		}
	}
	return false
}

func validate(e Environment, loc DescriptorLocation, ins ...interface{}) ValidationErrors {
	vErrs := ValidationErrors{}
	for _, in := range ins {
		validValue(e, loc, in, &vErrs)
		vOf := reflect.ValueOf(in)
		switch vOf.Kind() {
		case reflect.Map:
			validMap(e, loc, vOf, &vErrs)
		case reflect.Slice:
			validSlice(e, loc, vOf, &vErrs)
		}
	}
	return vErrs
}

func validMap(e Environment, loc DescriptorLocation, val reflect.Value, vErrs *ValidationErrors) {
	for _, key := range val.MapKeys() {
		val := val.MapIndex(key)
		okImpl := reflect.TypeOf(val.Interface()).Implements(reflect.TypeOf((*validatable)(nil)).Elem())
		if okImpl {
			concreteVal, ok := val.Interface().(validatable)
			if ok {
				vErrs.merge(concreteVal.validate(e, loc.appendPath(key.String())))
			}
		}
	}
}

func validSlice(e Environment, loc DescriptorLocation, val reflect.Value, vErrs *ValidationErrors) {
	for i := 0; i < val.Len(); i++ {
		val := val.Index(i)
		okImpl := reflect.TypeOf(val.Interface()).Implements(reflect.TypeOf((*validatable)(nil)).Elem())
		if okImpl {
			concreteVal, ok := val.Interface().(validatable)
			if ok {
				vErrs.merge(concreteVal.validate(e, loc.appendIndex(i)))
			}
		}
	}
}

func validValue(e Environment, loc DescriptorLocation, val interface{}, vErrs *ValidationErrors) {
	concreteVal, ok := val.(validatable)
	if ok {
		vErrs.merge(concreteVal.validate(e, loc))
	}
}
