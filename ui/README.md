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