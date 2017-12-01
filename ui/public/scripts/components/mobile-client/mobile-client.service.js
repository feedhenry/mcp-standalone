'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.MobileClientService
 * @description
 * # MobileClientService
 * MobileClientService
 */
angular.module('mobileControlPanelApp').service('MobileClientService', [
  'ProjectsService',
  'DataService',
  'BuildsService',
  'McpService',
  function(ProjectsService, DataService, BuildsService, McpService) {
    this.getData = function(projectId, appId) {
      return ProjectsService.get(projectId).then(projectInfo => {
        const [project = {}, projectContext = {}] = projectInfo;

        return Promise.all([
          Promise.resolve(project),
          Promise.resolve(projectContext),
          McpService.mobileApp(appId),
          DataService.list('buildconfigs', projectContext),
          DataService.list('builds', projectContext),
          DataService.list('secrets', projectContext)
        ]);
      });
    };

    this.startBuild = function(buildConfig) {
      return BuildsService.startBuild(buildConfig);
    };

    this.watch = function(name, projectContext, cb) {
      return DataService.watch(name, projectContext, cb);
    };

    this.unwatchAll = function(watches) {
      return DataService.unwatchAll(watches);
    };
  }
]);
