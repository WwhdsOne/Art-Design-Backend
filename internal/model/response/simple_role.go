package response

type SimpleRole struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}
