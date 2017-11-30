'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-get-started
 * @description
 * # mp-get-started
 */
angular.module('mobileControlPanelApp').component('mpGetStarted', {
  template: `<div class="blank-slate-pf" ng-if="$ctrl.hasData && $ctrl.overviewsEmpty">
              <div class="blank-slate-pf-icon">
                <span class="pficon pficon pficon-add-circle-o"></span>
              </div>
              <h1>
                Get Started with Mobile Apps & Services
              </h1>
              <p>
                You can create a Mobile App to enable Mobile Integrations with Mobile Enabled Services.
              </p>
              <p>
                You can provision or link a Mobile Enabled Service to enable a Mobile App Integration.
              </p>
              <p>
                Learn more about Mobile Apps & Services <a href="http://feedhenry.org/docs/">in the documentation</a>.
              </p>
              <div class="blank-slate-pf-main-action">
                <a ng-href="project/{{ $ctrl.projectName }}/create-mobileapp" class="btn btn-primary btn-lg">Create Mobile App</a>
              </div>
              <div class="blank-slate-pf-secondary-action">
                <a ng-href="/" class="btn btn-default">Provision Catalog Service</a>
                <mp-modal ng-init="created = $ctrl.created" class="btn-default" modal-open=$ctrl.options.modalOpen modal-class="'mp-service-create-modal'" launch="'Add External Service'" modal-title="'Add External Service'"" display-controls=false>
                  <div ng-include="'extensions/mcp/templates/create-service.template.html'"></div>
                </mp-modal>
              </div>
            </div>`,
  bindings: {
    overviews: '<',
    projectName: '<',
    serviceCreated: '&?',
    options: '<'
  },
  controller: [
    function() {
      this.overviewEmpty = false;

      this.$onChanges = function(changes) {
        if (!changes.overviews) {
          return;
        }

        const currentValue = changes.overviews.currentValue;
        const keys = Object.keys(currentValue);
        this.hasData = keys.every(key => currentValue[key].objects);
        if (this.hasData) {
          this.overviewsEmpty = keys.every(
            key => !currentValue[key].objects.length
          );
        }
      };

      this.created = function(err, service) {
        this.serviceCreated()(err, service);
      };
    }
  ]
});
