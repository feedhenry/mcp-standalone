'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-loading
 * @description
 * # mp-loading
 */
angular.module('mobileControlPanelApp').component('mpLoading', {
  template: `<div ng-if="$ctrl.loading">
                Loading...
              </div>
              <div ng-if="!$ctrl.loading">
                <ng-transclude></ng-transclude>
              </div>`,
  transclude: true,
  bindings: {
    loading: '<'
  }
});
