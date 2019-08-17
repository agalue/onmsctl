# onmsctl

A CLI tool for OpenNMS.

The following features have been implemented:

* Verify installed OpenNMS Version
* Manage provisioning requisitions (replacing `provision.pl`)
* Manage SNMP configuration (replacing `provision.pl`)
* Manage Foreign Source definitions
* Send events to OpenNMS (replacing `send-event.pl`)
* Reload configuration of OpenNMS daemons
* Enumerate collected resources and metrics (replacing `resourcecli`)
* Preliminar support for searching entities (work in progress)

The reason for implementing a CLI in `Go` is that the generated binaries are self-contained, and for the first time, Windows users will be able to control OpenNMS from the command line. For example, `provision.pl` or `send-events.pl` rely on having Perl installed with some additional dependencies, which can be complicated on the environment where this is either hard or impossible to have.

## Compilation

1. Make sure to have [GoLang](https://golang.org/dl/) installed on your system.

2. Make sure to have `Go Modules` enabled (recommended version: 1.12 or newer)

```bash
export GO111MODULE=on
```

3. Compile the source code for your desired operating system

For Linux:

```bash
GOOS=linux GOARCH=amd64 go build -o onmsctl onmsctl.go
```

For Mac:

```bash
GOOS=darwin GOARCH=amd64 go build -o onmsctl onmsctl.go
```

For Windows:

```bash
GOOS=windows GOARCH=amd64 go build -o onmsctl.exe onmsctl.go
```

For your own operating system, there is no need to specify parameters, as `go build` will be sufficient. Also, you can build targets for any operating system from any operating system, and the generated binary will work on itself, there is no need to install anything on the target device, besides copying the generated binary file.

Alternatively, in case you don't want to install `Go` on your system, but you have [Docker](https://www.docker.com) installed, you can use it to compile it:

```bash
➜ docker run -it --rm -e GO111MODULE=on -e GOOS=windows -e GOARCH=amd64 -v $(pwd):/app golang:1.12 bash
root@3854e5d2d67c:/go# cd /app
root@3854e5d2d67c:/app# go build -o onmsctl.exe
root@3854e5d2d67c:/app# exit
```

## Usage

The binary contains help for all commands and subcommands by using `-h` or `--help`. Everything should be self-explanatory.

1. Build a requisition like you would do it with `provision.pl`:

```bash
➜ onmsctl inv req add Local
➜ onmsctl inv node add Local srv01
➜ onmsctl inv intf add Local srv01 10.0.0.1
➜ onmsctl inv svc add Local srv01 10.0.0.1 ICMP
➜ onmsctl inv cat add Local srv01 Servers
➜ onmsctl inv assets set Local srv01 address1 home
➜ onmsctl inv node get Local srv01
nodeLabel: srv01
foreignID: srv01
interfaces:
- ipAddress: 10.0.0.1
  snmpPrimary: S
  status: 1
  services:
  - name: ICMP
categories:
- name: Servers
assets:
- name: address1
  value: home

➜ onmsctl inv req import Local
Importing requisition Local (rescanExisting? true)...
```

2. You can build requisitions in `YAML` and apply it like `kubernetes` workload with `kubectl`:

```bash
➜ cat <<EOF | onmsctl inv req apply -f -
name: Routers
nodes:
- foreignID: router01
  nodeLabel: Router-1
  interfaces:
  - ipAddress: 10.0.0.1
  categories:
  - name: Routers
- foreignID: router02
  nodeLabel: Router-2
  interfaces:
  - ipAddress: 10.0.0.2
  categories:
  - name: Routers
EOF
```

The above also works for individual nodes:

```bash
➜ cat <<EOF | onmsctl inv node apply -f - Local
foreignID: www.opennms.com
interfaces:
- ipAddress: www.opennms.com
categories:
- name: WebSites
EOF

www.opennms.com translates to [34.194.50.139], using the first entry.
Adding node www.opennms.com to requisition Local...

➜ onmsctl inv node get Local www.opennms.com
nodeLabel: www.opennms.com
foreignID: www.opennms.com
interfaces:
- ipAddress: 34.194.50.139
  snmpPrimary: N
  status: 1
categories:
- name: WebSites
```

As you can see, it is possible to specify FQDN instead of IP addresses, and they will be translated into IPs before sending the JSON payload to the ReST end-point for requisitions.

Additionally, for convenience, if the `node-label` is not specified, the `foreign-id` will be used.

To configure the tool, or to avoid specifying the URL, username and password for your OpenNMS server with each request, you can create a file with the following content on `$HOME/.onms/config.yaml` or add the file on any location and create an environment variable called `ONMSCONFIG` with the location of the file:

```yaml
url: demo.opennms.com
username: demo
password: demo
```

Make sure to protect the file, as the credentials are on plain text.

## Upcoming features

* Search for entities. The idea is to provide a way to build a search expression that will be translated into a [FIQL](https://fiql-parser.readthedocs.io/en/stable/usage.html) expression and use the ReST API v2 of OpenNMS to search for events, alarms, nodes, etc.

* Visualize tabular data with pagination (nodes, events, alarms, outages, notifications).
