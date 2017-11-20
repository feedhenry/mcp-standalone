'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.MobileServiceService
 * @description
 * # MobileServiceService
 * MobileServiceService
 */
angular.module('mobileControlPanelApp').service('MobileServiceService', [
  'ProjectsService',
  'ServiceClassesService',
  'McpService',
  function(ProjectsService, ServiceClassesService, McpService) {
    this.getData = function(projectId, serviceId) {
      return ProjectsService.get(projectId).then(projectInfo => {
        return Promise.all([
          McpService.mobileService(serviceId, 'true'),
          ServiceClassesService.list(projectInfo[1])
        ]);
      });
    };
  }
]);
