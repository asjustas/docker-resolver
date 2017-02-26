# Docker resolver - Docker DNS resolver

[![GoDoc](https://godoc.org/github.com/asjustas/docker-resolver?status.png)](https://godoc.org/github.com/asjustas/docker-resolver)

A simple DNS server + `/etc/hosts` file updater used to resolve names of local Docker containers.
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