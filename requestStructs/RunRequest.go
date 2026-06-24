package requeststructs

type Run struct {
	DataInput  string `json:"data_input" omitempty`
	Lang       string `json:"lang"`
	QuestionId string `json:"question_id"`
	TypedCode  string `json:"typed_code"`
}
