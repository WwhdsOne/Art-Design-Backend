package request

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" binding:"required,max=10"`
	Code        string `json:"code" binding:"required,min=5,max=10"`
	Description string `json:"description"`
	Status      int8   `json:"status"`
}
