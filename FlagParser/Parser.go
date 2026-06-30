package flagparser

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Flags struct {
	Language        string
	QuestionSlug    string
	QuestionId      string
	commandRegistry CommandRegistry
	flagRegistry    FlagRegistry
}

type CommandRegistry map[string]any

// flag with wether a value is expected or not
type FlagRegistry map[string]bool

var supportedCommandsList []string = []string{
	"init",
	"clear",
}

func (c *CommandRegistry) IsSupported(cmd string) bool {
	_, exists := (*c)[cmd]
	return exists
}

func (f *FlagRegistry) IsSupported(flag string) bool {
	_, exists := (*f)[flag]
	return exists
}

func (f *FlagRegistry) ValueRequied(flag string) bool {
	req, exists := (*f)[flag]

	return exists && req

}

var supportedFlag map[string]bool = map[string]bool{
	"--lang": true,
	"--test": false,
}

func NewCommandRegistry(commandList []string) CommandRegistry {
	reg := make(CommandRegistry)
	for _, cmd := range commandList {
		reg[cmd] = struct{}{}
	}
	return reg
}

func NewFlagRegistry(flagList map[string]bool) FlagRegistry {
	return flagList
}

func isCommand(cmd string) bool {
	res := !strings.HasPrefix(cmd, "--")
	fmt.Println("checking isCommand for", cmd, "result:", res)
	return res
}

func isFlag(flag string) bool {
	res := strings.HasPrefix(flag, "--")

	fmt.Println("Checking isFlag for", flag, "resut:", res)
	return res 
}

func (f *Flags) HandleCommand(cmd string) {
	if !f.commandRegistry.IsSupported(cmd) {
		fmt.Println("Command not supported", cmd)
		os.Exit(1)
	}
}

func (f *Flags) HandleFlag(flag string, flagValue string) {

	if !f.flagRegistry.IsSupported(flag) {
		fmt.Println("Unsupported flag ", flag)
		os.Exit(1)
	}

	if (flagValue == "" || strings.HasPrefix(flagValue, "--")) && f.flagRegistry.ValueRequied(flag) {
		fmt.Printf("Flag:%v needs value\n", flag)
		log.Fatalf("for flag %v  next value was `%v`", flag, flagValue)
	}

	switch flag {
		case "--lang": 
			f.Language = flagValue
	}

}

func Parse(args []string) *Flags {
	argsLength := len(args)
	f := &Flags{}
	f.commandRegistry = NewCommandRegistry(supportedCommandsList)
	f.flagRegistry = NewFlagRegistry(supportedFlag)

	fmt.Println(f.flagRegistry)
	for i := 0; i < argsLength; i++ {
		_arg := strings.ToLower(args[i])

		if isCommand(_arg) {
			f.HandleCommand(_arg)
		} else if isFlag(_arg) {
			flagValue := ""
			if f.flagRegistry.ValueRequied(_arg) && i+1 < argsLength {
				i++
				flagValue = args[i]
			}
			f.HandleFlag(_arg, flagValue)
		} else {
			fmt.Println("Unsupported flag/ command")
			log.Fatalln("Unsupported flag/ command")
		}
	}

	return f

}
