'use strict';

/**
 * @ngdoc component
 * @name mcp.component:overview
 * @description
 * # overview
 */
angular.module('mobileControlPanelApp').component('overview', {
  template: `<div ng-if="$ctrl.model.objects.length" class="overview container-fluid container-cards-pf">
              <div class="header">
                <h1>{{$ctrl.model.title}}</h1>
                <span class="page-header-link">
                  <a ng-href="http://feedhenry.org/docs/" target="_blank" href="http://feedhenry.org/docs/">
                  Learn More <i class="fa fa-external-link" aria-hidden="true"></i>
                  </a>
                </span>
                <div class="pull-right">
                  <div class="actions" ng-repeat="action in $ctrl.model.actions">
                    <a ng-if="!action.modal" ng-class="['btn', {'btn-default': !action.primary, 'btn-primary': action.primary}]" ng-click="action.action()" ng-if="action.canView()">
                      {{action.label}}
                    </a>
                    <modal modal-class="'control-panel'" ng-if="action.modal" ng-class="{'btn-default': !action.primary, 'btn-primary': action.primary}" modal-open=$ctrl.model.modalOpen launch=action.label modal-title=action.label display-controls=false ng-if="action.canView()">
                      <div class="content" ng-include=action.contentUrl></div>
                    </modal>
                  </div>
                </div>
              </div>

              <div class="row row-cards-pf">
                <div ng-repeat="object in $ctrl.model.objects" class="col-xs-12 col-sm-6 col-md-3">
                  <object-card object=object selected=objectSelected action-selected=actionSelected service-classes=$ctrl.model.serviceClasses></object-card>
                </div>
              </div>
            </div>

            <div ng-if="!$ctrl.model.objects.length" class="blank-slate-pf " id="">
              <div class="blank-slate-pf-icon">
                <span class="pficon pficon pficon-add-circle-o"></span>
              </div>
              <h1>
                Get Started with {{$ctrl.model.title}}
              </h1>
              <p>
                {{$ctrl.model.text}}
              </p>
              <p>
                Learn more about {{$ctrl.model.title}} <a href="http://feedhenry.org/docs/">in the documentation</a>.
              </p>
              <div ng-repeat="action in $ctrl.model.actions | orderBy: 'primary'" class="blank-slate-pf-main-action">
                <a ng-if="!action.modal" ng-class="['btn', {'btn-default': !action.primary, 'btn-primary': action.primary}]" ng-click="action.action()" ng-if="action.canView()">
                    {{action.label}}
                </a>
                <modal modal-class="'control-panel'" ng-if="action.modal" class="btn-default" modal-open=$ctrl.model.modalOpen launch=action.label modal-title=action.label display-controls=false ng-if="action.canView()">
                  <div ng-include=action.contentUrl></div>
                </modal>  
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
