## Mobile Server

The mobile server is PoC for a per Namespace service that helps developers create and integrate mobile applications on OpenShift.


### Install

To install you will need a running OpenShift. We recommend using ```oc cluster up``` find out more [here](https://github.com/openshift/origin/blob/master/docs/cluster_up_down.md)

Once you have a running openshift you can use the install script:

```
oc new-project myproject
cd install/openshift
./install.sh
````