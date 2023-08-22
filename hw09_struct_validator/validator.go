package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) String() string {
	return fmt.Sprintf("%s:%s", v.Field, v.Err)
}

type ValidationErrors []ValidationError

var (
	errNotStruct      = errors.New("v is not a struct")
	errLessThanMin    = errors.New("value is less than min")
	errGreaterThanMax = errors.New("value is greater than max")
	errNotInSet       = errors.New("value is not in set")
	errLenNotMatch    = errors.New("string length doesn't match")
	errRegexNotMatch  = errors.New("string doesn't match a regex")
)

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}
	for i, err := range v {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(err.String())
	}
	return sb.String()
}

func Validate(v interface{}) error {
	vType := reflect.TypeOf(v)
	vValue := reflect.ValueOf(v)

	if vType.Kind() != reflect.Struct {
		return errNotStruct
	}

	validationErrors := ValidationErrors{}

	for i := 0; i < vType.NumField(); i++ {
		fieldType := vType.Field(i)
		fieldVal := vValue.Field(i)
		constraints := strings.Split(fieldType.Tag.Get("validate"), "|")

		switch fieldType.Type.Kind() { //nolint:exhaustive
		case reflect.Int:
			validationResult, err := validateInt(fieldVal.Interface().(int), fieldType.Name, constraints)
			if err != nil {
				return err
			}
			validationErrors = append(validationErrors, validationResult...)
		case reflect.String:
			validationResult, err := validateString(fieldVal.String(), fieldType.Name, constraints)
			if err != nil {
				return err
			}
			validationErrors = append(validationErrors, validationResult...)
		case reflect.Array:
			validationResult, err := validateArray(fieldType, fieldVal, constraints)
			if err != nil {
				return err
			}
			validationErrors = append(validationErrors, validationResult...)
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateInt(number int, name string, rules []string) (ValidationErrors, error) {
	var result ValidationErrors
	for _, rule := range rules {
		keyValue := strings.Split(rule, ":")

		switch keyValue[0] {
		case "min":
			ruleValue, err := strconv.Atoi(keyValue[1])
			if err != nil {
				return nil, err
			}
			if number < ruleValue {
				result = append(result, ValidationError{Field: name, Err: errLessThanMin})
			}
		case "max":
			ruleValue, err := strconv.Atoi(keyValue[1])
			if err != nil {
				return nil, err
			}
			if number > ruleValue {
				result = append(result, ValidationError{Field: name, Err: errGreaterThanMax})
			}
		case "in":
			set := strings.Split(keyValue[1], ",")
			numInSet := false

			for _, strVal := range set {
				val, err := strconv.Atoi(strVal)
				if err != nil {
					return nil, err
				}
				if number == val {
					numInSet = true
					break
				}
			}
			if !numInSet {
				result = append(result, ValidationError{Field: name, Err: errNotInSet})
			}
		}
	}

	return result, nil
}

func validateString(str string, name string, rules []string) (ValidationErrors, error) {
	var result ValidationErrors
	for _, rule := range rules {
		keyValue := strings.Split(rule, ":")

		switch keyValue[0] {
		case "len":
			strLen, err := strconv.Atoi(keyValue[1])
			if err != nil {
				return nil, err
			}
			if len(str) != strLen {
				result = append(result, ValidationError{Field: name, Err: errLenNotMatch})
			}
		case "regexp":
			regex, err := regexp.Compile(keyValue[1])
			if err != nil {
				return nil, err
			}
			if !regex.MatchString(str) {
				result = append(result, ValidationError{Field: name, Err: errRegexNotMatch})
			}
		case "in":
			set := strings.Split(keyValue[1], ",")
			strInSet := false
			for _, val := range set {
				if str == val {
					strInSet = true
					break
				}
			}
			if !strInSet {
				result = append(result, ValidationError{Field: name, Err: errNotInSet})
			}
		}
	}
	return result, nil
}

func validateArray(fieldType reflect.StructField, fieldVal reflect.Value, constraints []string) (ValidationErrors, error) { //nolint:lll
	validationErrors := ValidationErrors{}
	switch fieldType.Type.Elem().Kind() { //nolint:exhaustive
	case reflect.Int:
		for _, val := range fieldVal.Interface().([]int) {
			validationResult, err := validateInt(val, fieldType.Name, constraints)
			if err != nil {
				return nil, err
			}
			validationErrors = append(validationErrors, validationResult...)
		}
	case reflect.String:
		for _, str := range fieldVal.Interface().([]string) {
			validationResult, err := validateString(str, fieldType.Name, constraints)
			if err != nil {
				return nil, err
			}
			validationErrors = append(validationErrors, validationResult...)
		}
	}
	return validationErrors, nil
}
