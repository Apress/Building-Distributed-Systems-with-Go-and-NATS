package domain

type MenuService struct {
	// This service holds the current menu, typically populated from the application layer.
	menu map[string]MenuItem
}

func NewMenuService(items []MenuItem) *MenuService {
	m := make(map[string]MenuItem)
	for _, item := range items {
		m[item.Name] = item
	}
	return &MenuService{menu: m}
}

func (s *MenuService) IsItemValid(itemName string) bool {
	_, exists := s.menu[itemName]
	return exists
}

func (s *MenuService) GetMenuItem(itemName string) (MenuItem, bool) {
	item, exists := s.menu[itemName]
	return item, exists
}