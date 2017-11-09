# Mobile Control Panel UI

## Overview

The Mobile Control Panel (MCP) UI is built using:
* AngularJS https://angularjs.org/
* Patternfly http://www.patternfly.org/

It is developed as a set of OpenShift UI Extensions.
All Javascript is delivered as a single file (mcp.js).
All CSS is delivered as a single file (mcp.css).

The extension config looks something like this in the OpenShift master-config.yaml

```yaml
assetConfig:
  extensionDevelopment: true
  extensionProperties: null
  extensionScripts:
    - /var/lib/origin/openshift.local.config/master/servicecatalog-extension.js
    - /var/lib/origin/openshift.local.config/public/mcp.js
  extensionStylesheets:
    - /var/lib/origin/openshift.local.config/public/mcp.css
  extensions:
    - name: mcp
      sourceDirectory: /var/lib/origin/openshift.local.config/public
```

When doing local development against an `oc cluster`, the `ui` folder would be mounted into the docker `origin` container as the config dir. This allows the path `/var/lib/origin/openshift.local.config/public` to point to the `./ui/public` folder.

## Style guide

MCP UI is an angular application and we are using angularjs version 1.5.
Angular promotes a component based application architecture to take advantage of component benefits like:

  - Reusability - Components are atomic units, and building with components allows for their reuse in future development cycles.
  - Consistency - Implementing reusable components helps keep design consistent and can provide clarity in organising code. Also keeps consistency of use for the user.
  - Maintainability - A set of well organised components can be quick to update, and you can be more confident about which areas will and won't be affected.
  - Scalability - Having a library of components to implement can make for speedy development.

Angular has the [component directive](https://docs.angularjs.org/guide/component) which should be the main building block for UI features.
[This style guide](https://github.com/toddmotto/angularjs-styleguide) will form the basis of the MCP UI style guide so have a read before starting to code.

NOTE: We are currently moving to this style code so the code may not match the guide in places. Also our app uses standard angular router instread of ui-router and babel is not configured so import/export is not available.
