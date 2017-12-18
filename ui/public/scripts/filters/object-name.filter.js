'use strict';

/**
 * @ngdoc filter
 * @name mobileControlPanelApp.objectName
 * @description
 * # Object Name Angular Filter
 */

angular.module('mobileControlPanelApp').filter('objectName', [
  'ServiceClassesService',
  function(ServiceClassesService) {
    return function(object, serviceClasses) {
      if (!object) {
        return '';
      }

      const objectIsService = !!object.integrations;
      if (!objectIsService) {
        return object.name || '';
      }

      const serviceClass = ServiceClassesService.getServiceClass(
        object,
        serviceClasses
      );

      if (!serviceClass) {
        return object.name || '';
      }

      return (
        serviceClass.spec.externalMetadata.displayName || object.name || ''
      );
    };
  }
]);
