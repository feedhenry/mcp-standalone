## Mobile Server

The mobile server is PoC for a per Namespace service that helps developers create and integrate mobile applications on OpenShift.


### Install

To install you will need a running OpenShift. We recommend using ```oc cluster up``` find out more [here](https://github.com/openshift/origin/blob/master/docs/cluster_up_down.md)

Once you have a running openshift you can use the install script:

```
oc new-project myproject
cd install/openshift
./install.sh <namespace>
```

### Local Development

- install go

You will go installed at version 1.7.x or later [goloang download](https://golang.org/dl/)

- install dep

You will also want to have dep installed for adding new dependencies [Install dep](https://github.com/golang/dep#setup)

- clone the repo into your ```$GOPATH```

```
mkdir -p $GOPATH/src/github.com/feedhenry
cd $GOPATH/src/github.com/feedhenry
git clone git@github.com:feedhenry/mcp-standalone.git
cd mcp-standalone
```

- Use make to run the tests and create a docker image

```
make test image
```

