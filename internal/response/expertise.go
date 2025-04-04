package response

type Expertise struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Tags  []Tag  `json:"tags"`
}
