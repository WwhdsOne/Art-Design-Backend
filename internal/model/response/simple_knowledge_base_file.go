package response

type SimpleKnowledgeBaseFile struct {
	ID       int64  `json:"id,string"`
	Filename string `json:"filename"`
}
