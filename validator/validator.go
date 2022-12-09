package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bngbngstwnd/library-go-chub/constant"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	v     *validator.Validate
	trans ut.Translator
)

func GetValidator() (*validator.Validate, ut.Translator) {
	return v, trans
}

func InitValidator() {
	translator := en.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	var found bool
	trans, found = uni.GetTranslator("en")
	if !found {
		fmt.Println("translator not found")
	}

	v = validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		fmt.Println(err)
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} tidak boleh kosong", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("numeric", trans, func(ut ut.Translator) error {
		return ut.Add("numeric", "{0} harus berupa angka", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("numeric", fe.Field())
		return t
	})

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	_ = v.RegisterValidation("campaigndate", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		_, err := time.Parse(constant.DATE_TIME_LAYOUT_NON_TZ, val)
		return err == nil
	})

	_ = v.RegisterTranslation("campaigndate", trans, func(ut ut.Translator) error {
		return ut.Add("campaigndate", "{0} must be in format 'yyyy-mm-dd hh:mm:ss'", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("campaigndate", fe.Field())
		return t
	})

	_ = v.RegisterValidation("datedob", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		_, err := time.Parse(constant.DATE_LAYOUT, val)
		return err == nil
	})

	_ = v.RegisterTranslation("datedob", trans, func(ut ut.Translator) error {
		return ut.Add("datedob", "{0} must be in format 'yyyy-mm-dd'", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("datedob", fe.Field())
		return t
	})

	_ = v.RegisterValidation("productid", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()

		regex, _ := regexp.Compile("[^-_.a-zA-Z0-9]+")
		result := regex.MatchString(val)
		return !result
	})

	_ = v.RegisterTranslation("productid", trans, func(ut ut.Translator) error {
		return ut.Add("productid", "{0} tidak boleh mengandung karakter spesial", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("productid", fe.Field())
		return t
	})
}

func GetAndValidatePersonalNumber(c *gin.Context) (string, error) {
	personalNumber := c.GetHeader("X-CHUB-PERSONAL-NUMBER")
	if len(personalNumber) != 8 {
		return "", errors.New("invalid personal number")
	}
	if _, err := strconv.ParseInt(personalNumber, 10, 64); err != nil {
		return "", errors.New("personal number should be numeric")
	}
	return personalNumber, nil
}

func GetChannel(c *gin.Context, required bool) (string, error) {
	channel := c.GetHeader("X-CHUB-CHANNEL")
	if required && len(channel) == 0 {
		return "", errors.New("channel is required")
	}
	return channel, nil
}

func GetPersonalNumber(c *gin.Context, required bool) (string, error) {
	pn := c.GetHeader("X-CHUB-PERSONAL-NUMBER")
	if required && len(pn) == 0 {
		return "", errors.New("pn is required")
	}
	if required && len(pn) != 8 {
		return "", errors.New("invalid personal number")
	}
	if _, err := strconv.ParseInt(pn, 10, 64); required && err != nil {
		return "", errors.New("personal number should be numeric")
	}
	return pn, nil
}

func Validate(req interface{}) error {
	v, trans := GetValidator()
	var errorMessages []string

	errData := v.Struct(req)
	if errData != nil {
		for _, e := range errData.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, e.Translate(trans))
		}
		errMessage := strings.Join(errorMessages, "\n")

		return errors.New(errMessage)
	}

	return nil
}
