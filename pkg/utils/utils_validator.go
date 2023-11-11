// Copyright The Jet authors. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utils

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

var validate = validator.New()

func Struct(s interface{}) error {
	return validate.Struct(s)
}

// ProcessErr processes validation errors and returns an error message.
// It handles custom rules and error messages for the go-validator parameter validator.
func ProcessErr(u interface{}, err error) string {
	if err == nil {
		return ""
	}

	if invalid, ok := err.(*validator.InvalidValidationError); ok {
		// If the error is an InvalidValidationError, it means the input parameter is invalid.
		return "Input parameter error: " + invalid.Error()
	}

	validationErrs := err.(validator.ValidationErrors)
	for _, validationErr := range validationErrs {
		// Get the field that doesn't match the format.
		fieldName := validationErr.Field()
		typeOf := reflect.TypeOf(u)

		// If the type is a pointer, get its underlying type.
		if typeOf.Kind() == reflect.Ptr {
			typeOf = typeOf.Elem()
		}

		// Get the field using reflection.
		if field, ok := typeOf.FieldByName(fieldName); ok {
			// Get the reg_error_info tag value for the field.
			errorInfo := field.Tag.Get("reg_error_info")
			if errorInfo == "" {
				// If reg_error_info is empty, fall back to reg_err_info.
				errorInfo = field.Tag.Get("reg_err_info")
			}
			// Return the error message.
			return fieldName + ": " + errorInfo
		} else {
			// The field is missing reg_error_info tag.
			return "Parameter validation failed on the " + validationErr.StructField()
		}
	}

	return ""
}
