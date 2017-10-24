'use strict';

/**
 * @ngdoc component
 * @name mcp.component:object-card
 * @description
 * # ObjectCard
 */
angular.module('mobileControlPanelApp').component('objectCard', {
  template: `<div class="card-pf card-pf-view card-pf-view-select card-pf-view-multi-select">
              <kebab actions=actions action-selected=actionSelected></kebab>
              <div class="card-pf-body">
                <div class="card-pf-top-element card-icon ng-scope" ng-click="selected($ctrl.object)">
                  <span class="card-pf-icon-circle icon fa {{getIcon($ctrl.object)}}"></span>
                </div>
                <h2 class="card-pf-title text-center">
                  {{$ctrl.object.name}}
                </h2>
                <p class="card-pf-info text-center"> {{$ctrl.object.description}}</p>
              </div>
              <div ng-if=objectIsService>
                <p>
                  <span >id: {{$ctrl.object.id}}</span>
                </p>
                <span ng-if="$ctrl.object.labels.external == 'true'" style="color: red;"><em> External Service</em></span>
                <span ng-if="$ctrl.object.labels.external != 'true'" style="color: blue;"><em> Local Service</em></span>
              </div>
            </div>`,
  bindings: {
    object: '<',
    serviceClasses: '<',
    selected: '&',
    actionSelected: '&'
  },
  controller: [
    '$scope',
    function($scope) {
      $scope.objectIsService = !!$scope.$ctrl.object.integrations;
      const actions = ['Delete'];
      $scope.actions = actions.map(action => ({
        label: action,
        value: $scope.$ctrl.object
      }));

      $scope.selected = function(value) {
        $scope.$ctrl.selected()(value);
      };

      $scope.actionSelected = function(value) {
        $scope.$ctrl.actionSelected()(value);
      };

      $scope.getIcon = function(object) {
        const objectIsService = !!object.integrations;
        if (objectIsService) {
          for (var serviceId in $scope.$ctrl.serviceClasses) {
            var serviceClass = $scope.$ctrl.serviceClasses[serviceId];
            var serviceName = serviceClass.spec.externalMetadata.serviceName;
            if (
              serviceName === object.name ||
              (serviceName &&
                serviceName.toLowerCase().indexOf(object.name) >= 0)
            ) {
              if (
                typeof serviceClass.spec.externalMetadata[
                  'console.openshift.io/iconClass'
                ] !== 'undefined'
              ) {
                return formatIconClasses(
                  serviceClass.spec.externalMetadata[
                    'console.openshift.io/iconClass'
                  ]
                );
              }
            }
          }
          return formatIconClasses('fa-clone');
        } else {
          return object.metadata.icon;
        }
      };

      formatIconClasses = function(icon) {
        bits = icon.split('-', 2);
        switch (bits[0]) {
          case 'font':
          case 'icon':
            return 'font-icon ' + icon;
          case 'fa':
            return 'fa ' + icon;
          default:
            return icon;
        }
      };
    }
  ]
});
