package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

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

func getCodeContent() (string, error) {
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

func createRunObject() requeststructs.Run {

	cConfig, err := getCodeConfig()
	if err != nil {
		log.Fatal("Cant read code config", err)
	}

	// read the code from the code file
	codeContent, err := getCodeContent()

	if err != nil {
		log.Fatal("Error reading code file", err)
	}

	runObject := requeststructs.Run{
		DataInput:  cConfig.DataInput,
		Lang:       cConfig.Lang,
		QuestionId: cConfig.QuestionId,
		TypedCode:  codeContent,
	}

	return runObject
}

func (a *API) Run() {
	// question session

	cConfig, err := getCodeConfig()
	runObject := createRunObject()

	fmt.Println("CODE: \n", runObject.TypedCode)

	runObjBytes, err := json.Marshal(runObject)

	if err != nil {
		log.Fatal("Unable to marshal runObject for making run request", err)
	}

	// send it and receive the response object

	submitLink := fmt.Sprintf("%sproblems/%s/interpret_solution/", constants.LEETCODE_BASE, cConfig.Slug)
	runRequest, err := request.NewRequest("POST", submitLink, bytes.NewBuffer(runObjBytes))

	response, err := request.MakeRequest(runRequest, a.session)

	if err != nil {
		log.Fatal("Error making run request", err)
	}

	defer response.Body.Close()

	var resStruct map[string]any

	err = json.NewDecoder(response.Body).Decode(&resStruct)

	if err != nil {
		log.Fatal("Error decoding response json", err)
	}

	interpretId := resStruct["interpret_id"].(string)

	fmt.Println("Interpret_ID:", interpretId)
	a.getRunState(interpretId)

}

func (a *API) getRunState(interpret_id string) {

	if interpret_id == "" {
		fmt.Println("no interpret_id passed, passed value '", interpret_id, "'")
	}

	state := ""
	runCheckUrl := strings.Replace(constants.RUN_CHECK_URL, constants.SUBMISSION_ID_PLACEHOLDER, interpret_id, 1)

	var resObj map[string]any
	for {
		checkRequest, _ := request.NewRequest("GET", runCheckUrl, nil)

		response, _ := request.MakeRequest(checkRequest, a.session)

		defer response.Body.Close()

		err := json.NewDecoder(response.Body).Decode(&resObj)

		if err != nil {
			log.Fatal("Unable to decode response json in getting run state", err)
		}

		newState := resObj["state"].(string)

		if newState != state {
			fmt.Println(resObj["state"])
		}
		state = newState

		if state != "PENDING" && state != "STARTED" {
			break
		}

	}

	fmt.Println("====")
	fmt.Printf("%v\n", resObj)

}

func (a *API) Submit() {
	runObject := createRunObject()
	codeConfig, _ := getCodeConfig()
	submitLink := strings.Replace(constants.SUBMIT_URL, constants.QUESTION_SLUG_PLACEHOLDER, codeConfig.Slug, 1)

	codeBytes, _ := json.Marshal(runObject)

	submitRequest, _ := request.NewRequest("POST", submitLink, bytes.NewReader(codeBytes))

	resp, _ := request.MakeRequest(submitRequest, a.session)

	defer resp.Body.Close()

	var submitRes map[string]any

	fmt.Println(resp)

	_ = json.NewDecoder(resp.Body).Decode(&submitRes)

	fmt.Println(submitRes)

	submissionId := strconv.Itoa(int(submitRes["submission_id"].(float64)))

	if submissionId == "" {
		log.Fatal("submissionId not found", submitRes)
	}

	a.checkSubmission(submissionId)

}

func (a *API) checkSubmission(submissionId string) {
	fmt.Println("Got submissionId:", submissionId)
	submitCheckUrl := strings.ReplaceAll(constants.SUBMISSION_CHECK_URL, constants.SUBMISSION_ID_PLACEHOLDER, submissionId)
	state := ""
	var resObj map[string]any

	for {
		subCheckReq, _ := request.NewRequest("GET", submitCheckUrl, nil)

		resp, _ := request.MakeRequest(subCheckReq, a.session)

		defer resp.Body.Close()

		json.NewDecoder(resp.Body).Decode(&resObj)
		if stateAny, ok := resObj["state"]; ok {
			newState := stateAny.(string)

			if newState != state {
				fmt.Println("state:", state)
			}
			state = newState

		} else {
			break
		}

		if state == "SUCCESS" {
			break
		}

	}

	resBytes, _ := json.MarshalIndent(resObj, "", " ")

	fmt.Println(string(resBytes))

}

func (a *API) ClearWorkspace(){

	err := os.RemoveAll(constants.CODE_DIR)

	if err != nil {
			log.Fatal("Cannot clear workspace", err)
	}else {
		fmt.Println("Workspace cleared")
	}

}
