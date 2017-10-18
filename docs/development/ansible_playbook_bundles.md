# FeedHenry Ansible Playbook Bundle Development Guide

The document is documenting best practises for the FeedHenry APB development! Some general information about the Ansible Service Broker (ASB) & Ansbile Playbook Bundles (APB) can be found[here](https://docs.openshift.com/container-platform/3.6/architecture/service_catalog/ansible_service_broker.html)

## Installation 

The `apb` CLI can be installed locally, _however_ it is highly recommended using it wrapped in a Linux container. Below is an example on how to create an alias for the _wrapped_`apb` CLI:

```bash
alias apb='docker run --rm --privileged -v $PWD:/mnt -v $HOME/.kube:/.kube -v /var/run/docker.sock:/var/run/docker.sock -u $UID docker.io/feedhenry/apb'
```


## Creating an APB

Using the `apb init` command you can create a template for your APB, like:

```
apb init DOCKERORG/NAME_OF_YOUR_STUFF-apb 

```

_NOTE:_ It's a good practise to have the APB end with `-apb`!


### Create a Makefile

Once you have made the changes for (de)provisioning or (un)binding, you _could_ use `apb prepare` and `apb build` to run the build. Using a `Makefile` is a much better idea. Currently our FeedHenry APBs do have a `Makefile`, like [here](https://raw.githubusercontent.com/feedhenry/fh-sync-server-apb/master/Makefile):

```
DOCKERHOST = docker.io
DOCKERORG = feedhenry
USER=$(shell id -u)
PWS=$(shell pwd)
build_and_push: apb_build docker_push

.PHONY: apb_build
apb_build:
	docker run --rm -u $(USER) -v $(PWD):/mnt:z feedhenry/apb prepare
	docker build -t $(DOCKERHOST)/$(DOCKERORG)/my-cool-apb .

.PHONY: docker_push
docker_push:
	docker push $(DOCKERHOST)/$(DOCKERORG)/my-cool-apb
```

### Local Build

With the `Makefile` in place, run a _local_ build like:

```
make apb_build
```

### Pushing to Dockerhub

If you want to push the _APB_ to your own Dockerhub account, run the following, to run a _complete_ build:

```
make DOCKERORG="my_org"
```

_NOTE:_ With `DOCKERHOST` you can also override the dockerhost; default is `docker.io`.

### Locally pushing APBs to Ansible Service Broker (ASB)

#### Setup of MCP

Create a MCP cluster pointing at your own docker organisation, this will also copy all the existing APBs in feedhenry to the dockerhub_org:
```
ansible-playbook playbook.yml -e "dockerhub_username=<dockerusername>" -e "dockerhub_password=<dockerpassword>" -e "dockerhub_org=<USE_THIS_VALUE>" -e "apb_sync=true" --ask-become-pass
```

#### Push to ASB

Once you have issued the (local) build, you can push your ABP to the running ASB of your MCP projct:

```
apb push
```

Afterwards your APB is ready to be used from the _Service Catalog_.

## Best Practises

### ServiceName label

Inside of the `apb.yml` file, make sure you use the `serviceName:` label, like:

```
...
metadata:
  displayName: FeedHenry Sync Server
  console.openshift.io/iconClass: font-icon icon-nodejs
  serviceName: fhsync
...
```

### Secret label

Inside of the Kubernetes/Openshift secrets, it's also recommended to use labels:

```
mcp:enabled
```

## Testing 

There are currently not much tools for testing. The APB team has a few open Github issues and proposals:

* [lint command for APB content](https://github.com/ansibleplaybookbundle/ansible-playbook-bundle/issues/131)
* [CI and Image tests](https://github.com/ansibleplaybookbundle/ansible-playbook-bundle/blob/master/docs/proposals/testing.md)
