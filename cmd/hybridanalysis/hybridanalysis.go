package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hsmtkk/musical-umbrella/cmd/env"
	"github.com/hsmtkk/musical-umbrella/hybridanalysis"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{}

var submitFileCommand = &cobra.Command{
	Use:  "submit-file file",
	Run:  submitFile,
	Args: cobra.ExactArgs(1),
}

var getReportCommand = &cobra.Command{
	Use:  "get-report id",
	Run:  getReport,
	Args: cobra.ExactArgs(1),
}

func init() {
	command.AddCommand(submitFileCommand)
	command.AddCommand(getReportCommand)
}

func main() {
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}

func submitFile(cmd *cobra.Command, args []string) {
	filePath := args[0]
	apiKey := env.RequiredEnv("API_KEY")
	h := hybridanalysis.New(apiKey)

	fileName := filepath.Base(filePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := h.SubmitFile(fileName, content)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(resp)
}

func getReport(cmd *cobra.Command, args []string) {
}
