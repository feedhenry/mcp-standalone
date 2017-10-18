## Local Walkthrough

This document will walk you through setting up the Mobile Control Panel on a local OpenShift cluster using ```oc cluster up```. This is targeted at someone wanting to try out the Mobile Control Panel rather than develop it.

### Prerequisites 

- [Ansible](http://docs.ansible.com/ansible/latest/intro_installation.html) >= 2.3.2.0
- [Docker Hub account](https://hub.docker.com/)
  - This is needed by the Ansible Service Broker to list APBs. The upstream APBs for the MCP are kept in the `feedhenry` organisation.
- [Docker](https://docs.docker.com/engine/installation/)
- [oc command line client](https://github.com/openshift/origin/releases)
- [Node.js](https://nodejs.org/en/) >= 6.10.0

### Setup

1) Clone this repository

```bash
git clone git@github.com:feedhenry/mcp-standalone.git
```

**Note:** it is a good idea to clone this into a valid $GOPATH, however it is not essential.

2) Install the required ansible dependencies: 
```sh
ansible-galaxy install -r ./installer/requirements.yml
```

3) Run the ansible installer. This installer sets up your OpenShift environment with the service catalog and the ansible service broker. 

```sh
export DOCKERHUB_USERNAME="<username>"
export DOCKERHUB_PASSWORD="<password>"
export DOCKERHUB_APBS_ORG="feedhenry"
cd installer/ && ansible-playbook playbook.yml -e "dockerhub_username=$DOCKERHUB_USERNAME" -e "dockerhub_password=$DOCKERHUB_PASSWORD" -e "dockerhub_org=$DOCKERHUB_APBS_ORG" --ask-become-pass
```

### Creating Mobile Apps

Once the installer is complete you should be able to access OpenShift at [https://192.168.37.1:8443/console/](https://192.168.37.1:8443/console/). You will need to accept the self-signed certificate.

You can login using `developer` and any password.

Once logged in you should be presented with a "catalog" of services. To create your first Mobile App:

- Choose the "Mobile" category & "Apps" sub-category.
- Choose "Android App"
- Fill in the required information and continue through the wizard.
  - This will provision the MCP Server (first time only) and create the Android App. The MCP Server is required before Mobile Apps can be created.
  - Any further Apps created from the Catalog will use the same MCP Server.
- You'll need to accept the self-signed cert for the MCP Server in your Browser. To do this:
  - Get the route for the MCP Server by running:
    - `oc get route mcp-standalone -n myproject --template "https://{{.spec.host}} "`
  - Visit the route in your browser and accept the cert.
  - The page might give a message like 'no token provided access denied'. This is OK as it show's the cert is now trusted and we're hitting the server.
- Back in your Project, click the "Mobile" tab on the left nav. You should now see the Mobile Overview screen and your Android App.


