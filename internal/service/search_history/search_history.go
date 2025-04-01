package searchhistory

var (
	Instance *searchHistory
)

const (
	DEFAULT_CAPACITY = 10
)

type searchHistory struct {
	capacity int
	items    []string
}

func New() *searchHistory {
	if Instance != nil {
		return Instance
	}

	Instance = &searchHistory{
		capacity: DEFAULT_CAPACITY,
		items:    make([]string, 0),
	}

	return Instance
}

func (h *searchHistory) GetAll() []string {
	return h.items
}

func (h *searchHistory) GetCount(count int) []string {
	if count > len(h.items) {
		count = len(h.items)
	}

	return h.items[:count]
}

func (h *searchHistory) GetSize() int {
	return len(h.items)
}

func (h *searchHistory) Insert(item string) {
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
