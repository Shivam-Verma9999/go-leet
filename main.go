package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shivam-Verma9999/go-leetcode/api"
	"github.com/Shivam-Verma9999/go-leetcode/config"
	"github.com/Shivam-Verma9999/go-leetcode/response"
	"github.com/Shivam-Verma9999/go-leetcode/session"
)

func saveCodeTemplateAndTestCases(questionResponse *response.QuestionResponse) {
	codeDir := "./code/"
	err := os.MkdirAll(codeDir, 0644)

	if err != nil {
		log.Fatalf("Error creating ./code directory %v", err)
	}

	if questionResponse == nil {
		log.Fatalln("Nil DataContainer for saving example test cases")
	}

	inputFile, err := os.OpenFile(codeDir+"input.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("Error opening input.txt for writing test cases")
	}

	defer inputFile.Close()

	outputFile, err := os.OpenFile(codeDir+"output.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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

	mainCodeFile, err := os.OpenFile(codeDir+"main.cpp", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error opening code file, %v", err)
	}

	defer mainCodeFile.Close()

	for _, snippet := range questionResponse.Data.Question.CodeSnippets {
		if snippet.LangSlug == "cpp" {
			mainCodeFile.WriteString(snippet.Code)
		}
	}

}

var questionSlug string

func main() {

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

	dataBody, err := api.GetQuestion()

	if err != nil {
		log.Fatal("Error in GetQuestion API", err)
	}

	fmt.Printf("\n question body\n %v", dataBody.Data.Question.Content)

	saveCodeTemplateAndTestCases(dataBody)
}
