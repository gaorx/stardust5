package sdecho

import (
	"context"
	"github.com/gaorx/stardust5/sdslices"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"slices"
)

type Menu struct {
	Id    string      `json:"id"`
	Items []*MenuItem `json:"items,omitempty"`
}

type MenuItem struct {
	Name     string      `json:"name"`
	Path     string      `json:"path,omitempty"`
	Href     string      `json:"href,omitempty"`
	Icon     string      `json:"icon,omitempty"`
	Children []*MenuItem `json:"children,omitempty"`
	Object   Object      `json:"-"`
	Tags     []string    `json:"-"`
}

func MenuReify(ctx context.Context, ec echo.Context, menuId string, tags []string) *Menu {
	menus := MustGet[Menus](ec, keyMenus)
	for _, menu := range menus {
		if menu != nil && menu.Id == menuId {
			token, err := TokenDecode(ctx, ec)
			if err != nil {
				return &Menu{Id: menuId, Items: []*MenuItem{}}
			}
			return menu.Reify(ctx, ec, token, tags)
		}
	}
	return nil
}

func (menu *Menu) Reify(ctx context.Context, ec echo.Context, token Token, tags []string) *Menu {
	return &Menu{
		Id:    menu.Id,
		Items: sdslices.Ensure(reifyMenuItems(ctx, ec, token, tags, menu.Items)),
	}
}

func (item *MenuItem) Reify(ctx context.Context, ec echo.Context, token Token, tags []string) *MenuItem {
	// children
	reifiedChildren := reifyMenuItems(ctx, ec, token, tags, item.Children)
	if len(reifiedChildren) <= 0 {
		err := AccessControlCheck(ctx, ec, token, item.Object, ActionShow)
		if err != nil {
			return nil
		}
	}
	if len(item.Tags) > 0 {
		// 如果菜单项本身有item.Tags，说明它可以通过tags来过滤
		if len(lo.Intersect(tags, item.Tags)) <= 0 {
			return nil
		}
	}
	mapper := contextExpandMapper(ec)
	return &MenuItem{ // clone
		Name:     item.Name,
		Path:     sdstrings.ExpandShellLike(item.Path, mapper),
		Href:     sdstrings.ExpandShellLike(item.Href, mapper),
		Icon:     item.Icon,
		Object:   item.Object,
		Children: sdslices.Ensure(reifiedChildren),
		Tags:     slices.Clone(item.Tags),
	}
}

func reifyMenuItems(ctx context.Context, ec echo.Context, token Token, tags []string, items []*MenuItem) []*MenuItem {
	var filteredItems []*MenuItem
	for _, item := range items {
		if item != nil {
			if filteredItem := item.Reify(ctx, ec, token, tags); filteredItem != nil {
				filteredItems = append(filteredItems, filteredItem)
			}
		}
	}
	return filteredItems
}
