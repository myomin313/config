package utils

import (
	"log"
	"math"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

/*
Validate is the global validator instance used across the application.
Initialize it once and reuse it everywhere.
*/
var Validate *validator.Validate

/*
Required locales for game translations.
Example expected input:

	{
		"en": "Dragon",
		"ch": "龙",
		"kh": "នាគ"
	}
*/
var requiredGameLocales = []string{
	"en",
	"ch",
	"kh",
}

/*
Map version for quick lookup when validating keys.
*/
var requiredGameLocalesMap = map[string]bool{
	"en": true,
	"ch": true,
	"kh": true,
}

/*
Precompiled regex for strict URL validation.

Allowed formats:
https://example.com
http://example.com
https://sub.domain.com
*/
var strictURLRegex = regexp.MustCompile(`^https?:\/\/([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)

/*
RequireAllGameLocales validates that:

1. Exactly 3 locales exist
2. Keys must be: en, ch, kh
3. Values must not be empty
4. No extra keys allowed
*/
func RequireAllGameLocales(fl validator.FieldLevel) bool {

	names, ok := fl.Field().Interface().(map[string]string)
	if !ok {
		return false
	}

	// Must contain exactly required locales
	if len(names) != len(requiredGameLocales) {
		return false
	}

	// Validate required locales exist and values are non-empty
	for _, locale := range requiredGameLocales {

		value, exists := names[locale]

		if !exists {
			return false
		}

		if strings.TrimSpace(value) == "" {
			return false
		}
	}

	// Ensure no extra locale keys exist
	for key := range names {
		if !requiredGameLocalesMap[key] {
			return false
		}
	}

	return true
}

/*
init initializes the validator and registers custom validations.
*/
func init() {

	Validate = validator.New()

	/*
		notblank validation

		Ensures a string is not empty or whitespace.

		Supports:
		- string
		- *string (optional JSON fields)
	*/
	err := Validate.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {

		field := fl.Field()

		switch field.Kind() {

		case reflect.Ptr:

			// Pointer not provided in JSON
			if field.IsNil() {
				return true
			}

			elem := field.Elem()

			// Ensure pointer is string
			if elem.Kind() != reflect.String {
				return false
			}

			return strings.TrimSpace(elem.String()) != ""

		case reflect.String:

			return strings.TrimSpace(field.String()) != ""

		default:
			return true
		}
	})
	if err != nil {
		log.Fatalf("failed to register validation 'notblank': %v", err)
	}

	/*
		strict_url validation

		Allows only clean domain URLs such as:

		http://example.com
		https://example.com
		https://domain.co
	*/

	if err := Validate.RegisterValidation("strict_url", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		if field.Kind() != reflect.String {
			return false
		}

		return strictURLRegex.MatchString(field.String())
	}); err != nil {
		log.Fatalf("failed to register validation 'strict_url': %v", err)
	}
	/*
		require_all_game_locales validation

		Used for validating translated fields like game names.
	*/
	if err := Validate.RegisterValidation("require_all_game_locales", RequireAllGameLocales); err != nil {
		log.Fatalf("failed to register validation 'require_all_game_locales': %v", err)
	}

	if err := Validate.RegisterValidation("nohyphen", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		if field.Kind() != reflect.String {
			return false
		}

		return !strings.Contains(field.String(), "-")

	}); err != nil {
		log.Fatalf("failed to register validation 'nohyphen': %v", err)
	}
	var tableNumberRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{2,20}$`)

	if err := Validate.RegisterValidation("table_number", func(fl validator.FieldLevel) bool {

		field := fl.Field()

		if field.Kind() != reflect.String {
			return false
		}

		value := strings.TrimSpace(field.String())

		return tableNumberRegex.MatchString(value)

	}); err != nil {
		log.Fatalf("failed to register validation 'table_number': %v", err)
	}

	if err := Validate.RegisterValidation("max3dp", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(float64)
		if !ok {
			return false
		}

		// Check at most 3 decimal places
		return math.Round(value*1000) == value*1000
	}); err != nil {
		log.Fatalf("failed to register validation 'max3dp': %v", err)
	}

}
