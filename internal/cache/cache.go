package cache

import (
	"sync"

	"github.com/d-jackalope/L0/pkg/models"
)

type Cache interface {
	Get(uid string) (*models.Order, bool)
	Set(uid string, order models.Order)
	Delete(uid string)
	Update(cache map[string]models.Order)

	RandomOrders() []models.Order
	CacheSize() int
}

type hashCache struct {
	mu    sync.RWMutex
	items map[string]models.Order
}

func New() Cache {
	cache := &hashCache{
		mu:    sync.RWMutex{},
		items: make(map[string]models.Order),
	}
	return cache
}

func (h *hashCache) Get(uid string) (*models.Order, bool) {
	h.mu.RLock()
	item, found := h.items[uid]
	h.mu.RUnlock()
	if !found {
		return nil, false
	}

	return &item, true
}

func (h *hashCache) Set(uid string, order models.Order) {
	h.mu.Lock()
	h.items[uid] = order
	h.mu.Unlock()
}

func (h *hashCache) Delete(uid string) {
	h.mu.Lock()
	_, found := h.items[uid]
	if found {
		delete(h.items, uid)
	}
	h.mu.Unlock()
}

func (h *hashCache) Update(cache map[string]models.Order) {
	h.items = cache
}

func (h *hashCache) RandomOrders() []models.Order {
	h.mu.Lock()
	lenght := len(h.items)
	if lenght == 0 {
		h.mu.Unlock()
		return []models.Order{}
	}

	var orders []models.Order
	count := 0
	for _, order := range h.items {
		if count > 20 {
			break
		}
		orders = append(orders, order)
		count++
	}
	h.mu.Unlock()
	return orders
}

func (h *hashCache) CacheSize() int {
	h.mu.Lock()
	lenght := len(h.items)
	h.mu.Unlock()
	return lenght
}
