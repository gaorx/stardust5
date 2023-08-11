package sdgorm

import (
	"slices"
	"sync"

	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdreflect"
	"gorm.io/gorm/schema"
)

var schemaCacheStore = sync.Map{}

func ParseModel(t any, namer schema.Namer) (*schema.Schema, error) {
	s, err := schema.Parse(t, &schemaCacheStore, namer)
	if err != nil {
		return nil, sderr.Wrap(err, "parse table ")
	}
	return s, nil
}

func ParsePrimaryColumnNames(t any, namer schema.Namer) ([]string, error) {
	s, err := ParseModel(t, namer)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return slices.Clone(s.PrimaryFieldDBNames), nil
}

func GetPrimaryKeys(t any, namer schema.Namer) ([]any, error) {
	s, err := ParseModel(t, namer)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	if len(s.PrimaryFields) <= 0 {
		return nil, sderr.New("no primary keys")
	}
	ids := make([]any, 0, len(s.PrimaryFields))
	for _, f := range s.PrimaryFields {
		id, err := sdreflect.StructGetFieldValue(t, f.Name)
		if err != nil {
			return nil, sderr.Wrap(err, "get primary key field error")
		}
		ids = append(ids, id)
	}
	return ids, nil
}
