package main

import (
	"fmt"
	"log"

	"github.com/hsmtkk/musical-umbrella/cmd/env"
	"github.com/hsmtkk/musical-umbrella/virustotal"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{}

var uploadFileCommand = &cobra.Command{
	Use: "upload-file",
	Run: uploadFile,
}

var getReportCommand = &cobra.Command{
	Use: "get-report",
	Run: getReport,
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
	apiKey := env.RequiredEnv("API_KEY")
	v := virustotal.New(apiKey)
	resp, err := v.UploadFile("eicar", []byte(`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resp))
}

func getReport(cmd *cobra.Command, args []string) {
}
