package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hsmtkk/musical-umbrella/cmd/env"
	"github.com/hsmtkk/musical-umbrella/virustotal"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{}

var uploadFileCommand = &cobra.Command{
	Use:  "upload-file file",
	Run:  uploadFile,
	Args: cobra.ExactArgs(1),
}

var getReportCommand = &cobra.Command{
	Use:  "get-report id",
	Run:  getReport,
	Args: cobra.ExactArgs(1),
}

func init() {
	command.AddCommand(uploadFileCommand)
	command.AddCommand(getReportCommand)
}

func main() {
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}

func uploadFile(cmd *cobra.Command, args []string) {
	filePath := args[0]
	apiKey := env.RequiredEnv("API_KEY")
	v := virustotal.New(apiKey)

	fileName := filepath.Base(filePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := v.UploadFile(fileName, content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resp))
}

func getReport(cmd *cobra.Command, args []string) {
	id := args[0]
	apiKey := env.RequiredEnv("API_KEY")
	v := virustotal.New(apiKey)

	resp, err := v.GetReport(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resp))
}
