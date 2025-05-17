package response

type SimpleMenu struct {
	ID       int64         `json:"id,string"`
	Title    string        `json:"title"`
	ParentID int64         `json:"parent_id"`
	Children []*SimpleMenu `json:"children"`
}
