## Code

This document will cover the design principals that should guide the server side component for the MCP.

The code development is guided as much as possible by the principals behind Onion architecture or Clean architecture:


- Dependency Inversion: dependencies should always point in to the layers below and never back out.
- Dependency injection: defining interfaces that can be swapped out at runtime for different concrete implementations
- The domain model should always be at the very center and depend on nothing external.

### Layers: 

#### Domain Model

This is the core layer and should depend on none of the other pkgs that make up the application. Our domain model is
contained in the mobile package as mobile is our core domain the actual model definitions are located in ``` pkg/mobile/types.go``` 
only core domain types specific to mobile should go here. Other packages can define types as needed for their own needs as long as it 
is not meant to make up part of the core domain model.
The core application interfaces are also located in the core mobile package under ```interfaces.go``` these interfaces, define how the core mobile 
layer expects to be able to interact with external services such as database or web service in order to fulfill its business logic requirements. It is 
the responsibility of these external layers to fulfill this interface.


#### Services (Business logic)

Our services or business logic are contained within the sub packages of the core mobile package. These packages can depend on the core mobile
package and domain model, but not the other way around, they should also not depend on any database specific implementation or know anything about http. 
They should only contain the mobile specific business logic and consume interfaces for interacting with things like a database or external services. They consume the and act on domain models and perform logic using these models.

For example our pkg/mobile/integration has business logic around integrating services with each other. However it knows nothing about kubernetes, instead it relies on implementations of 
its required core interfaces in order to fulfill its business logic.

#### External Layers (tests, web, cli, external services)

The remaining packages make up our concrete implementations of the core interfaces required by the Business logic layer. They can depend on any of the internal layers
they can know everything about a specific database or external service. They should not leak this specific knowledge.
A special note on tests. Even though they are within the same package as the code they are testing, they are an external interface acting on the code they are testing. As such  they often are part of the special
```_test``` package which is treated as a separate package and it only has access to the external API of the package it is testing.
  


```
pkg
├── cli #cli interface (external layer)
│   └── cmd #cli based parsing and handling.
├── clients #client builders for external services (external layer) 
├── data # concrete implementations of the mobile/interfaces.go for data manipulation
├── k8s  # kubernetes client and specific logs (external layer) 
├── mobile
│   ├── interfaces.go # domain logic interfaces (AppRepo, ServiceRepo etc)
│   ├── types.go # domain models (mobile app, mobile services, mobile config)
│   ├── app # core app based business logic
│   ├── integration # core integration business logic
│   └── sdk # core sdk based logic
├── mock # external for mocking in tests
├── openshift # external specific to openshift
└── web # external http handlers and parsers (only resposible for accepting the request parsing and handing off to business logic)
    ├── headers
    └── middleware

```


### How Authentication and Authorization works

**Authentication**

The ui for the MCP is embedded into the OpenShift UI, we are able to reuse the bearer token created when the user first logs into OpenShift.
From the UI, this token is sent with each request to the server. Using this token the access middleware first checks if it needs to check access to the currently
requested URL, if it does, it makes a request to the user read endpoint in OpenShift. If we get a valid response, then we know the user is a 
valid user for this instance of OpenShift and we can proceed.

**Authorization**

Once a user is authenticated, we use their bearer token in order to interact with the kubernetes and openshift API. In this way we offload 
authorisation to the cluster. However is some places it is useful to check access to a particular resource within a given namespace. To do this we 
do a localaccessreview. 

**Service Account**

In the MCP server we also have available a service account that is granted the edit role on the namespace where the server is deployed.
This service account is used when preforming tasks on behalf of a mobile client: for instance, a mobile client will request the sdk service 
configuration from the server. We cannot pass a bearer token from a mobile client. In this case we authenticate the mobile client by its
mobile clientID and its generated API Key we then use the service account token to configure our clients and act against the namespace.


### The Data Layer

MCP does not have a direct dependency on a database of its own. Instead, it levarages data structures provided by the Kubernetes API
in order to store data. So we use objects such as [configmaps](https://kubernetes-v1-4.github.io/docs/user-guide/configmap/) and [secrets](https://kubernetes-v1-4.github.io/docs/user-guide/secrets/) in order to store data that is used to represent our mobile
domain objects. For example a mobile app is backed by a configmap. However the data layer (contained in ```pkg/data```) is the only layer that knows there is a configmap,
every other layer uses an instance of ```mobile.App```. The data layer is responsible for taking a storing domain objects and returning domain objects from whatever backing store is used to store them. 
This separation is enforced via the core interfaces ``` mobile/interfaces.go ``` the data layer fulfills these interfaces, and the services layer depends on the interfaces to perform its logic.


### Adding and testing a new web handler

The web layer is responsible for
accepting requests, parsing the data sent in those requests into the correct domain objects, and handling other request response based actions, this layer should not contain business logic.

Web handlers live under the ```pkg/web``` directory as do the tests that test the handlers. As we base our authentication and authorization from the token passed with
each request, these handlers are also required to configure the kubernetes clients with the correct token to ensure we are acting as the requesting user.
These configured dependencies implement interfaces required by the service business logic. When a request comes in, a cofigured instance is passed through
to the business logic layer. Examples of this are present in most handlers but can be seen for example in the ```web/mobileappHandler.go```. 

Once a new handler is created, you can give it a route via ```web/routes.go``` If it is a completely new route, then a new route method should be added, otherwise the route
can be added to an existing handler.

Most of the handlers here have some level of tests present. These tests start up the server and make real requests against the API. The only layer that is mocked here, is the 
clients that would normally talk to kubernetes, as we don't want to rely on having a running kubernetes just to run our tests.



### Adding and testing new business logic

When adding business logic it should go into a package within the mobile package. Business logic is anything where we are dealing with
a mobile.App or a mobile.Service. Essentially it the mobile specific logic. Test for new business logic should live in the same package as the business logic itself.
The business logic layer shouldn't need to know about things like configmaps, secrets and so on, it should focus on dealing with the domain objects defined in our mobile types.
If it needs to access something from the data layer it should define a method in an existing interface or create a new interface within ```mobile/interfaces.go```. The data layer must then
be changed to add this new logic and method in order to fulfill the required interface. Tests at this layer, should be unit tests 
and execute and test the business logic. Again at this layer we mock the clients talking to the data stores so that we do not need
a running kubernetes in order to execute our tests. 

### Adding a new domain type

Domain types are contained in ```mobile/types.go``` only types that are core to mobile logic should go here. They should represent a 
mobile concern. Examples include ```mobile.App, mobile.Service, mobile.SDKConfig``` these are object that we need and have nothing to do
with kubernetes.

### Adding new data layer functionality

Data layer functionality goes into the data package. There are several good examples of this already present. They illustrate how we
move between configmaps and mobile.App or secrets and services. Tests at this layer also leverage the mock client. Again this is to avoid needing
a running kubernetes to run tests. 