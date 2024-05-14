package utils

import "nearbyassist/internal/response"

// Sort the services based on Suggestability in descending order
func BubbleSort(elements []*response.SearchResult) []*response.SearchResult {
	for i := 0; i < len(elements); i++ {
		for j := 0; j < len(elements)-1; j++ {
			if elements[j].Suggestability < elements[j+1].Suggestability {
				elements[j], elements[j+1] = elements[j+1], elements[j]
			}
		}
	}

	return elements
}
