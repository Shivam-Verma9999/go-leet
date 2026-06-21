package response

type CodeSnippet struct {
	Code string `json:"code"`
	Lang string `json:"lang"`
	LangSlug string `json:"langSlug"`
}
 

type QuestionContainer  struct {
	Content string `json:"content"`
	ExampleTestCases []string `json:"exampleTestcaseList"`   
	CodeSnippets []CodeSnippet `json:"codeSnippets"`
}

type DataContainer struct {
	Question QuestionContainer  `json:"question"`
}
type QuestionResponse struct {
	Data DataContainer `json:"data"`
}
