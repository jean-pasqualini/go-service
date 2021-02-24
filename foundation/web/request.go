package web

/**
en "github.com/go-playground/locales/en"
ut "github.com/go-playground/universal-translator"
validator "gopkg.in/go-playground/validator.v9"
en_translations "gopkg.in/go-playground/validator/v9/translations/en"
 */

import (
	"encoding/json"
	"github.com/dimfeld/httptreemux/v5"
	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	validator "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	"net/http"
	"reflect"
	"strings"
)

// validte holds the settings and caches for validating request struct value
var validate *validator.Validate

// translator is a cache of locale and translation information
var translator *ut.UniversalTranslator

func init() {
	// Instanciate the validator for use.
	validate = validator.New()

	// Instanciate the english locale for the validator library.
	enLocale := en.New()

	// Create a value using English as the fallback locale (first argument).
	// Provide one or more arguments for additional supported locales.
	translator := ut.New(enLocale, enLocale)

	// Register the english error messages for validation errors.
	lang, _ := translator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, lang)

	// Use JSON tag names for errors instead of Go struct names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Params returns the web call parameters from the request.
func Params(r *http.Request) map[string]string {
	return httptreemux.ContextParams(r.Context())
}

func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(val); err != nil {
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		lang, _ := translator.GetTranslator("en")

		var fields []FieldError
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(lang),
			}
			fields = append(fields, field)
		}

		return &Error{
			Err: errors.New("field validation error"),
			Status: http.StatusBadRequest,
			Fields: fields,
		}
	}

	return nil
}