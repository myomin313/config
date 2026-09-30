package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func GetValidationMessages(err error) map[string]string {
	errors := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		field := e.Field()
		tag := e.Tag()

		switch tag {
		case "required":
			errors[field] = fmt.Sprintf("%s is required", field)
		case "email":
			errors[field] = "Please enter a valid email address"
		case "min":
			errors[field] = fmt.Sprintf("%s must be at least %s characters long", field, e.Param())
		case "max":
			errors[field] = fmt.Sprintf("%s must be at most %s characters long", field, e.Param())
		case "gte":
			errors[field] = fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
		case "lte":
			errors[field] = fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
		default:
			errors[field] = fmt.Sprintf("%s is invalid", field)
		}
	}

	return errors
}

func ValidationErrorMessages(err error, obj interface{}) []string {
	var messages []string

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{err.Error()}
	}

	objType := reflect.TypeOf(obj)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	for _, ve := range validationErrors {

		fieldName := resolveJSONFieldName(obj, ve.StructNamespace())

		// Build human-readable message
		var message string
		switch ve.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required field", fieldName)

		case "strict_url":
			message = fmt.Sprintf("%s must be a valid URL", fieldName)

		case "url":
			message = fmt.Sprintf("%s must be a valid URL", fieldName)
		case "oneof":
			message = fmt.Sprintf(
				"%s must be one of %s",
				fieldName,
				ve.Param(),
			)
		case "gt":
			message = fmt.Sprintf(
				"%s must be greater than %s",
				fieldName,
				ve.Param(),
			)
		case "gte":
			message = fmt.Sprintf(
				"%s must be greater than or equal to %s",
				fieldName,
				ve.Param(),
			)
		case "email":
			message = "Please enter a valid email address"
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters long", fieldName, ve.Param())
		case "max":
			message = fmt.Sprintf("%s must be at most %s characters long", fieldName, ve.Param())
		case "lte":
			message = fmt.Sprintf("%s must be less than or equal to %s", fieldName, ve.Param())
		case "require_all_game_locales":
			message = fmt.Sprintf("%s must contain all required locales", fieldName)
		case "max3dp":
			message = fmt.Sprintf("%s maximum three decimal place is allowed %s", fieldName, ve.Param())
		default:
			message = fmt.Sprintf("%s is invalid", fieldName)
		}

		messages = append(messages, message)
	}

	return messages
}

// helper to find first comma in tag string
func findComma(s string) int {
	for i, ch := range s {
		if ch == ',' {
			return i
		}
	}
	return -1
}

func resolveJSONFieldName(obj interface{}, namespace string) string {
	parts := strings.Split(namespace, ".")

	currentType := reflect.TypeOf(obj)
	if currentType.Kind() == reflect.Ptr {
		currentType = currentType.Elem()
	}

	var jsonPath []string

	for _, part := range parts[1:] { // skip root struct

		var fieldName string
		var indexPart string

		// Handle slice index
		if idx := strings.Index(part, "["); idx != -1 {
			fieldName = part[:idx]
			indexPart = part[idx:] // [1]
		} else {
			fieldName = part
		}

		sf, ok := currentType.FieldByName(fieldName)
		if !ok {
			continue // skip if not found (safe fallback)
		}

		// Extract JSON tag safely
		jsonTag := sf.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			jsonTag = fieldName // fallback to struct field name
		} else {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		// Append with index if exists
		if indexPart != "" {
			jsonPath = append(jsonPath, jsonTag+indexPart)
		} else {
			jsonPath = append(jsonPath, jsonTag)
		}

		// Move to next type safely
		fieldType := sf.Type

		// Handle pointer
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		// Handle slice
		if fieldType.Kind() == reflect.Slice {
			fieldType = fieldType.Elem()
		}

		currentType = fieldType
	}

	return strings.Join(jsonPath, ".")
}
