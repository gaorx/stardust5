package sdcasbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	stringadapter "github.com/casbin/casbin/v2/persist/string-adapter"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/samber/lo"
	"strings"
)

type RbacEnforcerBuilder struct {
	options   RbacEnforcerBuilderOptions
	userIds   *idSet
	roleIds   *idSet
	objectIds *idSet
	urEntries []builderEntry
	orEntries []builderEntry
	rrEntries []builderEntry
}

type RbacEnforcerBuilderOptions struct {
	DisableCheck        bool
	SuperuserId         string
	AllowObjectUserLink bool
}

type Object struct {
	Id               string
	AvailableActions []string
}

type builderEntry struct {
	v1, v2, v3 string
}

func NewRbacEnforcerBuilder(opts *RbacEnforcerBuilderOptions) *RbacEnforcerBuilder {
	opts1 := lo.FromPtr(opts)
	return &RbacEnforcerBuilder{
		options:   opts1,
		userIds:   newIdSet(),
		roleIds:   newIdSet(),
		objectIds: newIdSet(),
	}
}

func (b *RbacEnforcerBuilder) Build() (*casbin.Enforcer, error) {
	var m model.Model
	if b.options.SuperuserId != "" {
		m = rbacWithSuperuserModel(b.options.SuperuserId)
	} else {
		m = rbacModel()
	}
	policies := b.GeneratePolicies()
	if strings.TrimSpace(policies) != "" {
		//fmt.Println("---------")
		//fmt.Println(policies)
		//fmt.Println("---------")
		enforcer, err := casbin.NewEnforcer(m, stringadapter.NewAdapter(policies))
		if err != nil {
			return nil, sderr.WithStack(err)
		}
		return enforcer, nil
	} else {
		return casbin.NewEnforcer(m)
	}
}

func (b *RbacEnforcerBuilder) GeneratePolicies() string {
	q := func(v string) string {
		needQuote := false
		if strings.Contains(v, ",") {
			needQuote = true
		}
		if strings.Contains(v, "\"") {
			needQuote = true
			v = strings.Replace(v, "\"", "\"\"", -1)
		}
		if needQuote {
			return "\"" + v + "\""
		} else {
			return v
		}
	}
	var buff strings.Builder
	for _, entry := range b.orEntries {
		_, _ = fmt.Fprintf(&buff, "p,%s,%s,%s\n", q(entry.v1), q(entry.v2), q(entry.v3))
	}
	for _, entry := range b.rrEntries {
		_, _ = fmt.Fprintf(&buff, "g,%s,%s\n", q(entry.v1), q(entry.v2))
	}
	for _, entry := range b.urEntries {
		_, _ = fmt.Fprintf(&buff, "g,%s,%s\n", q(entry.v1), q(entry.v2))
	}
	return buff.String()
}

func (b *RbacEnforcerBuilder) Clear() *RbacEnforcerBuilder {
	b.userIds.clear()
	b.roleIds.clear()
	b.objectIds.clear()
	b.urEntries = nil
	b.orEntries = nil
	b.rrEntries = nil
	return b
}

func (b *RbacEnforcerBuilder) AddUserIds(uids ...string) *RbacEnforcerBuilder {
	b.userIds.add(uids)
	return b
}

func (b *RbacEnforcerBuilder) AddRoleIds(roleIds ...string) *RbacEnforcerBuilder {
	b.roleIds.add(roleIds)
	return b
}

func (b *RbacEnforcerBuilder) AddObjects(objs ...Object) *RbacEnforcerBuilder {
	for _, obj := range objs {
		b.objectIds.addWith(obj.Id, obj.AvailableActions)
	}
	return b
}

func (b *RbacEnforcerBuilder) AddObjectIds(objIds []string, availableActions []string) *RbacEnforcerBuilder {
	for _, objId := range objIds {
		b.objectIds.addWith(objId, availableActions)
	}
	return b
}

func (b *RbacEnforcerBuilder) AddUserRoleLink(userIds []string, roleIds []string) error {
	if len(userIds) <= 0 || len(roleIds) <= 0 {
		return nil
	}
	for _, userId := range userIds {
		if err := b.checkUserId(userId); err != nil {
			return sderr.WithStack(err)
		}
	}
	for _, roleId := range roleIds {
		if err := b.checkRoleId(roleId); err != nil {
			return sderr.WithStack(err)
		}
	}

	for _, userId := range userIds {
		for _, roleId := range roleIds {
			b.urEntries = append(b.urEntries, builderEntry{v1: userId, v2: roleId})
		}
	}
	return nil
}

