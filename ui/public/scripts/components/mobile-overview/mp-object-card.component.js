'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-object-card
 * @description
 * # mp-object-card
 */
angular.module('mobileControlPanelApp').component('mpObjectCard', {
  template: `<div class="mp-object-card card-pf card-pf-view card-pf-view-select card-pf-view-multi-select">
              <mp-kebab actions=$ctrl.actions action-selected=$ctrl.actionSelected()></mp-kebab>
              <div class="card-pf-body">
                <div class="card-pf-top-element card-icon ng-scope" ng-click="$ctrl.selected()($ctrl.object)">
                  <mp-object-icon object=$ctrl.object service-classes=$ctrl.serviceClasses></mp-object-icon>
                </div>
                <h2 class="card-pf-title text-center">
                  {{ $ctrl.object | objectName:$ctrl.serviceClasses }}
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
    function() {
      this.objectIsService = !!this.object.integrations;
      const actions = ['Delete'];
      this.actions = actions.map(action => ({
        label: action,
        value: this.object
      }));
    }
  ]
});
