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

__Note:__ On Linux you also need to have the `libselinux-python` package installed!

Execute these commands to clone the repo to the correct location.
```sh
mkdir -p ~/go/src/github.com/feedhenry/mcp-standalone && cd ~/go/src/github.com/feedhenry/mcp-standalone
git clone git@github.com:<YOUR_FORK>/mcp-standalone.git .
```

If you don't already have a Go environment setup you will need to add it to the path:
```sh
export PATH="$PATH:~/go/bin"
```
You will want to add the path permanently to your `.bashrc` or `.bashprofile`.

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

#### Docker Firewall setup

Next we need to configure Docker registry _and_ ports required as part of the cluster setup:
 * Linux: Follow steps 2 _and_ 3 [here](https://github.com/openshift/origin/blob/master/docs/cluster_up_down.md#linux)
 * Mac: Follow steps 2 _and_ 3 [here](https://github.com/openshift/origin/blob/master/docs/cluster_up_down.md#macos-with-docker-for-mac)

For Linux we also need to add an extra port to the `dockerc` zone:
```sh
firewall-cmd --permanent --zone dockerc --add-port 443/tcp
firewall-cmd --reload
```

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
./node_modules/.bin/grunt local
```
If you see an `ENOSPC` error, you may need to increase the number of files your user can watch by running this command:

```sh
echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p
```

*NOTE*: Running `./node_modules/.bin/grunt local` will *not* run `uglify` (to help with local dev), and *will* include `scripts/config.local.js`. This file is used to point to a local running MCP server rather than the default of looking up a Route names `mcp-standalone` and using that as the MCP server host.

### Step 3, Launch Services

Open your browser and point it at:
```
https://192.168.37.1:8443/console/
```

#### Step 3.1, Add a mobile app

In your project (the `localmcp` project), click on the mobile tab located in the left hand navigation. 

##### Provision the Mobile Control Panel

**Note:** If you the `Provision Mobile Control Panel` screen, when doing _Local Development_, it is due to an issue with the self-signed certificate of the locally running MCP server. You will get around it by accepting the certificate in the browser, by hitting `https://127.0.0.1:3001`!

##### Creating a Mobile Application

Next click ```Create Mobile App``` fill in some details and pick cordova as the type. Then click create.

**NOTE:** _You will need to accept the cert of the mobile server. You can do this by hitting https://127.0.0.1:3001 in your browser._  It's also recommended to open the browsers developer tools for more infos, in case there are problems with the self-signed certificate, or other issues.

*NOTE*: The cert can get out of sync at times. If you see a bad cert error you can clear the cert cache.

Linux:
```sh
rm ~/.pki/nssdb/cert9.db
```
And then restart your browser.

#### Step 3.2, Launch Keycloak

From the Service Catalog, select `Keycloak (APB)`, and either enter values for the username and password or accept the defaults, select `localmcp` as the project to add it to, then click `next`.

Select `Do not bind at this time` and click `Create` and then `Close`.

#### Step 3.2, Launch FeedHenry Sync

Select `FeedHenry Sync (Persistent)`, choose `localmcp` as the project and click `next`.

Select `Create a secret in localmcp to be used later ` and click `Create` and then `View Project`.

Wait for all the pods to have come up, before proceeding (this can take a few minutes and is indicated by a blue hoop next to each pod name).

### Step 4, Mount Keycloak secret into Feedhenry Sync

Log in to the web console for Openshift and click into the `localmcp` project, then select `Mobile`. Under `Mobile Enabled Services` select `fh-sync-server` and click `Integrations`. 
Under `Mobile Service Integrations`, click on `Create Integration` next to `Keycloak`.

