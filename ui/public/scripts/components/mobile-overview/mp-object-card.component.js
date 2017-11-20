'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-object-card
 * @description
 * # mp-object-card
 */
angular.module('mobileControlPanelApp').component('mpObjectCard', {
  template: `<div class="mp-object-card card-pf card-pf-view card-pf-view-select card-pf-view-multi-select">
              <mp-kebab actions=actions action-selected=actionSelected></mp-kebab>
              <div class="card-pf-body">
                <div class="card-pf-top-element card-icon ng-scope" ng-click="selected($ctrl.object)">
                  <span class="card-pf-icon-circle icon fa {{getIcon($ctrl.object)}}"></span>
                </div>
                <h2 class="card-pf-title text-center">
                  {{getDisplayName($ctrl.object)}}
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
    'ServiceClassService',
    function($scope, ServiceClassService) {
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

      $scope.getDisplayName = function(object) {
        const objectIsService = !!object.integrations;
        if (!objectIsService) {
          return object.name;
        }
        const serviceClass = ServiceClassService.retrieveServiceClass(
          object,
          $scope.$ctrl.serviceClasses
        );
        return ServiceClassService.retrieveDisplayName(
          serviceClass,
          object.name
        );
      };

      $scope.getIcon = function(object) {
        const objectIsService = !!object.integrations;
        if (!objectIsService) {
          return object.metadata.icon;
        }
        const serviceClass = ServiceClassService.retrieveServiceClass(
          object,
          $scope.$ctrl.serviceClasses
        );
        return ServiceClassService.retrieveIcon(serviceClass);
      };
    }
  ]
});
