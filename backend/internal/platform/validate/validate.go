package validate

import "github.com/go-playground/validator/v10"

type Custom struct{ validator *validator.Validate }

func New() *Custom { return &Custom{validator: validator.New()} }

func (c *Custom) Validate(value any) error { return c.validator.Struct(value) }
