'use strict';

(function() {
  // This is the default configuration for the dev mode of the web console.
  // A generated version of this config is created at run-time when running
  // the web console from the openshift binary.
  //
  // To change configuration for local development, copy this file to
  // assets/app/config.local.js and edit the copy.
  var masterPublicHostname = '127.0.0.1:8443';

  window.MCP_CONFIG = {
    api: {
      host: 'https://127.0.0.1:3001'
    }
  };

  window.OPENSHIFT_CONFIG = {
    apis: {
      hostPort: masterPublicHostname,
      prefix: "/apis"
    },
    api: {
      openshift: {
        hostPort: masterPublicHostname,
        prefix: "/oapi"
      },
      k8s: {
        hostPort: masterPublicHostname,
        prefix: "/api"
      }
    },
    auth: {
      oauth_authorize_uri: 'https://' + masterPublicHostname + "/oauth/authorize",
      oauth_redirect_base: window.MCP_CONFIG.api.host + "/console",
      oauth_client_id: "system:serviceaccount:myproject:mobile-server",
      oauth_token_uri: window.MCP_CONFIG.api.host + "/oauth/token",
      logout_uri: ""
    },
    loggingURL: "",
    metricsURL: ""
  };

  window.OPENSHIFT_VERSION = {
    openshift: "dev-mode",
    kubernetes: "dev-mode"
  };

})();
