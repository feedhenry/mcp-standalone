'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-mobile-service-integration
 * @description
 * # mp-mobile-service-integration
 */
angular
  .module('mobileControlPanelApp')
  .component('mpMobileServiceIntegration', {
    template: `<mp-integrations-list integrations=$ctrl.integrations apps=$ctrl.apps integration-toggled=$ctrl.integrationToggled project-name=$ctrl.projectName></mp-integrations-list>`,
    controller: [
      'ServiceIntegrationsService',
      '$routeParams',
      function(ServiceIntegrationsService, $routeParams) {
        this.$onInit = function() {
          this.service = {};
          this.projectName = $routeParams.project;

          ServiceIntegrationsService.getIntegrationInfo($routeParams.service)
            .then(integrationInfo => {
              const [service = {}, apps = []] = integrationInfo;

              this.service = service;
              this.apps = apps;
              this.integrations = Object.keys(service.integrations).map(key => {
                return Object.assign(service.integrations[key], {
                  target: service
                });
              });
            })
            .catch(err => {
              console.error('Error loading integration info', err);
            });
        };

        this.integrationToggled = function(integration, enabled) {
          let promise;
          if (enabled) {
            promise = ServiceIntegrationsService.enableIntegration(integration);
          } else {
            promise = ServiceIntegrationsService.disableIntegration(
              integration
            );
          }

          promise
            .then(() => {
              integration.enabled = enabled;
            })
            .catch(err => {
              console.log('Error with service integration ', err);
            });
        };
      }
    ]
  });
