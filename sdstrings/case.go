package sdstrings

import (
	"github.com/fatih/camelcase"
	"github.com/iancoleman/strcase"
)

var (
	ToSnakeL     = strcase.ToSnake
	ToSnakeU     = strcase.ToScreamingSnake
	ToKebabL     = strcase.ToKebab
	ToKebabU     = strcase.ToScreamingKebab
	ToDelimitedL = strcase.ToDelimited
	ToDelimitedU = strcase.ToScreamingDelimited
	ToCamelL     = strcase.ToLowerCamel
	ToCamelU     = strcase.ToCamel
	SplitCamel   = camelcase.Split
)
