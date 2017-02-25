package main

import (
	"fmt"
	"runtime"

	"github.com/chuckpreslar/emission"
)

const (
	HOSTS_FILE = "/tmp/hosts"
)

type App struct {
	emitter *emission.Emitter
	records map[string]string
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)
	fmt.Println("Starting app")

	app := new(App)
	app.emitter = emission.NewEmitter()
	app.records = make(map[string]string)

	containerStart := func(domains []string, ip string) {
		fmt.Printf("ContainerStart %s\n%s\n\n", domains, ip)
	}

	containerStop := func(domains []string) {
		fmt.Printf("ContainerStop %s\n\n", domains)
	}

	app.emitter.On("container-start", containerStart)
	app.emitter.On("container-stop", containerStop)

	go app.startDockerListener()
	app.startHostsWriter()
	app.startDNSServer()
}

func (app *App) registerDomains(domains []string, ip string) {
	for _, domain := range domains {
		app.records[domain] = ip
	}
}

func (app *App) removeDomains(domains []string) {
	for _, domain := range domains {
		delete(app.records, domain)
	}
}
