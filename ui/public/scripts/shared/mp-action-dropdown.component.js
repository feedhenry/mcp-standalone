'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-action-dropdown
 * @description
 * # mp-action-dropdown
 */
angular.module('mobileControlPanelApp').component('mpActionDropdown', {
  template: `<div class="dropdown">
              <button type="button" class="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-expanded="true">
              Actions
                <span class="caret"></span>
              </button>
              <ul class="dropdown-menu" aria-labelledby="dropdownMenu1">
                <li ng-repeat="action in $ctrl.actions"><a href="#" ng-click="select(action.value, $index, action)">{{action.label}}</a></li>
              </ul>
            </div>`,
  bindings: {
    actions: '<',
    selected: '&'
  },
  controller: [
    '$scope',
    function($scope) {
      $scope.select = function(value, index, action) {
        $scope.$ctrl.selected()(value, index, action);
      };
    }
  ]
});