func (b *RbacEnforcerBuilder) AddObjectRoleLink(objId string, roleOrUserIds []string, actions []string) error {
	if err := b.checkObjectId(objId); err != nil {
		return sderr.WithStack(err)
	}
	if b.options.AllowObjectUserLink {
		for _, roleId := range roleOrUserIds {
			if err := b.checkRoleOrUserId(roleId); err != nil {
				return sderr.WithStack(err)
			}
		}
	} else {
		for _, roleId := range roleOrUserIds {
			if err := b.checkRoleId(roleId); err != nil {
				return sderr.WithStack(err)
			}
		}
	}

	availableActions := b.objectIds.getVals(objId)
	actions1 := expandActions(actions, availableActions, nil)
	for _, action := range actions1 {
		if err := b.checkAction(objId, action); err != nil {
			return sderr.WithStack(err)
		}
	}
	for _, roleOrUserId := range roleOrUserIds {
		for _, action := range actions1 {
			b.orEntries = append(b.orEntries, builderEntry{
				v1: roleOrUserId,
				v2: objId,
				v3: action,
			})
		}
	}
	return nil
}

func (b *RbacEnforcerBuilder) AddRoleRoleLink(childRoleId string, parentRoleIds []string) error {
	if len(parentRoleIds) <= 0 {
		return nil
	}

	if err := b.checkRoleId(childRoleId); err != nil {
		return sderr.WithStack(err)
	}
	for _, parentRoleId := range parentRoleIds {
		if err := b.checkRoleId(parentRoleId); err != nil {
			return sderr.WithStack(err)
		}
	}
	for _, parentRoleId := range parentRoleIds {
		b.rrEntries = append(b.rrEntries, builderEntry{
			v1: childRoleId,
			v2: parentRoleId,
		})
	}
	return nil
}

func (b *RbacEnforcerBuilder) checkUserId(userId string) error {
	if !b.options.DisableCheck && !b.userIds.has(userId) {
		return sderr.WrapWith(ErrIllegalUserId, "", sderr.Attrs{"user_id": userId})
	} else {
		return nil
	}
}

func (b *RbacEnforcerBuilder) checkRoleId(roleId string) error {
	if b.options.DisableCheck {
		return nil
	}
	if b.roleIds.has(roleId) {
		return nil
	}
	if b.options.SuperuserId != "" && roleId == b.options.SuperuserId {
		return nil
	}
	return sderr.WrapWith(ErrIllegalRoleId, "", sderr.Attrs{"role_id": roleId})
}

func (b *RbacEnforcerBuilder) checkRoleOrUserId(userOrRoleId string) error {
	if b.options.DisableCheck {
		return nil
	}
	if b.userIds.has(userOrRoleId) || b.roleIds.has(userOrRoleId) {
		return nil
	}
	return sderr.WrapWith(ErrIllegalRoleId, "", sderr.Attrs{"user_id/role_id": userOrRoleId})
}

func (b *RbacEnforcerBuilder) checkObjectId(objId string) error {
	if !b.options.DisableCheck && !b.objectIds.has(objId) {
		return sderr.WrapWith(ErrIllegalObjectId, "", sderr.Attrs{"object_id": objId})
	} else {
		return nil
	}
}

func (b *RbacEnforcerBuilder) checkAction(objId, action string) error {
	if !b.options.DisableCheck && !b.objectIds.hasVal(objId, action) {
		return sderr.WrapWith(ErrIllegalAction, "", sderr.Attrs{"object_id": objId, "action": action})
	} else {
		return nil
	}
}

func rbacModel() model.Model {
	text := sdstrings.TrimMargin(`
		|[request_definition]
		|r = sub, obj, act
		|
		|[policy_definition]
		|p = sub, obj, act
		|
		|[role_definition]
		|g = _, _
		|
		|[policy_effect]
		|e = some(where (p.eft == allow))
		|
		|[matchers]
		|m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
	`, "|")
	m, err := model.NewModelFromString(text)
	if err != nil {
		panic(sderr.New("define rbac model error"))
	}
	return m
}

func rbacWithSuperuserModel(superuserId string) model.Model {
	text := sdstrings.TrimMargin(fmt.Sprintf(`
		|[request_definition]
		|r = sub, obj, act
		|
		|[policy_definition]
		|p = sub, obj, act
		|
		|[role_definition]
		|g = _, _
		|
		|[policy_effect]
		|e = some(where (p.eft == allow))
		|
		|[matchers]
		|m = (g(r.sub, p.sub) || g(r.sub, "%s")) && r.obj == p.obj && r.act == p.act
	`, superuserId), "|")
	m, err := model.NewModelFromString(text)
	if err != nil {
		panic(sderr.New("define rbac(superuser) model error"))
	}
	return m
}
