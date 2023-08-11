package sdcasbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/gaorx/stardust5/sdconcur"
	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
	"slices"
	"sync"
)

type ReloadableRbac struct {
	mtx                   sync.RWMutex
	enforcer              *casbin.Enforcer
	polices               string
	loaders               []RbacLoader
	loaderFactories       []RbacLoaderFactory
	enforceBuilderOptions RbacEnforcerBuilderOptions
}

var _ Rbac = &ReloadableRbac{}

func NewReloadableRbac(
	loaders []RbacLoader,
	loadFactories []RbacLoaderFactory,
	enforcerBuilderOpts RbacEnforcerBuilderOptions,
) *ReloadableRbac {
	return &ReloadableRbac{
		enforcer:              lo.Must(NewRbacEnforcerBuilder(enforcerBuilderOpts).Build()),
		polices:               "",
		loaders:               loaders,
		loaderFactories:       loadFactories,
		enforceBuilderOptions: enforcerBuilderOpts,
	}
}

func (rbac *ReloadableRbac) Reload() error {
	b := NewRbacEnforcerBuilder(rbac.enforceBuilderOptions)
	loaders := rbac.getLoaders()
	if len(loaders) > 0 {
		if ok := lo.Try0(func() {
			rbacLoad(b, loaders)
		}); !ok {
			return sderr.New("casbin reload rbac error")
		}
	}
	e, err := b.Build()
	if err != nil {
		return sderr.WithStack(err)
	}
	sdconcur.LockW(&rbac.mtx, func() {
		rbac.enforcer = e
		rbac.polices = b.GeneratePolicies()
	})
	return nil
}

func (rbac *ReloadableRbac) IsGranted(userOrRoleId, objId, action string) bool {
	var e *casbin.Enforcer
	sdconcur.LockR(&rbac.mtx, func() {
		e = rbac.enforcer
	})
	if e == nil {
		panic(sderr.New("static rbac is not ready"))
	}
	ok, err := e.Enforce(userOrRoleId, objId, action)
	if err != nil {
		panic(sderr.Wrap(err, "casbin enforce error"))
	}
	return ok
}

func (rbac *ReloadableRbac) Policies() string {
	var polices string
	sdconcur.LockR(&rbac.mtx, func() {
		polices = rbac.polices
	})
	return polices
}

func (rbac *ReloadableRbac) getLoaders() []RbacLoader {
	loaders := slices.Clone(rbac.loaders)
	for _, loaderFactory := range rbac.loaderFactories {
		if loaderFactory == nil {
			continue
		}
		loader := loaderFactory.NewRbacLoader()
		if loader == nil {
			continue
		}
		loaders = append(loaders, loader)
	}
	return loaders
}
