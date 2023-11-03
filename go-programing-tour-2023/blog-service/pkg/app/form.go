package app

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
)

// 对入参校验的方法进行二次封装

type ValidError struct {
	Key     string
	Message string
}

type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

func BindAndValid(ctx *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := ctx.ShouldBind(v)
	if err != nil {
		v := ctx.Value("trans")
		trans, _ := v.(ut.Translator)
		var verrs val.ValidationErrors
		ok := errors.As(err, &verrs)
		if !ok {
			return false, errs
		}
		for vek, vev := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key:     vek,
				Message: vev,
			})
		}
		return false, errs
	}
	return true, nil
}
