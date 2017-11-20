'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.MobileOverviewService
 * @description
 * # MobileOverviewService
 * MobileOverviewService
 */
angular.module('mobileControlPanelApp').service('MobileOverviewService', [
  'ServiceClassesService',
  'ProjectsService',
  'AuthorizationService',
  'McpService',
  function(
    ServiceClassesService,
    ProjectsService,
    AuthorizationService,
    McpService
  ) {
    this.getOverview = function(projectId) {
      return ProjectsService.get(projectId).then(projectInfo => {
        return Promise.all([
          Promise.resolve(projectInfo[0]),
          Promise.resolve(projectInfo[1]),
          McpService.mobileApps(),
          McpService.mobileServices(),
          ServiceClassesService.list(projectInfo[1])
        ]);
      });
    };

    this.getServices = function() {
      return McpService.mobileServices();
    };

    this.getApps = function() {
      return McpService.mobileApps();
    };

    this.deleteService = function(object) {
      return McpService.deleteService(object).then(result =>
        McpService.mobileServices()
      );
    };

    this.deleteApp = function(object) {
      return McpService.deleteApp(object).then(result =>
        McpService.mobileApps()
      );
    };

    this.canViewService = function(projectContext) {
      return AuthorizationService.canI(
        'services',
        'create',
        projectContext.projectName
      );
    };
  }
]);
