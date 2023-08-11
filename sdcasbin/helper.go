package sdcasbin

import (
	"github.com/samber/lo"
	"slices"
)

type idSet struct {
	ids   []string
	idMap map[string][]string
}

func newIdSet() *idSet {
	return &idSet{idMap: map[string][]string{}}
}

func (s *idSet) clear() {
	s.ids, s.idMap = nil, map[string][]string{}
}

func (s *idSet) has(id string) bool {
	_, ok := s.idMap[id]
	return ok
}

func (s *idSet) getVals(id string) []string {
	vals, ok := s.idMap[id]
	if !ok {
		return nil
	}
	return vals
}

func (s *idSet) hasVal(id string, val string) bool {
	vals := s.getVals(id)
	if len(vals) <= 0 {
		return false
	}
	return lo.Contains(vals, val)
}

func (s *idSet) add(ids []string) {
	for _, id := range ids {
		if id == "" {
			continue
		}
		if s.has(id) {
			continue
		}
		s.ids = append(s.ids, id)
		s.idMap[id] = nil
	}
}

func (s *idSet) addWith(id string, vals []string) {
	if id == "" {
		return
	}
	if s.has(id) {
		return
	}
	s.ids = append(s.ids, id)
	s.idMap[id] = slices.Clone(vals)
}

func (s *idSet) remove(ids []string) {
	for _, id := range ids {
		if s.has(id) {
			s.ids = lo.Without(s.ids, id)
			delete(s.idMap, id)
		}
	}
}

var allAlias = []string{"*", "all", "ALL", "All"}

func expandActions(actions []string, all []string, aliases map[string][]string) []string {
	if len(actions) <= 0 {
		return actions
	}
	isAll := false
	for _, action := range actions {
		if lo.Contains(allAlias, action) {
			isAll = true
			break
		}
	}
	if isAll {
		return slices.Clone(all)
	}
	if len(aliases) <= 0 {
		return slices.Clone(actions)
	}
	set := newIdSet()
	for _, action := range actions {
		alias, ok := aliases[action]
		if ok {
			set.add(alias)
		} else {
			set.add([]string{action})
		}
	}
	return set.ids
}
