package persistence

import (
	"chapter_8/chapter-8/domain" // Depends on domain layer
)

// InMemoryMenuItemRepository implements an interface (implicitly or explicitly defined in application layer)
// for fetching menu items.
type InMemoryMenuItemRepository struct {
	items []domain.MenuItem
}

func NewInMemoryMenuItemRepository(initialItems []domain.MenuItem) *InMemoryMenuItemRepository {
	return &InMemoryMenuItemRepository{
		items: initialItems,
	}
}

func (r *InMemoryMenuItemRepository) GetAll() ([]domain.MenuItem, error) {
	// Return a copy to ensure immutability from external changes, protecting the internal state.
	copiedItems := make([]domain.MenuItem, len(r.items))
	copy(copiedItems, r.items)
	return copiedItems, nil
}

func (r *InMemoryMenuItemRepository) FindByName(name string) (domain.MenuItem, bool) {
	for _, item := range r.items {
		if item.Name == name {
			return item, true
		}
	}
	return domain.MenuItem{}, false
}