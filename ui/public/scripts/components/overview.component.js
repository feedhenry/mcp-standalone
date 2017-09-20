'use strict';

/**
 * @ngdoc component
 * @name mcp.component:overview
 * @description
 * # overview
 */
angular.module('mobileControlPanelApp').component('overview', {
  template: `<div class="overview container-fluid container-cards-pf">
              <div class="header">
                <h1>{{$ctrl.model.title}}</h1>
                <span class="page-header-link">
                  <a ng-href="http://feedhenry.org/docs/" target="_blank" href="http://feedhenry.org/docs/">
                  Learn More <i class="fa fa-external-link" aria-hidden="true"></i>
                  </a>
                </span>
                <div class="pull-right">
                  <a ng-repeat="action in $ctrl.model.actions" ng-class="['btn', {'btn-default': !action.primary, 'btn-primary': action.primary}]" ng-click="action.action()" ng-if="action.canView()">
                    {{action.label}}
                  </a>
                </div>
              </div>

              <div class="row row-cards-pf">
                <div ng-repeat="object in $ctrl.model.objects" class="col-xs-12 col-sm-6 col-md-3">
                  <object-card object=object selected=objectSelected action-selected=actionSelected service-classes=$ctrl.model.serviceClasses></object-card>
                </div>
              </div>
            </div>`,
  bindings: {
    model: '<',
    objectSelected: '&',
    actionSelected: '&'
  },
  controller: [
    '$scope',
    function($scope) {
      $scope.objectSelected = function(value) {
        $scope.$ctrl.objectSelected()(value);
      };

      $scope.actionSelected = function(value) {
        $scope.$ctrl.actionSelected()(value);
      };
    }
  ]
});
