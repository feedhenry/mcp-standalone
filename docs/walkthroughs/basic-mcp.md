## Basic MCP (Mobile Control Panel) Walkthrough

This document will walk you through setting up the Mobile Control Panel on a local OpenShift cluster using ```oc cluster up```. This is targeted at someone wanting to try out the Mobile Control Panel rather than develop it.

### Prerequisites 

- Ansible version 2.3.2.0 or greater
- Docker
- OC command line tool

### Setup

1) Clone this repository  ```git clone git@github.com:feedhenry/mcp-standalone.git ``` *note* it is a good idea to clone this into a valid $GOPATH, however it is not essential.

2) Install the required ansible dependencies: 
```sh
sudo ansible-galaxy install -r requirements.yml
```

3) Run the ansible installer. This installer sets up your OpenShift environment with the service catalog and the ansible service broker.

```sh
cd installer/ && ansible-playbook playbook.yml -e "dockerhub_username=<your docker login>" -e "dockerhub_password=<your docker password>" -e "dockerhub_org=<docker org containing APBs>" --ask-become-pass
```

### Deploying MCP via the catalog

Once the installer is complete you should be able to access OpenShift at [https://192.168.37.1:8443/console/](https://192.168.37.1:8443/console/) 

You can login using the username ```developer``` and pass ```anypass```

Once logged in you should be presented with a "catalog" of services. You should be able to see a ```mobile``` cataogory. Click on this and select apps.

- Choose android app
- Fill in the required information and click ok.
- Jump into your namespace you should see a new route created. Open this route in your browser and accept the cert.
- Back in  your namespace, click on the mobile tab on the left hand side nav, you should now be able to see your android app.


