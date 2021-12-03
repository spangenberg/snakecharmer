package snakecharmer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

type Config interface {
	PreValidate(validate *validator.Validate)
}

func Validate(cfg Config) error {
	v := validator.New()
	trans, err := registerTranslations(cfg, v)
	if err != nil {
		return err
	}
	cfg.PreValidate(v)
	if err = v.Struct(cfg); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		var fieldErrors []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, "  "+fieldError.Translate(trans))
		}
		return fmt.Errorf("\n%s\n", strings.Join(unique(fieldErrors), "\n"))
	}
	return nil
}

type translation struct {
	tag             string
	translation     string
	override        bool
	customRegisFunc validator.RegisterTranslationsFunc
	customTransFunc validator.TranslationFunc
}

func registerTranslations(cfg Config, v *validator.Validate) (_ ut.Translator, err error) {
	trans, _ := ut.New(en.New()).GetTranslator("en")
	if err = entranslations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, err
	}
	translations := []translation{
		{
			tag:         "required",
			translation: "{0} is required",
			override:    true,
		},
		{
			tag:         "hostname_port",
			translation: "{0} must to be a valid HostPort",
			override:    true,
		},
	}
	for _, t := range translations {
		if t.customTransFunc != nil && t.customRegisFunc != nil {
			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc)
		} else if t.customTransFunc != nil && t.customRegisFunc == nil {
			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), t.customTransFunc)
		} else if t.customTransFunc == nil && t.customRegisFunc != nil {
			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, translateFuncFactory(cfg))
		} else {
			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFuncFactory(cfg))
		}
		if err != nil {
			return nil, err
		}
	}
	return trans, nil
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		if err := ut.Add(tag, translation, override); err != nil {
			return err
		}
		return nil
	}
}

func translateFuncFactory(cfg Config) func(ut ut.Translator, fe validator.FieldError) string {
	return func(ut ut.Translator, fe validator.FieldError) string {
		paths := strings.Split(fe.StructNamespace(), ".")
		field := "--" + tagInformation(reflect.TypeOf(cfg).Elem(), paths[1:], "")
		t, err := ut.T(fe.Tag(), field)
		if err != nil {
			panic(fmt.Sprintf("warning: error translating FieldError: %#v", fe))
			return fe.(error).Error()
		}
		return t
	}
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
