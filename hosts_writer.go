package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

func (app *App) startHostsWriter() {
	fmt.Println("Starting hosts file writer")

	app.emitter.On("domains-updated", app.updateHostsFile)
}

func (app *App) updateHostsFile() {
	cleanOldRecords(app)
	writeNewRecords(app)
}

func cleanOldRecords(app *App) {
	hosts, err := ioutil.ReadFile(*app.hostsFile)

	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile("(?m)^.*" + "docker-resolver" + ".*$[\n]+")
	text := re.ReplaceAllString(string(hosts), "")

	f, err := os.OpenFile(*app.hostsFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func writeNewRecords(app *App) {
	f, err := os.OpenFile(*app.hostsFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	text := createHostsString(app)
	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func createHostsString(app *App) string {
	text := ""

	for domain, ip := range app.records {
		text = text + fmt.Sprintf("%s  %s   # added by docker-resolver\n", ip, domain)
	}

	return text
}
