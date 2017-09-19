'use strict';
/**
 * @ngdoc component
 * @name mcp.component:mobile-overview
 * @description
 * # mobile-overview
 */
angular.module('mobileControlPanelApp').component('mobileOverview', {
  template: `<div class="cards-pf">
              <overview ng-repeat="overview in $ctrl.overviews" model=overview object-selected=objectSelected action-selected=actionSelected></overview>
            </div>`,
  bindings: {
    overviews: '<',
    objectSelected: '&',
    actionSelected: '&'
  },
  controller: [
    '$scope',
    function($scope) {
      $scope.actionSelected = function(object) {
        $scope.$ctrl.actionSelected()(object);
      };

      $scope.objectSelected = function(object) {
        $scope.$ctrl.objectSelected()(object);
      };
    }
  ]
});
