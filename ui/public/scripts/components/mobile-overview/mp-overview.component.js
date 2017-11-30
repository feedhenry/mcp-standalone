'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-overview
 * @description
 * # mp-overview
 */
angular.module('mobileControlPanelApp').component('mpOverview', {
  template: `<div class="mp-overview">
              <div ng-if="$ctrl.model.objects.length" class="container-fluid container-cards-pf">
                <div class="header object-list row">
                  <div class="column col-xs-12 col-sm-12 col-md-6">
                    <h1>{{$ctrl.model.title}}</h1>
                    <span class="page-header-link">
                      <a ng-href="http://feedhenry.org/docs/" target="_blank" href="http://feedhenry.org/docs/">
                      Learn More <i class="fa fa-external-link" aria-hidden="true"></i>
                      </a>
                    </span>
                  </div>
                  <div class="column col-xs-12 col-sm-12 col-md-6">
                    <div class="pull-right">
                      <div class="actions" ng-repeat="action in $ctrl.model.actions" ng-init="created = action.action" >
                        <a ng-if="!action.modal" ng-class="['btn', {'btn-default': !action.primary, 'btn-primary': action.primary}]" ng-click="action.action()" ng-if="action.canView()">
                          {{action.label}}
                        </a>
                        <mp-modal modal-class="'mp-service-create-modal'" ng-if="action.modal" ng-class="{'btn-default': !action.primary, 'btn-primary': action.primary}" modal-open=$ctrl.model.modalOpen launch=action.label modal-title=action.label display-controls=false ng-if="action.canView()">
                          <div class="content" ng-include=action.contentUrl></div>
                        </mp-modal>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="row row-cards-pf">
                  <div ng-repeat="object in $ctrl.model.objects" class="col-xs-12 col-sm-6 col-md-3">
                    <mp-object-card object=object selected=$ctrl.objectSelected() action-selected=$ctrl.actionSelected() service-classes=$ctrl.model.serviceClasses></mp-object-card>
                  </div>
                </div>
              </div>

              <div ng-if="$ctrl.model.objects.length === 0" class="blank-slate-pf no-objects">
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
                <div ng-repeat="action in $ctrl.model.actions | orderBy: 'primary'" ng-init="created = action.action" class="blank-slate-pf-main-action">
                  <a ng-if="!action.modal" ng-class="['btn', {'btn-default': !action.primary, 'btn-primary': action.primary}]" ng-click="action.action()" ng-if="action.canView()">
                      {{action.label}}
                  </a>
                  <mp-modal modal-class="'mp-service-create-modal'" ng-if="action.modal" class="btn-default" modal-open=$ctrl.model.modalOpen launch=action.label modal-title=action.label display-controls=false ng-if="action.canView()">
                    <div ng-include=action.contentUrl></div>
                  </mp-modal>  
                </div>
              </div>
            </div>`,
  bindings: {
    model: '<',
    objectSelected: '&',
    actionSelected: '&'
  }
});
