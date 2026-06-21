package requeststructs

type Run struct {
	DataInput  string `json:"data_input"`
	Lang       string `json:"lang"`
	QuestionId string `json:"question_id"`
	TypeCode   string `json:"typed_code"`
}
