'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.ServiceClassService
 * @description
 * # ServiceClass Angular Service
 */
angular.module('mobileControlPanelApp').service('ServiceClassService', [
  function() {
    return {
      retrieveServiceClass: function(service, serviceClasses) {
        for (var serviceId in serviceClasses) {
          const serviceClass = serviceClasses[serviceId];
          const serviceName = serviceClass.spec.externalMetadata.serviceName;
          if (
            serviceName === service.name ||
            (serviceName &&
              serviceName.toLowerCase().indexOf(service.name) >= 0)
          ) {
            return serviceClass;
          }
        }
      },
      retrieveIcon: function(serviceClass) {
        const iconClass =
          serviceClass.spec.externalMetadata['console.openshift.io/iconClass'];
        if (!serviceClass || typeof iconClass === 'undefined') {
          return this.formatIconClass('fa-clone');
        }
        return this.formatIconClass(iconClass);
      },
      retrieveDisplayName: function(serviceClass, defaultName) {
        if (!serviceClass) {
          return defaultName || 'No Name';
        }
        return serviceClass.spec.externalMetadata.displayName || defaultName;
      },
      formatIconClass: function(icon) {
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
      }
    };
  }
]);
