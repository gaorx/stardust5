package sdecho

import (
	"context"
	"github.com/gaorx/stardust5/sdslices"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/labstack/echo/v4"
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
}

func MenuReify(ctx context.Context, ec echo.Context, menuId string) *Menu {
	menus := MustGet[Menus](ec, keyMenus)
	for _, menu := range menus {
		if menu != nil && menu.Id == menuId {
			token, err := TokenDecode(ctx, ec)
			if err != nil {
				return &Menu{Id: menuId, Items: []*MenuItem{}}
			}
			return menu.Reify(ctx, ec, token)
		}
	}
	return nil
}

func (menu *Menu) Reify(ctx context.Context, ec echo.Context, token Token) *Menu {
	return &Menu{
		Id:    menu.Id,
		Items: sdslices.Ensure(reifyMenuItems(ctx, ec, token, menu.Items)),
	}
}

func (item *MenuItem) Reify(ctx context.Context, ec echo.Context, token Token) *MenuItem {
	// children
	filteredChildren := reifyMenuItems(ctx, ec, token, item.Children)
	if len(filteredChildren) <= 0 {
		err := AccessControlCheck(ctx, ec, token, item.Object, ActionShow)
		if err != nil {
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
		Children: sdslices.Ensure(filteredChildren),
	}
}

func reifyMenuItems(ctx context.Context, ec echo.Context, token Token, items []*MenuItem) []*MenuItem {
	var filteredItems []*MenuItem
	for _, item := range items {
		if item != nil {
			if filteredItem := item.Reify(ctx, ec, token); filteredItem != nil {
				filteredItems = append(filteredItems, filteredItem)
			}
		}
	}
	return filteredItems
}
