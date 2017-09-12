# Mobile Control Panel (MCP)

The Mobile Control Panel is PoC for a 'per namespace' service that helps developers discover, create and integrate Mobile Apps and Services on OpenShift.

* Mobile SDKs are developed and maintained in their respective repos
* The Backend is a Golang server in this repo
* The Frontend is a set of AngularJS services, controllers, views etc... in this repo that extend the OpenShift Web Console UI (via extensions)
* Services are developeed and maintained in their respective repos. They leverage the Service Catalog and various brokers to help provision them

The MCP brings all of these componments together to create a unified Mobile developer experience on top of OpenShift.

## Contributing

You can develop [locally on your host](#local-development).
Please include as much info as possible in Issues and Pull Requests.
Merging to master requires approval from a reviewer and a passing CI build.

## Communication

Daily communication happens on #feedhenry on [freenode IRC](https://webchat.freenode.net/).
The [feedhenry-dev@redhat.com mailing list](http://feedhenry-dev.2363497.n4.nabble.com/) is also used for team-wide & community comms.
Issues are tracked in both [Jira](https://issues.jboss.org/secure/RapidBoard.jspa?rapidView=4143&view=planning.nodetail) and Github Issues. Where issues are duplicates, they should be linked so that only 1 source of info exists (automation would be nice here). Typically the core Red Hat team will create and work from Jira Issues.

## Onboarding Resources

* [Local Development](#local-development)
* Mobile SDKs
  * [Android Sync SDK](https://github.com/feedhenry/fh-sync-android)
  * [Cordova/Browser Sync SDK](https://github.com/feedhenry/fh-sync-js)
  * [Push SDKs](https://www.aerogear.org/docs/specs/#push)
  * [Keycloak JS Adapter](https://www.npmjs.com/package/keycloak-js)
* Backend Resources
  * [Tour of Go](https://tour.golang.org/welcome/1)
* Frontend Resources
  * [UI src](https://github.com/feedhenry/mcp-standalone/tree/master/ui)
  * [AngularJS PhoneCat Tutorial](https://docs.angularjs.org/tutorial)
  * [AngularJS API Docs](https://docs.angularjs.org/api)
  * [Patternfly](http://www.patternfly.org/)
  * [OpenShift Web Console](https://github.com/openshift/origin-web-console)
  * [Customising the OpenShift Web Console (Extensions](https://docs.openshift.com/container-platform/3.6/install_config/web_console_customization.html)
  * [Service Catalog/OpenShift Mall UI](https://github.com/openshift/origin-web-catalog)
* Catalog/Mall, Brokers & Services
  * [Service Catalog](https://docs.openshift.com/container-platform/3.6/architecture/service_catalog/index.html)
  * [Ansible Service Broker (ASB) & Ansbile Playbook Bundles (APB)](https://docs.openshift.com/container-platform/3.6/architecture/service_catalog/ansible_service_broker.html)
  * [Template Service Broker](https://docs.openshift.com/container-platform/3.6/architecture/service_catalog/template_service_broker.html)
  * [fh-sync-server](https://github.com/feedhenry/fh-sync-server)
  * [fh-sync-server Template](https://github.com/feedhenry/fh-sync-server/blob/master/fh-sync-server-DEVELOPMENT.yaml)
  * [Keycloak](https://github.com/keycloak/keycloak)
  * [Keycloak APB](https://github.com/feedhenry/keycloak-apb)

## Local Development

This document is intended to walk you through setting up a local openshift development cluster, deploying two mobile services (in particular `Feedhenry Sync` and `Keycloak`) to it via the ASB then executing the mobile control panel server configured to communicate with this cluster to connect those two services together.

### Requirements

- `oc` tool [installed](https://github.com/openshift/origin/releases/tag/v3.6.0).
- `go` programming language [installed](https://golang.org/dl/).
- `ansible-playbook` tools [installed](http://docs.ansible.com/ansible/latest/intro_installation.html)
- Local clone of this repo

Execute these commands to clone the repo to the correct location.
```sh
mkdir -p ~/go/src/github.com/feedhenry/mcp-standalone && cd ~/go/src/github.com/feedhenry/mcp-standalone
git clone git@github.com:<YOUR_FORK>/mcp-standalone.git .
export PATH="$PATH:~/go/bin"
```

### Setup the cli 

there is a very basic cli at ```cmd/mcp-cli``` you can build this and use it by running
``` make build_cli ``` this will drop a binary in your current dir which you can then use to exercise the api.


### Local Development
### Step 1, Creating a Local Cluster

First we will use the ansible-playbooks included in this repo to create a local oc cluster which is running the Ansible Service Broker. 

### Prerequisites

* A DockerHub account, credentials are required to set up the Ansible Service
Broker.
* User with sudo permissions on machine.

First check that oc cluster is down:
```sh
oc cluster down
```

Now we need to install any dependencies. The next step is executed from inside the `installer` directory in this repo:
```sh
sudo ansible-galaxy install -r requirements.yml
```

Next we need to configure Docker to accept an insecure registry required as part of the cluster setup.

**Linux**

Add to the file `/etc/docker/daemon.json` (create the file if it doesn't exist)
```json
{
    "insecure-registries" : [ "172.30.0.0/16" ]
}
```
Then restart Docker
```sh
sudo systemctl daemon-reload
sudo systemctl restart docker
```

**Mac**

Click the Docker icon in the tray to open Preferences. Click on the Daemon tab and add your insecure registries in Insecure registries section.
Don't forget to Apply & Restart and you're ready to go.

The next step is again executed from inside the `installer` directory:
```sh
cd installer/
ansible-playbook playbook.yml \
  -e "dockerhub_username=<your docker login>" \
  -e "dockerhub_password=<your docker password>" \
  -e "dockerhub_org=<docker org containing APBs>" \
  --ask-become-pass
```

This will set up your cluster for you - note that it is possible for this to fail on the first attempt, as the cluster up check may fail waiting for the images to be pulled - if this happens, re-run `oc cluster down` and execute the playbook again.

### Step 2, Run the MCP Server

We are now ready to compile the MCP Server so that we can execute against our new cluster. Compiling this is easy with the Makefile, this should be executed in the root directory of this repo:
```sh
cd ..
make run_server NAMESPACE=localmcp
```

In another terminal, bundle the MCP UI extension for OpenShift, watching for changes.
This is required to produce the mcp extension files referenced in master-config.yaml, and keep them up to date whenever changed during development.

```
cd ui
grunt local
```

*NOTE*: Running `grunt local` will *not* run `uglify` (to help with local dev), and *will* include `scripts/config.local.js`. This file is used to point to a local running MCP server rather than the default of looking up a Route names `mcp-standalone` and using that as the MCP server host.

### Step 3, Launch Services

Open your browser and point it at:
```
https://192.168.37.1:8443/console/
```

#### Step 3.1, Add a mobile app

In your project, click on the mobile tab located in the left hand navigation. Next click ```create app``` fill in some details and pick cordova as the type. Then click create. (note) you will need to accept the cert of the mobile server. You can do this by hitting https://localhost:3001 in your browser.

#### Step 3.2, Launch Keycloak

From the Service Catalog, select `Keycloak (APB)`, and either enter values for the username and password or accept the defaults, select `localmcp` as the project to add it to, then click `next`.

Select `Do not bind at this time` and click `Create` and then `Close`.

#### Step 3.2, Launch FeedHenry Sync

Select `FeedHenry Sync (Persistent)`, choose `localmcp` as the project and click `next`.

Select `Create a secret in localmcp to be used later ` and click `Create` and then `View Project`.

Wait for all the pods to have come up, before proceeding (this can take a few minutes and is indicated by a blue hoop next to each pod name).

### Step 4, Mount Keycloak secret into Feedhenry Sync

TBD

