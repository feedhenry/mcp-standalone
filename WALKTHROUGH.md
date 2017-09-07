# Mobile Control Panel Demonstration
This document is intended to walk you through setting up a local openshift development cluster, deploying two mobile services (in particular `Feedhenry Sync` and `Keycloak`) to it via the ASB then executing the mobile control panel server configured to communicate with this cluster to connect those two services together.

## Requirements
- `oc` tool [installed](https://github.com/openshift/origin/releases/tag/v3.6.0).
- `go` programming language [installed](https://golang.org/dl/).
- `ansible-playbook` tools [installed](http://docs.ansible.com/ansible/latest/intro_installation.html)
- Local clone of this repo

## Step 1, Creating a Local Cluster
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

## Step 2, Compile the MCP Server
We are now ready to compile the MCP Server so that we can execute against our new cluster. Compiling this is easy with the Makefile, this should be executed in the root directory of this repo:
```sh
make run_server NAMESPACE=demoproject
```

## Step 3, Launch Services
Open your browser and point it at:
```
https://192.168.37.1:8443/console/
```

### Step 3.1, Launch Keycloak
Select `Keycloak (APB)`, and either enter values for the username and password or accept the defaults, select `demoproject` as the project to add it to, then click `next`.

Select `Create a secret in demoproject to be used later` and click `Create` and then `Close`.

### Step 3.2, Launch FeedHenry Sync
Select `FeedHenry Sync (Persistent)` and either enter some values here, or leave them on defaults, select `demoproject` as the project and click `next`.

Select `Create a secret in demoproject to be used later ` and click `Create` and then `View Project`.

Wait for all the pods to have come up, before proceeding (this can take a few minutes and is indicated by a blue hoop next to each pod name).

## Step 4, Mount Keycloak secret into Feedhenry Sync
Now we will make a POST request to the end-point of the MCP server, and it will configure FeedHenry Sync to be able to make use of Keycloak.

```sh
curl -k \
  -H "Authorization: $(oc whoami -t | tr -d '[[:space:]]')" \
  -H "Content-Type: Application/JSON" \
  -X POST \
  https://127.0.0.1:3001/mobileservice/configure/fh-sync-server/keycloak-public-client
```

The response to this should be a JSON block. 

## Step 5, Confirm Keycloak secret is mounted in FeedHenry Sync
If you look in the OpenShift Console for the `demoproject` now, you will see a fresh deploy of `FeedHenry Sync` has been triggered. Once this is running, connect to the remote shell for the new pod:
```sh
oc rsh <fh-sync-server-pod-name> cat /etc/secrets/keycloak-public-client/installation
```

This shows the fh-sync-server now has access to the data required to connect to the keycloak service.