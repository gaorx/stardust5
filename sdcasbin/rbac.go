package sdcasbin

type Rbac interface {
	IsGranted(userOrRoleId, objId, action string) bool
}

type RbacLoader interface {
	Prepare(b *RbacEnforcerBuilder)
	Load(b *RbacEnforcerBuilder)
}

type RbacLoaderFactory interface {
	NewRbacLoader() RbacLoader
}

type (
	RbacLoaderFunc        func(b *RbacEnforcerBuilder)
	RbacLoaderFactoryFunc func() RbacLoader
)

func (f RbacLoaderFunc) Prepare(b *RbacEnforcerBuilder) {}
func (f RbacLoaderFunc) Load(b *RbacEnforcerBuilder) {
	if f != nil {
		f(b)
	}
}
func (f RbacLoaderFactoryFunc) NewRbacLoader() RbacLoader {
	return f()
}

func rbacLoad(b *RbacEnforcerBuilder, loaders []RbacLoader) {
	for _, loader := range loaders {
		if loader == nil {
			continue
		}
		loader.Prepare(b)
	}
	for _, loader := range loaders {
		if loader == nil {
			continue
		}
		loader.Load(b)
	}
}

func RbacLoaderAsFactory(f func(b *RbacEnforcerBuilder)) RbacLoaderFactory {
	return RbacLoaderFactoryFunc(func() RbacLoader {
		return RbacLoaderFunc(f)
	})
}
