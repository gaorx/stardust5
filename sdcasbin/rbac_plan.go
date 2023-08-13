package sdcasbin

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

type RbacPlan struct {
	SuperRole string
	Users     []User
	Objects1  []Object
	Objects2  []Object
	Objects3  []Object
	Grants1   []Grant
	Grants2   []Grant
	Grants3   []Grant
}

type User struct {
	Id      string
	RoleIds []string
}

type ObjectSet struct {
	Ids              []string
	AvailableActions []string
}

type Grant struct {
	RoleId   string
	ObjectId string
	Actions  []string
}

func U(uid string, roleIds ...string) User {
	return User{Id: uid, RoleIds: roleIds}
}

func O(oid string, availableActions ...string) Object {
	return Object{Id: oid, AvailableActions: availableActions}
}

func G(roleId string) Grant {
	return Grant{RoleId: roleId}
}

func (g Grant) To(oid string) Grant {
	g.ObjectId = oid
	return g
}

func (g Grant) On(actions ...string) Grant {
	g.Actions = lo.Uniq(append(g.Actions, actions...))
	return g
}

func (plan RbacPlan) Apply(b *RbacEnforcerBuilder) error {
	uidSet, roleIdSet := map[string]int{}, map[string]int{}
	for _, u := range plan.Users {
		if u.Id == "" {
			continue
		}
		uidSet[u.Id] = 0
		for _, roleId := range u.RoleIds {
			if roleId != "" {
				roleIdSet[roleId] = 0
			}
		}
	}

	// add users
	b.AddUserIds(maps.Keys(uidSet)...)

	// add roles
	b.AddRoleIds(maps.Keys(roleIdSet)...)

	// add objects
	b.AddObjects(plan.Objects1...)
	b.AddObjects(plan.Objects2...)
	b.AddObjects(plan.Objects3...)

	// add user <-> role links
	for _, u := range plan.Users {
		if u.Id != "" && len(u.RoleIds) > 0 {
			if err := b.AddUserRoleLink([]string{u.Id}, u.RoleIds); err != nil {
				return sderr.WithStack(err)
			}
		}
	}

	// add role <-> object links
	var grants []Grant
	grants = append(grants, plan.Grants1...)
	grants = append(grants, plan.Grants2...)
	grants = append(grants, plan.Grants3...)
	for _, g := range grants {
		if g.RoleId != "" && g.ObjectId != "" {
			actions := g.Actions
			if len(actions) <= 0 {
				actions = []string{"*"}
			}
			if err := b.AddObjectRoleLink(g.ObjectId, []string{g.RoleId}, actions); err != nil {
				return sderr.WithStack(err)
			}
		}
	}
	return nil
}

func (plan RbacPlan) Prepare(b *RbacEnforcerBuilder) {}
func (plan RbacPlan) Load(b *RbacEnforcerBuilder) {
	lo.Must0(plan.Apply(b))
}

func (plan RbacPlan) ToRbac() (Rbac, error) {
	staticRbac, err := NewStaticRbac([]RbacLoader{plan}, &RbacEnforcerBuilderOptions{
		SuperuserId: plan.SuperRole,
	})
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return staticRbac, nil
}

func (plan RbacPlan) MustToRbac() Rbac {
	return lo.Must(plan.ToRbac())
}
