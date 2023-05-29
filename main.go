package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	gpt3 "github.com/PullRequestInc/go-gpt3"

	//viper & cobra are 2 widely used tools when developing CLIs in GO
	//ex: Github CLI is an opensource tool for using Github from your computer's command line
	"github.com/spf13/cobra" //commander for Go CLI interactions
	"github.com/spf13/viper"
)

func GetResponse(client gpt3.Client, ctx context.Context, question string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			question,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	fmt.Printf("\n")
}

type NullWriter int // Null writer implements the io.Write interface but doesn't do anything

func (NullWriter) Write([]byte) (int, error) { return 0, nil } //Write implements the io.Write interface but is no-op (output)

func main() {
	log.SetOutput(new(NullWriter))
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apikey := viper.GetString("API_KEY")
	if apikey == "" {
		panic("Missing API Key")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apikey)

	//now we are going to write code for the CLI
	rootcmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with chatGpt in console",
		Run: func(cmd *cobra.Command, args []string) { //this just like pvsm in java
			scanner := bufio.NewScanner(os.Stdin) //NewScanner-> returns a new scanner to read from r ,(A scanner reads a string (or file), and converts the string into a stream of tokens for the a language.)
			quit := false

			for !quit {
				fmt.Print("say something ('quit' to end)")
				if !scanner.Scan() { //mtlb kuch likha hi nai ya kuch Scan ni hua toh break
					break
				}

				question := scanner.Text() // jo question likha hai wo
				switch question {
				case "quit":
					quit = true

				default:
					GetResponse(client, ctx, question)
				}
			}
		},
	}
	rootcmd.Execute()
}
