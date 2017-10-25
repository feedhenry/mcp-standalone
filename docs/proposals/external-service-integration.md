#External Service Integrations

## Background
https://github.com/kubernetes-incubator/service-catalog/issues/1455

## Abstract

This document aims to provide solutions to the following use cases: 

1) A mobile enabled service has been provisioned via the service catalog into namespace A. It is intended for this to be a shared service. Think something like
Authentication. 
An MCP is provisioned in namespace B along with another mobile enabled service that is capable of integrating with the mobile service in namespace A. How do we support
this kind of integration?


2) A user or org has an existing service such as Keycloak that exists entirely separately to OpenShift, but is made accessible. It is the intention that this would be a
shared service. A developer then creates a mobile project with MCP provisioned on an OpenShift instance, in what circumstances can we provide integration and what form
should that integration take?


## Cross Namespace integrations

We are leveraging the machinery provided by the service catalog to power our integrations. This document assumes some level of familiarity with ServiceBindings and ServiceInstances and the ansible service broker.

Currently a binding can only be created in the same namespace as an existing provisioned a service instance:
https://github.com/kubernetes-incubator/service-catalog/issues/1455

The end result of a successful binding is a secret, that contains a set of credentials that is then created within the same namespace as the service instance.
The broker is responsible for doing any work for a binding and then returning any credentials back to be injected into the secret. What we want to achieve is for the ansible play book to be informed that there is an existing 
service that it should use for the binding that is not in its current namespace. Until the issue, outlined in the background, is resolved, we will also need a way to create a service instance
in the namespace that does not trigger a brand new instance of the service from the broker but rather refers to the existing service in a different namespace.   

### Options

As the service integration is to an external service, the service instance will need to be created in the MCP project but reflect the external service. 

-  We already have an option for adding an external service and specifying which namespace it is in. If the user adds the required credentials here, these can be used as part of the 
integration process to create the service instance and the binding. In theory, we could offer a "import from namespace" option, that would look at a specified namespace and import and of the 
known service secrets into the current namespace. 
This process would work as follows:
- User adds external service via the UI add the required credentials
- User chooses to integrate a service in the same namespace as MCP with this external service
- MCP sees this is an external service in another namespace
- MCP reads the credentials created when the external service was added
- MCP creates a provision service instance request passing the credentials as parameters to the service catalog (there is support for this already)
- Service APB provison, when it sees these credentials follows a different flow, no longer creating pods etc but perhaps reading some data from the service and creating any 
required objects that are not new pods services routes etc
- MCP next makes a binding call to service catalog, as per the broker spec, the provision credentials will be passed on by the broker to the binding. The binding does what is needed
except now it is making requests against a service in a different namespace. The end result is any required info from the binding is setup as a new secret in the MCP namespace ready to be consumed.
 