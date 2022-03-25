package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"

	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	// en_translations "github.com/go-playground/validator/v10/translations/en"
)

// ================================================================================

// validate holds the settings and caches for validating request struct values.
var validate = validator.New()

// translator is a cache of locale and translation information.
var translator *ut.UniversalTranslator

// ================================================================================
func init() {
	// Instantiate the english locale for the validator library.
	enLocale := en.New()

	// Create a value using English as the fallback locale (first argument).
	// Provide one or more arguments for additional supported locales.
	translator = ut.New(enLocale, enLocale)

	// Register the english error messages for validation errors.
	lang, _ := translator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, lang)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

}

// ================================================================================
// Decode
func Decode(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	// ==============================

	if err := validate.Struct(v); err != nil {

		verrs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		lang, _ := translator.GetTranslator("en")

		var fields []FieldError
		for _, verr := range verrs {
			field := FieldError{
				Error: verr.Field(),
				Field: verr.Translate(lang),
			}
			fields = append(fields, field)
		}

		return &Error{
			Err:    errors.New("field validation error"),
			Status: http.StatusBadRequest,
			Fields: fields,
		}

	}

	return nil
}
