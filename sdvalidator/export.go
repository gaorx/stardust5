package sdvalidator

import (
	"context"
	"github.com/gaorx/stardust5/sderr"

	"github.com/go-playground/validator/v10"
)

var defaultValidate = validator.New()

func Default() *validator.Validate {
	return defaultValidate
}

func SetDefault(v *validator.Validate) {
	defaultValidate = v
}

func FieldsErrors(err error) []validator.FieldError {
	if err == nil {
		return nil
	}
	if fieldErrs, ok := sderr.AsT[validator.ValidationErrors](sderr.Cause(err)); ok {
		return fieldErrs
	} else {
		return nil
	}
}

func Struct(s any) error {
	return defaultValidate.Struct(s)
}

func StructCtx(ctx context.Context, s any) error {
	return defaultValidate.StructCtx(ctx, s)
}

func StructPartial(s any, fields []string) error {
	return defaultValidate.StructPartial(s, fields...)
}

func StructPartialCtx(ctx context.Context, s any, fields []string) error {
	return defaultValidate.StructPartialCtx(ctx, s, fields...)
}

func Var(v any, tag string) error {
	return defaultValidate.Var(v, tag)
}
