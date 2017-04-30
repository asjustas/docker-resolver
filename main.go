package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/chuckpreslar/emission"
)

type BindFlags []string

type App struct {
	emitter   *emission.Emitter
	records   map[string]string
	hostsFile *string
	dnsBinds  BindFlags
}

func (i *BindFlags) String() string {
	return "my string representation"
}

func (i *BindFlags) Set(value string) error {
	*i = append(*i, value)

	return nil
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
	app.hostsFile = flag.String("hosts-file", "/etc/hosts", "Host file location")
	flag.Var(&app.dnsBinds, "dns-bind", "Dns server bind address")
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
