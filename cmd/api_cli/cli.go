package main

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	//	"github.com/SecurityDo/cloudv2/farm"

	api "github.com/SecurityDo/ingext_api/client"

	//lv3API "github.com/SecurityDo/cloudv2/lv3/api"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
	//"github.com/aws/aws-sdk-go-v2/service/s3"
)

var client *api.IngextClient

func PrettyPrintJSON(x interface{}) {
	pretty, _ := json.MarshalIndent(x, "", "   ")
	fmt.Printf("%s\n", pretty)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			println("Sending SIGINT to everyone. (just in case?)")
			syscall.Kill(syscall.Getpid()*(-1), syscall.SIGINT)
			println("Done")
		}
	}()

	VERSION := "0.1.0"
	PROMPT := "api_cli v" + VERSION + " $ "

	shell := ishell.NewWithConfig(&readline.Config{Prompt: PROMPT})
	shell.SetHomeHistoryPath(".apicli_shell_history")

	shell.AddCmd(&ishell.Cmd{
		Name: "set",
		Help: "set: <url> <token>",
		Func: func(c *ishell.Context) {
			//node.Commit()

			if len(c.Args) != 2 {
				fmt.Println("Usage: set: <url> <token>")
				return
			}
			url := c.Args[0]
			token := c.Args[1]

			client = api.NewIngextClient(url, token, false, nil)
			shell.Println()
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "call",
		Help: "call <function> <payload>",
		Func: func(c *ishell.Context) {
			//node.Commit()
			if len(c.Args) != 2 {
				fmt.Printf("Usage: call <function> <payload>")
				shell.Println()
				return
			}
			functionName := c.Args[0]
			payload := c.Args[1]
			b, _ := os.ReadFile(payload)

			var m map[string]interface{}
			json.Unmarshal(b, &m)

			result, err := client.GenericCall("api/ds", functionName, m)
			if err != nil {
				fmt.Printf("Error calling function %s: %s\n", functionName, err.Error())
				shell.Println()
				return
			}
			PrettyPrintJSON(result)

			shell.Println()
		},
	})

	if len(os.Args) > 1 {
		err := shell.Process(os.Args[1:]...)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		shell.Println("Interactive Shell")
		// start shell
		shell.Run()
		// teardown
		shell.Close()
	}

}
