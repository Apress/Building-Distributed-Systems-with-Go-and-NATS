package domain

import "errors"

// MenuItem is a Value Object because it's identified by its attributes (Name, Price)
// and has no independent identity. It should be immutable.
type MenuItem struct {
	Name  string
	Price float64
}

func NewMenuItem(name string, price float64) (MenuItem, error) {
	if name == "" || price <= 0 {
		return MenuItem{}, errors.New("menu item name and price must be valid")
	}
	return MenuItem{Name: name, Price: price}, nil
}