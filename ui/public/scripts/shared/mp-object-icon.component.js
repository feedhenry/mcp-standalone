'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-object-icon
 * @description
 * # mp-object-icon
 */
angular.module('mobileControlPanelApp').component('mpObjectIcon', {
  template: `<span class="card-pf-icon-circle icon fa {{$ctrl.getIconClass()}}"></span>`,
  bindings: {
    object: '<',
    serviceClasses: '<'
  },
  controller: [
    'ServiceClassesService',
    function(ServiceClassesService) {
      this.getIconClass = function() {
        const objectIsService = !!this.object.integrations;
        if (!objectIsService) {
          return this.object.metadata.icon;
        }

        if (!this.serviceClasses) {
          return '';
        }

        const serviceClass = ServiceClassesService.getServiceClass(
          this.object,
          this.serviceClasses
        );
        const iconClass =
          serviceClass &&
          serviceClass.spec.externalMetadata['console.openshift.io/iconClass'];
        if (iconClass) {
          return this.formatIconClass(iconClass);
        }

        return this.formatIconClass('fa-clone');
      };

      this.formatIconClass = function(icon) {
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
