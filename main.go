package main

import (
	"bytes"
	"encoding/json"
	"erp-job/cmd"
	"erp-job/models/fararavand"
	"fmt"
	"io"
	"net/http"

	"github.com/robfig/cron"
)

func main() {
	cmd.Execute()

	c := cron.New()
	// Schedule cron job to run every hour
	c.AddFunc("0 0 * * * *", func() {
		// Receive data from first API
		res, err := http.Get("http://firstapi.com/data")
		if err != nil {
			fmt.Println("Error receiving data from API:", err)
			return
		}
		defer res.Body.Close()

		// Read response body and unmarshal into struct
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var data []fararavand.BaseData
		json.Unmarshal(body, &data)

		// Send data to second API
		url := "http://secondapi.com/data"
		jsonData, _ := json.Marshal(data)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending data to second API:", err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			return
		}
		defer resp.Body.Close()

		// Print response from second API
		fmt.Println("Response from second API:", resp.Status)
	})
	c.Start()

	//Wait indefinitely so the program doesn't exit
}

func NewFararavand() {
	panic("unimplemented")
}
