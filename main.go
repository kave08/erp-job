package main

import (
	"erp-job/cmd"

	"github.com/robfig/cron"
)

func main() {
	cmd.Execute()

	c := cron.New()
	// Schedule cron job to run every hour
	c.AddFunc("0 0 * * * *", func() {
		
	})
	c.Start()

	//Wait indefinitely so the program doesn't exit
}
