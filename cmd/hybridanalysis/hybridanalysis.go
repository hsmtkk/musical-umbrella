package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

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
	Use:  "get-report id | get-report sha256 environment-id",
	Run:  getReport,
	Args: cobra.MinimumNArgs(1),
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
	fileName := filepath.Base(filePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := hybridanalysis.New(apiKey).SubmitFile(fileName, content)
	if err != nil {
		log.Fatal(err)
	}
	parsed, err := hybridanalysis.NewParser().SubmitFile(resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(parsed)
}

func getReport(cmd *cobra.Command, args []string) {
	apiKey := env.RequiredEnv("API_KEY")
	var resp []byte
	var err error
	if len(args) == 1 {
		jobID := args[0]
		resp, err = hybridanalysis.New(apiKey).GetReportByJobID(jobID)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		sha256 := args[0]
		environmentID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal(err)
		}
		resp, err = hybridanalysis.New(apiKey).GetReportBySHA256(sha256, environmentID)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Print(string(resp))
}
