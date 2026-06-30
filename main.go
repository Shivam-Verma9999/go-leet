package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	flagparser "github.com/Shivam-Verma9999/go-leetcode/FlagParser"
	"github.com/Shivam-Verma9999/go-leetcode/api"
	"github.com/Shivam-Verma9999/go-leetcode/config"
	"github.com/Shivam-Verma9999/go-leetcode/constants"
	response "github.com/Shivam-Verma9999/go-leetcode/responseStructs"
	"github.com/Shivam-Verma9999/go-leetcode/session"
	sharedstructs "github.com/Shivam-Verma9999/go-leetcode/sharedStructs"
)

func createWorkspace(questionResponse *response.QuestionResponse) {
	err := os.MkdirAll(constants.CODE_DIR, 0644)

	if err != nil {
		log.Fatalf("Error creating ./code directory %v", err)
	}

	if questionResponse == nil {
		log.Fatalln("Nil DataContainer for saving example test cases")
	}

	inputFile, err := os.OpenFile(path.Join(constants.CODE_DIR, "input.txt"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("Error opening input.txt for writing test cases")
	}

	defer inputFile.Close()

	outputFile, err := os.OpenFile(path.Join(constants.CODE_DIR, "output.txt"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("Error opening output.txt for writing test cases")
	}
	defer outputFile.Close()

	for _, tCase := range questionResponse.Data.Question.ExampleTestCases {
		inp, out, _ := strings.Cut(tCase, "\n")

		fmt.Println("tcases ")
		fmt.Println(inp)
		fmt.Println(out)

		replacer := strings.NewReplacer(
			",", " ",
			"[", "",
			"]", "")

		inputFile.WriteString(replacer.Replace(inp))
		inputFile.WriteString("\n")
		outputFile.WriteString(out)
		outputFile.WriteString("\n")

	}

	mainCodeFile, err := os.OpenFile(path.Join(constants.CODE_DIR, "main.cpp"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error opening code file, %v", err)
	}

	defer mainCodeFile.Close()

	for _, snippet := range questionResponse.Data.Question.CodeSnippets {
		if snippet.LangSlug == "cpp" {
			mainCodeFile.WriteString(snippet.Code)
		}
	}

	codeConfigFile, err := os.OpenFile(path.Join(constants.CODE_DIR, "codeConfig.json"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalln("Error creating codeConfig", err)
	}
	defer codeConfigFile.Close()

	codeConfigObject := &sharedstructs.CodeConfig{
		DataInput:  "",
		Lang:       "cpp",
		QuestionId: questionResponse.Data.Question.QuestionId,
		Slug:       questionResponse.Data.Question.Slug,
	}
	codeConfigEncoder := json.NewEncoder(codeConfigFile)
	codeConfigEncoder.SetIndent("", "  ")
	codeConfigEncoder.Encode(codeConfigObject)

}

var questionSlug string

func main() {

	// parse params

	fmt.Println("length", len(os.Args))
	fmt.Println(os.Args)

	flags := flagparser.Parse(os.Args[1:])
	fmt.Println("Parsed flags", flags)

	os.Exit(0)
	
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Fallback to sddilencing if file cannot be created
		log.SetOutput(io.Discard)
		return
	} else {
		log.SetOutput(logFile)
	}

	//os.Args[0] == program name

	cfg, err := config.Load()

	if err != nil {
		log.Fatalln("Error loading config", err)
	}

	session, err := session.New(cfg)
	if err != nil {
		log.Fatalln("Error creating session", err)

	}

	api := api.New(session)

	if len(os.Args) > 1 {
		questionSlug = os.Args[1]
		fmt.Println("Args", os.Args[1])
	}
	fmt.Println(os.Args)

	dataBody, err := api.GetQuestion("two-sum")

	if err != nil {
		log.Fatal("Error in GetQuestion API", err)
	}

	fmt.Printf("\n question body\n %v\n", dataBody.Data.Question.Content)
	createWorkspace(dataBody)

	api.Run()
	//	api.Submit()

}
