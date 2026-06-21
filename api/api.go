package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	request "github.com/Shivam-Verma9999/go-leetcode/Request"
	"github.com/Shivam-Verma9999/go-leetcode/constants"
	requeststructs "github.com/Shivam-Verma9999/go-leetcode/requestStructs"
	responsestructs "github.com/Shivam-Verma9999/go-leetcode/responseStructs"
	"github.com/Shivam-Verma9999/go-leetcode/session"
	sharedstructs "github.com/Shivam-Verma9999/go-leetcode/sharedStructs"
)

type API struct {
	session *session.Session
}

func New(session *session.Session) *API {
	return &API{
		session: session,
	}
}

var langFileExtensionMap map[string]string = map[string]string{
	"cpp": "cpp",
	"go":  "go",
}

func (a *API) GetQuestion(questionName string) (*responsestructs.QuestionResponse, error) {

	if questionName == "" {
		fmt.Println("Question name required")
		return nil, fmt.Errorf("Expected QuestionName")
	}

	payload, _ := os.ReadFile("./queryPayloads/questionDetail.json")

	payload = bytes.Replace(payload, []byte("<<questionName>>"), []byte(questionName), 1)

	questtionReq, _ := request.NewRequest("POST", constants.LEETCODE_GQL_URL, bytes.NewBuffer(payload))

	//fmt.Println(questtionReq.Body)

	qRes, _ := request.MakeRequest(questtionReq, a.session)
	defer qRes.Body.Close()

	//b, _ := io.ReadAll( sfcqRes.Body)
	//fmt.Println(string(b))

	var dataBody responsestructs.QuestionResponse

	if err := json.NewDecoder(qRes.Body).Decode(&dataBody); err != nil {
		fmt.Println("Error parsing json body", err)
		return nil, err
	}

	return &dataBody, nil
}

func getCodeConfig() (*sharedstructs.CodeConfig, error) {

	codeConfigContent, err := os.ReadFile(path.Join(constants.CODE_DIR, "codeConfig.json"))

	if err != nil {
		fmt.Println("Error Reading CodeConfig", err)
		return nil, err
	}

	codeConfigDecoder := json.NewDecoder(bytes.NewReader(codeConfigContent))
	var codeConfig sharedstructs.CodeConfig
	if err := codeConfigDecoder.Decode(&codeConfig); err != nil {
		fmt.Println("Unable to parse codeConfigFile", err)
		return nil, err
	}

	return &codeConfig, nil
}

func getCode() (string, error) {
	cConfig, err := getCodeConfig()

	if err != nil {
		fmt.Println("Error getting code config while getting code", err)
		return "", err
	}
	ext, exists := langFileExtensionMap[cConfig.Lang]

	if !exists {
		fmt.Println("Unsupported code file")
		return "", fmt.Errorf("Unsupported code file")
	}

	codeContent, err := os.ReadFile(path.Join(constants.CODE_DIR, "main."+ext))
	return string(codeContent), nil

}

func (a *API) Run() {
	// question session
	cConfig, err := getCodeConfig()

	if err != nil {
		log.Fatal("Cant read code config", err)
	}

	// read the code from the code file
	codeContent, err := getCode()

	if err != nil {
		log.Fatal("Error reading code file", err)
	}
	// create Run Request

	runObject := requeststructs.Run{
		DataInput:  cConfig.DataInput,
		Lang:       cConfig.Lang,
		QuestionId: cConfig.QuestionId,
		TypeCode:   codeContent,
	}

	runObjBytes, err := json.Marshal(runObject)

	if err != nil {
		log.Fatal("Unable to marshal runObject for making run request", err)
	}

	// send it and receive the response object

	submitLink := fmt.Sprintf("%sproblems/%s/interpret_solution/", constants.LEETCODE_BASE, cConfig.Slug)
	fmt.Println("submit link", submitLink)

	runRequest, err := request.NewRequest("POST", submitLink, bytes.NewBuffer(runObjBytes))

	response, err := request.MakeRequest(runRequest, a.session)

	if err != nil {
		log.Fatal("Error making run request", err)
	}

	defer response.Body.Close()

	res, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response", err)
	}

	fmt.Println(string(res))

	// poll back the result

}
