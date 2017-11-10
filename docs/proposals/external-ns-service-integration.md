#External Service Integrations

## Background
https://github.com/kubernetes-incubator/service-catalog/issues/1455

## Abstract

This document aims to provide solutions to the following use case: 

A mobile enabled service has been provisioned via the service catalog into namespace A. It is intended for this to be a shared service. Think something like
Authentication. 
An MCP is provisioned in namespace B along with another mobile enabled service that is capable of integrating with the mobile service in namespace A. How do we support
this kind of integration?


## Cross Namespace integrations

What we want to achieve is for the ansible playbook to be informed that there is an existing service that it should use for the binding that is not in the target namespace.

We are leveraging the machinery provided by the service catalog to power our integrations. This document assumes some level of familiarity with [Bindings](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#binding) and [Provision ServiceInstances](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#provisioning)and the [ansible service broker](https://github.com/openshift/ansible-service-broker).

Currently a binding can only be created in the same namespace as an existing service instance:
https://github.com/kubernetes-incubator/service-catalog/issues/1455

The end result of a successful binding is a secret, that contains a set of credentials and is created within the same namespace as the service instance.
The broker is responsible for doing any work for a binding and then returning any credentials back to be injected into the secret. 

Until the issue, outlined in the background, is resolved, we will also need a to create a service instance in any namespace where we want to provide a binding.
We will need to indicated to the playbook that there should not be a newly provisioned service created but rather that this provision request refers to the existing service in a different namespace.
This means that the playbook will only need to perform actions that would happen after the service has been sucesfully provisioned. For instance, in the our case, we would create a new secret
in the target namespace that represents the service just as we would have the first time it was provisioned.    

### Options

As the service integration is to an external service, the service instance will need to be created in the MCP project but reflect the external service. 

-  We already have an option for adding an external service and specifying which namespace it is in. If the user adds the required credentials here, these can be used as part of the 
integration process to create the service instance and the binding. In theory, we could offer a "import from namespace" option, that would look at a specified namespace and import any of the 
known service secrets into the current namespace. 

This process would work as follows:
- User provisions a service (keycloak for example) from the catalog into a separate namespace.
- In a namespace with an MCP, the user adds an external service via the MCP UI and imports or adds the required credentials for the keycloak they provisioned earlier.
- User now chooses to integrate a service in the same namespace as MCP with this external service (fh-sync for example)
- MCP sees this is an external service in another namespace
- MCP reads the credentials created when the external service was added
- MCP creates a provision service instance request passing the credentials as parameters to the service catalog (there is support for this already)
- In the APB provision playbook, when it sees these credentials, it will follow a different flow: no longer creating pods etc but perhaps reading some data from the service and creating any 
required objects that are not new pods services routes etc (for example a secret that could be used if an MCP where deployed to this namespace at a later stage)
- MCP next makes a binding call to service catalog, and the provision credentials will be passed on by the broker as part of the the binding params. 
The binding does what is needed except now it is making requests against a service in a different namespace. The end result is any required info from the binding is setup as a new secret in the MCP namespace ready to be consumed.
- MCP creates a pod preset (as it already does) to allow the target pods to consume this information.
 