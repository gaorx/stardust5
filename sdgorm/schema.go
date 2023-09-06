package sdgorm

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
	"gorm.io/gorm/schema"
	"sync"
)

type Schema struct {
	*schema.Schema
}

var schemaCacheStore = sync.Map{}

func ParseSchema(model any, namer schema.Namer) (Schema, error) {
	if namer == nil {
		namer = schema.NamingStrategy{
			TablePrefix:         "",
			SingularTable:       false,
			NameReplacer:        nil,
			NoLowerCase:         false,
			IdentifierMaxLength: 64,
		}
	}
	s, err := schema.Parse(model, &schemaCacheStore, namer)
	if err != nil {
		return Schema{}, sderr.Wrap(err, "parse table ")
	}
	return Schema{Schema: s}, nil
}

func (s Schema) IsNil() bool {
	return s.Schema == nil
}

func (s Schema) DBFieldNames() []string {
	return lo.Map(s.Fields, func(f *schema.Field, _ int) string {
		return f.Name
	})
}
