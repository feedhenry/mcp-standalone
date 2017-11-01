'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-kebab
 * @description
 * # mp-kebab
 */
angular.module('mobileControlPanelApp').component('mpKebab', {
  template: `<div class="dropdown  dropdown-kebab-pf pull-right">
                <button class="btn btn-link dropdown-toggle" type="button" id="dropdownKebab" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
                  <span class="fa fa-ellipsis-v"></span>
                </button>
                <ul class="dropdown-menu " aria-labelledby="dropdownKebab">
                  <li ng-repeat="action in $ctrl.actions">
                    <a ng-click="actionSelected(action.value, $index, action)" href="">{{action.label || action.value}}</a>
                  <li/>
                </ul>
              </div>`,
  bindings: {
    actions: '<',
    actionSelected: '&'
  },
  controller: [
    '$scope',
    function($scope) {
      $scope.actionSelected = function(value, index, option) {
        $scope.$ctrl.actionSelected()(value, index, option);
      };
    }
  ]
});
