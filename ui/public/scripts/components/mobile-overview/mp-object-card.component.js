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
                  <span class="card-pf-icon-circle icon fa {{$ctrl.getIcon($ctrl.object)}}"></span>
                </div>
                <h2 class="card-pf-title text-center">
                  {{$ctrl.object.displayName || $ctrl.object.name}}
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

      this.getIcon = function(object) {
        const objectIsService = !!object.integrations;
        if (objectIsService) {
          for (var serviceId in this.serviceClasses) {
            var serviceClass = this.serviceClasses[serviceId];
            var serviceName = serviceClass.spec.externalMetadata.serviceName;
            if (
              serviceName === object.name ||
              (serviceName &&
                serviceName.toLowerCase().indexOf(object.name) >= 0)
            ) {
              if (
                typeof serviceClass.spec.externalMetadata[
                  'console.openshift.io/iconClass'
                ] !== 'undefined'
              ) {
                return formatIconClasses(
                  serviceClass.spec.externalMetadata[
                    'console.openshift.io/iconClass'
                  ]
                );
              }
            }
          }
          return formatIconClasses('fa-clone');
        } else {
          return object.metadata.icon;
        }
      };

      formatIconClasses = icon => {
        bits = icon.split('-', 2);
        switch (bits[0]) {
          case 'font':
          case 'icon':
            return 'font-icon ' + icon;
          case 'fa':
            return 'fa ' + icon;
          default:
            return icon;
        }
      };
    }
  ]
});
