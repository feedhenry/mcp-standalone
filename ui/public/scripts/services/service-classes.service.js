'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.ServiceClassesService
 * @description
 * # ServiceClasses Angular Service
 */
angular.module('mobileControlPanelApp').service('ServiceClassesService', [
  'DataService',
  function(DataService) {
    return {
      list: function(projectContext) {
        return DataService.list(
          {
            group: 'servicecatalog.k8s.io',
            resource: 'clusterserviceclasses'
          },
          projectContext
        );
      },

      getServiceClass: function(object, serviceClasses) {
        for (var serviceId in serviceClasses) {
          const serviceClass = serviceClasses[serviceId];
          const serviceName = serviceClass.spec.externalMetadata.serviceName;
          if (
            serviceName === object.name ||
            (serviceName && serviceName.toLowerCase().indexOf(object.name) >= 0)
          ) {
            return serviceClass;
          }
        }
      }
    };
  }
]);
