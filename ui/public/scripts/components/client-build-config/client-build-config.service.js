'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.ClientBuildConfigService
 * @description
 * # ClientBuildConfigService
 * ClientBuildConfigService
 */
angular.module('mobileControlPanelApp').service('ClientBuildConfigService', [
  'ProjectsService',
  'McpService',
  'DataService',
  'BuildsService',
  function(ProjectsService, McpService, DataService, BuildsService) {
    this.getData = function(projectId) {
      return ProjectsService.get(projectId).then(function(projectInfo) {
        const projectContext = projectInfo[1];
        return Promise.all([
          Promise.resolve(projectContext),
          DataService.list('buildconfigs', projectContext)
        ]);
      });
    };

    this.createBuildConfig = function(appConfig, projectContext) {
      return McpService.createBuildConfig(appConfig).then(response => {
        return DataService.get('buildconfigs', appConfig.name, projectContext);
      });
    };

    this.updateBuildConfig = function(appConfig, projectContext) {
      return DataService.update(
        'buildconfigs',
        appConfig.metadata.name,
        appConfig,
        projectContext
      ).then(() => {
        return DataService.get(
          'buildconfigs',
          appConfig.metadata.name,
          projectContext
        );
      });
    };

    this.startBuild = function(buildConfig) {
      return BuildsService.startBuild(buildConfig);
    };
  }
]);
