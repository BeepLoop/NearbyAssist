package searchhistory

import "slices"

var (
	instance *searchHistory
)

const (
	DEFAULT_CAPACITY = 10
)

type searchHistory struct {
	capacity int
	items    []string
}

/*
Returns an instance of this. If instance is nil, creates an instance
then return it.
*/
func New() *searchHistory {
	if instance != nil {
		return instance
	}

	instance = &searchHistory{
		capacity: DEFAULT_CAPACITY,
		items:    make([]string, 0),
	}

	return instance
}

func Destroy() {
	if instance != nil {
		instance = nil
	}
}

func (h *searchHistory) GetAll() []string {
	return h.items
}

func (h *searchHistory) GetTopKElements(k int) []string {
	if k > len(h.items) {
		k = len(h.items)
	}

	return h.items[:k]
}

func (h *searchHistory) GetSize() int {
	return len(h.items)
}

func (h *searchHistory) Insert(item string) {
	if slices.Contains(h.items, item) {
		return
	}

	size := len(h.items)

	if size >= h.capacity {
		newItems := make([]string, 0)
		newItems = append(newItems, item)

		for _, oldItem := range h.items[:size-1] {
			newItems = append(newItems, oldItem)
		}

		h.items = newItems
		return
	}

	newItems := make([]string, 0)
	newItems = append(newItems, item)
	newItems = append(newItems, h.items...)
	h.items = newItems
}
