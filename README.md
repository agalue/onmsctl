# onmsctl [![Go Report Card](https://goreportcard.com/badge/github.com/OpenNMS/onmsctl)](https://goreportcard.com/report/github.com/OpenNMS/onmsctl)

A CLI tool for OpenNMS.

The following features have been implemented:

* Manage multiple OpenNMS servers
* Verify installed OpenNMS Version
* Manage provisioning requisitions (replacing `provision.pl`)
* Manage SNMP configuration (replacing `provision.pl`)
* Manage Foreign Source definitions
* Send events to OpenNMS (replacing `send-event.pl`)
* Reload configuration of OpenNMS daemons
* Enumerate collected resources and metrics (replacing `resourcecli`)
* Manually manage the inventory (bypassing the provisioning system), useful when it is not possible to use Provisioning or Auto-Discover.
* Support for searching entities using [FIQL](https://fiql-parser.readthedocs.io/en/stable/usage.html) (work in progress)

The reason for implementing a CLI in `Go` is that the generated binaries are self-contained, and for the first time, Windows users will be able to control OpenNMS from the command line. For example, `provision.pl` or `send-events.pl` rely on having Perl installed with some additional dependencies, which can be complicated in the environment where this is either hard or impossible to have.

## Compilation

Make sure to have [GoLang](https://golang.org/dl/) installed on your system. Recommended version: 1.16 or newer.

To compile the source code for your desired operating system:

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

There is no need to specify parameters for your own operating system, as `go build` will be sufficient. You can also build targets for any operating system from any operating system, and the generated binary will work on itself; there is no need to install anything on the target device besides copying the generated binary file.

Alternatively, in case you don't want to install `Go` on your system, but you have [Docker](https://www.docker.com) installed, you can use it to compile it:

```bash
➜ docker run -it --rm -e GO111MODULE=on -e GOOS=windows -e GOARCH=amd64 -v $(pwd):/app golang bash
root@3854e5d2d67c:/go# cd /app
root@3854e5d2d67c:/app# go build -o onmsctl.exe
root@3854e5d2d67c:/app# exit
```

## Usage

The binary contains help for all commands and subcommands by using `-h` or `--help`. Everything should be self-explanatory.

The following outlines several examples.

0. Deploy OpenNMS

There are several ways to do that, and I'm going to use Kind and Helm:

```bash
kind create cluster --name onms
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add opennms https://opennms.github.io/helm-charts
helm install onms-db bitnami/postgresql --version 15.5.38 \
  --set global.postgresql.auth.postgresPassword=P0stgr3s
kubectl rollout status sts/onms-db-postgresql
helm install onms opennms/horizon \
  --set domain=test.local \
  --set dependencies.postgresql.hostname=onms-db-postgresql.default.svc \
  --set dependencies.postgresql.sslmode=disable
kubectl rollout status sts/onms-core
kubectl port-forward svc/onms-core 8980
```

> That will take a few minutes to be ready.

Your OpenNMS server would be reachable at `http://localhost:8980/opennms`.

Once you're done, use `ctrl+c` to stop the port forward and then remove the cluster:

```bash
kind delete cluster -n onms
```

1. Configure OpenNMS servers

To configure the tool, or to avoid specifying the URL, username, and password for your OpenNMS server with each request, you can create a file with the following content on `$HOME/.onms/config.yaml` or add the file on any location and create an environment variable called `ONMSCONFIG` with the location of the file.

To manipulate the content of the configuration file, please use the `onmsctl config` subcommand.

For instance, to add a new configuration entry and make it the default:

```bash
➜  onmsctl config set --name M2021 --url http://192.168.205.200:8980/opennms --user admin --passwd admin
➜  onmsctl config default M2021
➜  onmsctl config list
Default	Name		User	URL
*	M2021		admin	http://192.168.205.200:8980/opennms
	local		admin	http://localhost:8980/opennms

➜  onmsctl info
displayVersion: 2021.1.1
version: 2021.1.1
packageName: meridian
packageDescription: OpenNMS Meridian
datetimeFormat:
  zoneId: America/New_York
  format: yyyy-MM-dd'T'HH:mm:ssxxx
```

> Make sure to protect the file, as the credentials are on plain text.

2. Verify the installed version of OpenNMS

```bash
onmsctl info
```

The output would be something like this:

```
displayVersion: 26.2.2
version: 26.2.2
packageName: opennms
packageDescription: OpenNMS
datetimeFormat:
  zoneId: America/New_York
  format: yyyy-MM-dd'T'HH:mm:ssxxx
```

3. Build a requisition like you would do it with `provision.pl`:

```bash
➜ onmsctl inv req add Local
➜ onmsctl inv node add Local srv01
➜ onmsctl inv intf add Local srv01 10.0.0.1
➜ onmsctl inv svc add Local srv01 10.0.0.1 ICMP
➜ onmsctl inv cat add Local srv01 Servers
➜ onmsctl inv asset set Local srv01 address1 home
```

To visualize the content of the requisition:

```bash
➜ onmsctl inv node get Local srv01
nodeLabel: srv01
foreignID: srv01
interfaces:
- ipAddress: 10.0.0.1
  snmpPrimary: "N"
  status: 1
  services:
  - name: ICMP
categories:
- name: Servers
assets:
- name: address1
  value: home
```

To import the requisition:

```bash
➜ onmsctl inv req import Local
```

4. Build requisitions in `YAML` and apply it (similar to `kubernetes` workload with `kubectl`):

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
➜ cat <<EOF | onmsctl inv node apply -f=- Local
 foreignID: www.opennms.com
 interfaces:
 - ipAddress: www.opennms.com
 categories:
 - name: WebSites
EOF
```

Or,

```bash
➜ onmsctl inv node apply Local '
foreignID: www.opennms.com
interfaces:
- ipAddress: www.opennms.com
categories:
- name: WebSites'
```

> Note that an FQDN was used instead of an IP Address (more on this below).

To get the actual content of the node:

```bash
➜ onmsctl inv node get Local www.opennms.com
nodeLabel: www.opennms.com
foreignID: www.opennms.com
interfaces:
- ipAddress: 141.193.213.20
  snmpPrimary: "N"
  status: 1
categories:
- name: WebSites
```

As you can see, it is possible to specify FQDN instead of IP addresses, and they will be translated into IPs before sending the JSON payload to the ReST end-point for requisitions.

Additionally, for convenience, if the `node-label` is not specified, the `foreign-id` will be used.

5. Configure SNMP credentials

Obtain the current credentials for a given IP address:

```bash
➜ onmsctl snmp get 12.0.0.1
```

> For nodes behind Minions, you can specify the location as a command option.

The output would be:

```
version: v2c
port: 161
retries: 1
timeout: 1800
community: public
maxRequestSize: 65535
maxRepetitions: 2
maxVarsPerPdu: 10
```

Change the credentials for a given IP address:

```
➜ onmsctl snmp set -v v2c -r 3 -t 2000 -c c0mpl1x 12.0.0.1
```

> For nodes behind Minions, you can specify the location as a command option.

## Upcoming features

* Visualize tabular data with pagination (nodes, events, alarms, outages, notifications).

* Configure scheduled outages.
