'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-overview-list
 * @description
 * # mp-overview-list
 */
angular.module('mobileControlPanelApp').component('mpOverviewList', {
  template: `<div class="cards-pf" ng-if="$ctrl.overviews.apps.objects.length > 0 || $ctrl.overviews.services.objects.length > 0">
              <mp-overview ng-repeat="overview in $ctrl.overviews" model=overview object-selected=$ctrl.objectSelected() action-selected=$ctrl.actionSelected()></mp-overview>
            </div>`,
  bindings: {
    overviews: '<',
    objectSelected: '&',
    actionSelected: '&'
  }
});
