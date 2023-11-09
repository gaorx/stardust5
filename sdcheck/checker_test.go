package sdcheck

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestAll(t *testing.T) {
	// func
	assert.NoError(t, Func(nil).Check())
	assert.NoError(t, Func(func() error {
		return nil
	}).Check())
	assert.Error(t, Func(func() error {
		return sderr.New("FUNC")
	}).Check())

	// true
	assert.NoError(t, True(true, "TRUE").Check())
	assert.Error(t, True(false, "TRUE").Check())

	// false
	assert.NoError(t, False(false, "FALSE").Check())
	assert.Error(t, False(true, "FALSE").Check())

	// not
	assert.Error(t, Not(True(true, "TRUE"), "NOT").Check())
	assert.NoError(t, Not(True(false, "TRUE"), "NOT").Check())

	// all
	assert.NoError(t, All().Check())
	assert.NoError(t, All(
		True(true, "TRUE"),
		True(true, "TRUE"),
	).Check())
	assert.Error(t, All(
		True(false, "TRUE"),
		True(true, "TRUE"),
	).Check())
	assert.Error(t, All(
		True(true, "TRUE"),
		True(false, "TRUE"),
	).Check())

	// And
	assert.NoError(t, And(nil, "AND").Check())
	assert.NoError(t, And([]Checker{
		True(true, "TRUE"),
	}, "AND").Check(), "AND")
	assert.Error(t, And([]Checker{
		True(false, "TRUE"),
	}, "AND").Check(), "AND")
	assert.NoError(t, And([]Checker{
		True(true, "TRUE"),
		True(true, "TRUE"),
	}, "AND").Check())
	assert.Error(t, And([]Checker{
		True(false, "TRUE"),
		True(true, "TRUE"),
	}, "AND").Check())
	assert.Error(t, And([]Checker{
		True(true, "TRUE"),
		True(false, "TRUE"),
	}, "AND").Check())

	// Or
	assert.NoError(t, Or(nil, "OR").Check())
	assert.NoError(t, Or([]Checker{
		True(true, "TRUE"),
	}, "OR").Check())
	assert.Error(t, Or([]Checker{
		True(false, "TRUE"),
	}, "OR").Check())
	assert.NoError(t, Or([]Checker{
		True(true, "TRUE"),
		True(true, "TRUE"),
	}, "OR").Check())
	assert.NoError(t, Or([]Checker{
		True(false, "TRUE"),
		True(true, "TRUE"),
	}, "AND").Check())
	assert.NoError(t, Or([]Checker{
		True(true, "TRUE"),
		True(false, "TRUE"),
	}, "OR").Check())
	assert.Error(t, Or([]Checker{
		True(false, "TRUE"),
		True(false, "TRUE"),
	}, "OR").Check())

	// if
	assert.NoError(t, If(true, True(true, "TRUE")).Check())
	assert.NoError(t, If(false, True(true, "TRUE")).Check())
	assert.Error(t, If(true, True(false, "TRUE")).Check())
	assert.NoError(t, If(false, True(false, "TRUE")).Check())

	// For
	var a, b int
	a, b = 0, 0
	assert.NoError(t, All(
		For(func() (int, error) { return 3, nil }, &a),
		For(func() (int, error) { return 4, nil }, &b),
	).Check())
	assert.Equal(t, 3, a)
	assert.Equal(t, 4, b)
	a, b = 0, 0
	assert.Error(t, All(
		For(func() (int, error) { return 3, sderr.New("xx") }, &a),
		For(func() (int, error) { return 4, nil }, &b),
	).Check())
	assert.Equal(t, 0, a)
	assert.Equal(t, 0, b)

	// lazy
	a, b = 0, 0
	assert.NoError(t, All(
		For(func() (int, error) { return 3, nil }, &a),
		Lazy(func() Checker {
			// 在lazy中可以使用被上一个checker修改过的a，不放在lazy中则不行
			return For(func() (int, error) { return 7 + a, nil }, &b)
		}),
	).Check())
	assert.Equal(t, 3, a)
	assert.Equal(t, 10, b)
	a, b = 0, 0
	assert.Error(t, All(
		For(func() (int, error) { return 3, nil }, &a),
		Lazy(func() Checker {
			return For(func() (int, error) { return 7 + a, sderr.New("XX") }, &b)
		}),
	).Check())
	assert.Equal(t, a, 3)
	assert.Equal(t, b, 0)

	// required
	assert.Error(t, Required(nil, "REQUIRED").Check())
	assert.Error(t, Required((func() int)(nil), "REQUIRED").Check())
	assert.Error(t, Required((*int)(nil), "REQUIRED").Check())
	assert.NoError(t, Required(1, "REQUIRED").Check())
	assert.Error(t, Required(0, "REQUIRED").Check())
	assert.NoError(t, Required(true, "REQUIRED").Check())
	assert.Error(t, Required(false, "REQUIRED").Check())
	assert.NoError(t, Required("a", "REQUIRED").Check())
	assert.Error(t, Required("", "REQUIRED").Check())
	assert.Error(t, Required([]int{}, "REQUIRED").Check())
	assert.Error(t, Required([]int(nil), "REQUIRED").Check())
	assert.NoError(t, Required([]int{0}, "REQUIRED").Check())
	assert.Error(t, Required(map[string]int{}, "REQUIRED").Check())
	assert.Error(t, Required(map[string]int(nil), "REQUIRED").Check())
	assert.NoError(t, Required(map[string]int{"": 0}, "REQUIRED").Check())
	assert.Error(t, Required(struct{}{}, "REQUIRED").Check())
	assert.NoError(t, Required(&struct{}{}, "REQUIRED").Check())

	// len
	assert.NoError(t, Len([]int{}, 0, 2, "LEN").Check())
	assert.NoError(t, Len([]int{0}, 0, 2, "LEN").Check())
	assert.NoError(t, Len([]int{0, 0}, 0, 2, "LEN").Check())
	assert.Error(t, Len([]int{0, 0, 0}, 0, 2, "LEN").Check())
	assert.NoError(t, Len("", 0, 2, "LEN").Check())
	assert.NoError(t, Len("a", 0, 2, "LEN").Check())
	assert.NoError(t, Len("aa", 0, 2, "LEN").Check())
	assert.Error(t, Len("aaa", 0, 2, "LEN").Check())

	// collection
	assert.Error(t, In("a", []string{}, "IN").Check())
	assert.NoError(t, NotIn("a", []string{}, "NOT_IN").Check())
	assert.NoError(t, In("a", []string{"b", "a"}, "IN").Check())
	assert.Error(t, In("a", []string{"b", "c"}, "IN").Check())
	assert.Error(t, NotIn("a", []string{"b", "a"}, "NOT_IN").Check())
	assert.NoError(t, NotIn("a", []string{"b", "c"}, "NOT_IN").Check())
	assert.NoError(t, HasKey("a", map[string]int{"a": 0}, "HAS_KEY").Check())
	assert.Error(t, HasKey("a", map[string]int{"b": 0}, "HAS_KEY").Check())
	assert.Error(t, NotHasKey("a", map[string]int{"a": 0}, "NOT_HAS_KEY").Check())
	assert.NoError(t, NotHasKey("a", map[string]int{"b": 0}, "NOT_HAS_KEY").Check())

	// string
	assert.NoError(t, MatchRegexp("abc", "a[bd]c", "MATCH_REGEXP").Check())
	assert.NoError(t, MatchRegexp("adc", "a[bd]c", "MATCH_REGEXP").Check())
	assert.Error(t, MatchRegexp("aec", "a[bd]c", "MATCH_REGEXP").Check())
	assert.NoError(t, MatchRegexpPattern("abc", regexp.MustCompile("a([bd])c"), "MATCH_REGEXP_PATTERN").Check())
	assert.NoError(t, MatchRegexpPattern("adc", regexp.MustCompile("a([bd])c"), "MATCH_REGEXP_PATTERN").Check())
	assert.Error(t, MatchRegexpPattern("aec", regexp.MustCompile("a([bd])c"), "MATCH_REGEXP_PATTERN").Check())
	assert.NoError(t, HasSub("abc", "b", "HAS_SUB").Check())
	assert.Error(t, HasSub("adc", "b", "HAS_SUB").Check())
	assert.NoError(t, HasPrefix("ab", "a", "HAS_PREFIX").Check())
	assert.Error(t, HasPrefix("cb", "a", "HAS_PREFIX").Check())
	assert.NoError(t, HasSuffix("ba", "a", "HAS_PREFIX").Check())
	assert.Error(t, HasSuffix("bc", "a", "HAS_SUFFIX").Check())
}
