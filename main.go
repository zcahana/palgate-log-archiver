package main

import (
	"log"

	"github.com/zcahana/palgate-log-archiver/sink/googlesheets"
	"github.com/zcahana/palgate-sdk"
)

func main() {
	log.Printf("Loading PalGate configuration")
	config, err := palgate.InitConfig()
	if err != nil {
		log.Fatalf("Error parsing configuration: %v", err)
	}

	err = config.Validate()
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	client := palgate.NewClient(config)

	log.Printf("Acuiring log records from PalGate server")
	logResp, err := client.Log()
	if err != nil {
		log.Fatalf("Error executing palgate command: %v", err)
	}

	if logResp.Status != palgate.ResponseStatusSuccess {
		log.Fatalf("Error executing palgate command: status=%s, error=%s, message=%s",
			logResp.Status, logResp.Error, logResp.Message)
	}

	log.Printf("Storing log records in Google Sheets")
	sink, err := googlesheets.NewSink()
	if err != nil {
		log.Fatalf("Error initializing Google Sheets sink: %v", err)
	}

	rowsWritten, err := sink.Receive(logResp.Records)
	if err != nil {
		log.Fatalf("Error processing palgate log records: %v", err)
	}

	log.Printf("Archived %d new log records", rowsWritten)
}
