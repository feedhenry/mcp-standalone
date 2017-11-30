'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-mobile-overview
 * @description
 * # mp-mobile-overview
 */
angular.module('mobileControlPanelApp').component('mpMobileOverview', {
  template: `<mp-setup has-mcp-server=!$ctrl.mcpError></mp-setup>
            <mp-get-started overviews=$ctrl.overviews project-name=$ctrl.projectContext.projectName service-created=$ctrl.getStartedServiceCreated options=$ctrl.getStartedOptions></mp-get-started>
            <mp-overview-list overviews=$ctrl.overviews object-selected=$ctrl.objectSelected action-selected=$ctrl.actionSelected></mp-overview-list>`,
  bindings: {},
  controller: [
    '$routeParams',
    '$location',
    'MobileOverviewService',
    function($routeParams, $location, MobileOverviewService) {
      const ctrl = this;

      Object.assign(ctrl, {
        mcpError: false,
        getStartedOptions: {
          modalOpen: false
        },
        overviews: {
          apps: {},
          services: {}
        }
      });

      ctrl.$onInit = function() {
        MobileOverviewService.getOverview($routeParams.project)
          .then(overview => {
            const [
              project = {},
              projectContext = {},
              apps = [],
              services = [],
              serviceClasses
            ] = overview;

            ctrl.project = project;
            ctrl.projectContext = projectContext;

            ctrl.overviews.apps = {
              type: 'app',
              title: 'Mobile Apps',
              text:
                'You can create a Mobile App to enable Mobile Integrations with Mobile Enabled Services',
              actions: [
                {
                  label: 'Create Mobile App',
                  primary: true,
                  action: $location.path.bind(
                    $location,
                    `project/${projectContext.projectName}/create-mobileapp`
                  ),
                  canView: () => true
                }
              ]
            };
            ctrl.overviews.services = {
              type: 'service',
              title: 'Mobile Enabled Services',
              modalOpen: false,
              text:
                'You can provision or link a Mobile Enabled Service to enable a Mobile App Integration.',
              actions: [
                {
                  label: 'Add External Service',
                  modal: true,
                  contentUrl:
                    'extensions/mcp/templates/create-service.template.html',
                  action: function(err) {
                    if (err) {
                      return;
                    }

                    MobileOverviewService.getServices().then(services => {
                      ctrl.overviews.services.objects = services;
                      ctrl.overviews.services.modalOpen = false;
                      ctrl.overviews = Object.assign({}, ctrl.overviews);
                    });
                  },
                  canView: () =>
                    MobileOverviewService.canViewService(projectContext)
                },
                {
                  label: 'Provision Catalog Service',
                  primary: true,
                  action: $location.path.bind($location, `/`),
                  canView: () =>
                    MobileOverviewService.canViewService(projectContext)
                }
              ]
            };

            ctrl.overviews.apps.objects = apps;
            ctrl.overviews.services.objects = services;
            ctrl.overviews.services.serviceClasses = serviceClasses['_data'];

            ctrl.overviews = Object.assign({}, ctrl.overviews);
          })
          .catch(err => {
            console.error('Error getting overview ', err);
            ctrl.mcpError = true;
          });
      };

      ctrl.actionSelected = function(object) {
        const objectIsService = !!object.integrations;
        const actionFn = objectIsService ? 'deleteService' : 'deleteApp';
        const objectType = objectIsService ? 'services' : 'apps';

        MobileOverviewService[actionFn](object).then(objects => {
          ctrl.overviews[objectType].objects = objects;
          ctrl.overviews = Object.assign({}, ctrl.overviews);
        });
      };

      ctrl.getStartedServiceCreated = function(err) {
        if (err) {
          return;
        }

        return MobileOverviewService.getServices().then(services => {
          ctrl.getStartedOptions.modalOpen = false;
          ctrl.overviews.services.objects = services;
          ctrl.overviews = Object.assign({}, ctrl.overviews);
        });
      };

      ctrl.objectSelected = function(object) {
        const objectIsService = !!object.integrations;
        const objectRoute = objectIsService ? 'mobileservices' : 'mobileapps';
        $location.path(
          `project/${$routeParams.project}/browse/${objectRoute}/${object.id}`
        );
      };
    }
  ]
});
