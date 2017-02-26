package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/chuckpreslar/emission"
)

type App struct {
	emitter   *emission.Emitter
	records   map[string]string
	hostsFile *string
	dnsBind   *string
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)

	app := new(App)
	app.emitter = emission.NewEmitter()
	app.records = make(map[string]string)

	parseFlags(app)

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

func parseFlags(app *App) {
	app.dnsBind = flag.String("dns-bind", ":53", "Dns server bind address")
	app.hostsFile = flag.String("hosts-file", "/etc/hosts", "Host file location")
	help := flag.Bool("help", false, "Show usage")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(1)
	}
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
