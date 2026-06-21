package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	request "github.com/Shivam-Verma9999/go-leetcode/Request"
	"github.com/Shivam-Verma9999/go-leetcode/constants"
	"github.com/Shivam-Verma9999/go-leetcode/response"
	"github.com/Shivam-Verma9999/go-leetcode/session"
)

type API struct {
	session *session.Session
}

func New(session *session.Session) *API {
	return &API{
		session: session,
	}
}

func (a *API) GetQuestion(questionName string) (*response.QuestionResponse, error) {

	if questionName == "" {
		fmt.Println("Question name required")
		return nil, fmt.Errorf("Expected QuestionName")
	}

	payload, _ := os.ReadFile("./queryPayloads/questionDetail.json")

	payload = bytes.Replace(payload, []byte("<<questionName>>"), []byte(questionName), 1)

	questtionReq, _ := request.NewRequest("POST", constants.LEETCODE_GQL_URL, bytes.NewBuffer(payload))

	//fmt.Println(questtionReq.Body)

	qRes, _ := request.MakeRequest(questtionReq, a.session.Client)
	defer qRes.Body.Close()

	//b, _ := io.ReadAll(qRes.Body)
	//fmt.Println(string(b))

	var dataBody response.QuestionResponse

	if err := json.NewDecoder(qRes.Body).Decode(&dataBody); err != nil {
		fmt.Println("Error parsing json body", err)
		return nil, err
	}

	return &dataBody, nil
}
