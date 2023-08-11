package sdblueprint

import (
	"github.com/gaorx/stardust5/sdcodegen"
	"github.com/gaorx/stardust5/sderr"
)

type Generator interface {
	GenerateTo(buffs *sdcodegen.Buffers, bp *Blueprint) error
}

func (bp *Blueprint) GenerateTo(buffs *sdcodegen.Buffers, generators ...Generator) error {
	for _, g := range generators {
		if g != nil {
			if err := g.GenerateTo(buffs, bp); err != nil {
				return sderr.WithStack(err)
			}
		}
	}
	return nil
}

func (bp *Blueprint) Generate(generators ...Generator) (*sdcodegen.Buffers, error) {
	buffs := sdcodegen.NewBuffers()
	if err := bp.GenerateTo(buffs, generators...); err != nil {
		return nil, sderr.WithStack(err)
	}
	return buffs, nil
}
