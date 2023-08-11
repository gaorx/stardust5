package sdcasbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
)

type StaticRbac struct {
	enforcer *casbin.Enforcer
	polices  string
}

var _ Rbac = &StaticRbac{}

func NewStaticRbac(loaders []RbacLoader, enforcerBuilderOpts RbacEnforcerBuilderOptions) (*StaticRbac, error) {
	b := NewRbacEnforcerBuilder(enforcerBuilderOpts)
	if len(loaders) > 0 {
		if ok := lo.Try0(func() {
			rbacLoad(b, loaders)
		}); !ok {
			return nil, sderr.New("casbin load rbac error")
		}
	}
	e, err := b.Build()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return &StaticRbac{
		enforcer: e,
		polices:  b.GeneratePolicies(),
	}, nil
}

func NewStaticRbacWithFunc(f func(b *RbacEnforcerBuilder), enforcerBuilderOpts RbacEnforcerBuilderOptions) (*StaticRbac, error) {
	return NewStaticRbac([]RbacLoader{RbacLoaderFunc(f)}, enforcerBuilderOpts)
}

func (rbac *StaticRbac) IsGranted(userOrRoleId, objId, action string) bool {
	ok, err := rbac.enforcer.Enforce(userOrRoleId, objId, action)
	if err != nil {
		panic(sderr.Wrap(err, "casbin enforce error"))
	}
	return ok
}
