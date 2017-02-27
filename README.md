# Docker resolver - Docker DNS resolver

[![GoDoc](https://godoc.org/github.com/asjustas/docker-resolver?status.png)](https://godoc.org/github.com/asjustas/docker-resolver)

A simple DNS server + `/etc/hosts` file updater used to resolve names of local Docker containers.

This app listens to docker events and automatically updates your `/etc/hosts` file to allow you easy access running containers.

Also you can configure docker daemon to use build in dns server. This can be used to communicate with other docker container

with a known port bound to the Docker bridge using domain names.  

## Container Registration

`docker-resolver` uses `hostname`, `container name` and `DOMAIN_NAME`, `DNSDOCK_ALIAS` env variables to register containers.

For example, the following container would be available as: 

* `container.docker`
* `container.demo`
* `container.dev`
* `container.test`
* `container.io`:

```yml
symfony:
    container_name: container
    hostname: container.demo
    build: docker/web
    volumes:
        - .:/var/www/html
    environment:
        DNSDOCK_ALIAS: container.dev,container.test
        DOMAIN_NAME: container.io
```

## Installation

```sh
# Compile application
go build

# Move file to /usr/bin/
sudo cp docker-resolver /usr/bin/
```

## Systemd service installation

```sh
# Copy service file to systemd
sudo cp systemd/docker-resolver.service /etc/systemd/system/

# Restart systemd configurations
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable docker-resolver

# Start service
sudo systemctl start docker-resolver
```