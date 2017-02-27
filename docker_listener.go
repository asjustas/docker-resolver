package main

import (
	"fmt"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

func (app *App) startDockerListener() {
	fmt.Println("Starting docker events listener")

	client, err := docker.NewClient("unix:///var/run/docker.sock")

	if err != nil {
		panic(err)
	}

	registerRunningContainers(app, client)

	events := make(chan *docker.APIEvents)
	err = client.AddEventListener(events)

	if err != nil {
		panic(err)
	}

	for event := range events {
		if "start" == event.Action || "stop" == event.Action || "kill" == event.Action || "die" == event.Action {
			domains := getDomains(client, event.ID)

			if "start" == event.Action {
				ip := getContainerIp(client, event.ID)

				app.registerDomains(domains, ip)
				app.emitter.Emit("container-start", domains, ip)
			} else {
				app.removeDomains(domains)
				app.emitter.Emit("container-stop", domains)
			}

			app.emitter.Emit("domains-updated")
		}
	}
}

func registerRunningContainers(app *App, client *docker.Client) {
	fmt.Println("Registering running containers")

	containers, err := client.ListContainers(docker.ListContainersOptions{})

	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		domains := getDomains(client, container.ID)
		ip := getContainerIp(client, container.ID)

		app.registerDomains(domains, ip)
	}

	app.emitter.Emit("domains-updated")
}

func getDomains(client *docker.Client, ID string) []string {
	domains := []string{}
	container, _ := client.InspectContainer(ID)

	if "" != container.Config.Domainname {
		domains = append(domains, container.Config.Hostname+"."+container.Config.Domainname)
	}

	domains = append(domains, container.Name[1:]+".docker")
	envDomains := getDomainsFromEnv(container.Config.Env)

	for _, domain := range envDomains {
		domains = append(domains, domain)
	}

	return domains
}

func getContainerIp(client *docker.Client, ID string) string {
	container, _ := client.InspectContainer(ID)

	return container.NetworkSettings.IPAddress
}

func getDomainsFromEnv(args []string) []string {
	domains := []string{}

	for _, arg := range args {
		env := strings.Split(arg, "=")

		if "DOMAIN_NAME" == env[0] || "DNSDOCK_ALIAS" == env[0] {
			if strings.Contains(env[1], ",") {
				splited := strings.Split(env[1], ",")

				for _, domain := range splited {
					domains = append(domains, domain)
				}
			} else {
				domains = append(domains, env[1])
			}
		}
	}

	return domains
}
