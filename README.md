# MCP (mobile control panel) Standalone 

The MCP Standalone is PoC for a per Namespace service that helps developers create and integrate mobile applications on OpenShift.


## Local Development

This document is intended to walk you through setting up a local openshift development cluster, deploying two mobile services (in particular `Feedhenry Sync` and `Keycloak`) to it via the ASB then executing the mobile control panel server configured to communicate with this cluster to connect those two services together.

### Requirements

- `oc` tool [installed](https://github.com/openshift/origin/releases/tag/v3.6.0).
- `go` programming language [installed](https://golang.org/dl/).
- `ansible-playbook` tools [installed](http://docs.ansible.com/ansible/latest/intro_installation.html)
- Local clone of this repo

### Step 1, Creating a Local Cluster

First we will use the ansible-playbooks included in this repo to create a local oc cluster which is running the Ansible Service Broker. 

First check that oc cluster is down:
```sh
oc cluster down
```

The next step is executed from inside the `installer` directory in this repo:
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
npm i && bower install && grunt local
```

*NOTE*: Running `grunt local` will *not* run `uglify` (to help with local dev), and *will* include `scripts/config.local.js`. This file is used to point to a local running MCP server rather than the default of looking up a Route names `mcp-standalone` and using that as the MCP server host.

### Step 3, Launch Services

Open your browser and point it at:
```
https://192.168.37.1:8443/console/
```

#### Step 3.1, Launch Keycloak

From the Service Catalog, select `Keycloak (APB)`, and either enter values for the username and password or accept the defaults, select `localmcp` as the project to add it to, then click `next`.

Select `Do not bind at this time` and click `Create` and then `Close`.

#### Step 3.2, Launch FeedHenry Sync

Select `FeedHenry Sync (Persistent)`, choose `localmcp` as the project and click `next`.

Select `Create a secret in localmcp to be used later ` and click `Create` and then `View Project`.

Wait for all the pods to have come up, before proceeding (this can take a few minutes and is indicated by a blue hoop next to each pod name).

### Step 4, Mount Keycloak secret into Feedhenry Sync

TBD

