package main

import (
	"ccruncher"
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"
)

type App struct {
	GUID     string
	Requests []Request
}

type Request struct {
	RequestID  string
	LogEntries []ccruncher.LogEntry
}

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	log, err := os.Open(os.Args[1])

	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not open file ", os.Args[1])
		os.Exit(1)
	}

	ccLog, _ := ccruncher.ParseLog(log)

	outFile, err := os.Create("ccLog.yml")

	if err != nil {
		os.Exit(2)
	}

	defer outFile.Close()

	appGuids := ccLog.Apps()

	for _, guid := range appGuids {
		if guid != "unspecified" {
			app := &App{
				GUID: guid,
			}
			ids := ccLog.RequestsForApp(guid)

			for _, id := range ids {
				app.Requests = append(app.Requests, Request{
					RequestID:  id,
					LogEntries: ccLog.EntriesForRequest(id),
				})

				e, _ := candiedyaml.Marshal(app)
				outFile.Write(e)
				// entries := ccLog.EntriesForRequest(id)
				//
				// for _, entry := range entries {
				// 	e, _ := entry.Render()
				// 	outFile.Write(e)
				// }
			}
		}
	}
}
