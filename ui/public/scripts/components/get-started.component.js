'use strict';

/**
 * @ngdoc component
 * @name mcp.component:get-started
 * @description
 * # get-started
 */
angular.module('mobileControlPanelApp').component('getStarted', {
  template: `<div class="blank-slate-pf" id="">
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
                <a ng-href="project/{{ $ctrl.projectName }}/create-mobileservice" class="btn btn-default">Add External Service</a>
              </div>
            </div>`,
  bindings: {
    projectName: '<'
  },
  controller: ['$scope', function($scope) {}]
});
