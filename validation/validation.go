package validation

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"log/slog"
	"reflect"
	"strings"
)

var (
	validate  *validator.Validate
	uni       *ut.UniversalTranslator
	translate ut.Translator
)

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Register "translations" to output validation errors in user-friendly messages
	enLocale := en.New()
	uni = ut.New(enLocale, enLocale)
	translate, _ = uni.GetTranslator("en")
	err := en_translations.RegisterDefaultTranslations(validate, translate)
	panicIf(err)

	// Register structNamespace processor for JSON tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})
}

func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			slog.Warn("Invalid validation", "err", err)
			return errors.New("validation failed")
		}

		errMsgs := make([]string, 0, len(err.(validator.ValidationErrors)))
		for _, e := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, e.Translate(translate))
		}
		return fmt.Errorf("validation error: %s", strings.Join(errMsgs, "; "))
	}
	return nil
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
