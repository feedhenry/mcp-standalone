'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.ClientBuildHistoryService
 * @description
 * # ClientBuildHistoryService
 * ClientBuildHistoryService
 */
angular.module('mobileControlPanelApp').service('ClientBuildHistoryService', [
  'ProjectsService',
  'DataService',
  function(ProjectsService, DataService) {
    this.getData = function(projectId) {
      return ProjectsService.get(projectId).then(projectInfo => {
        const projectContext = projectInfo[1];

        return Promise.all([
          Promise.resolve(projectContext),
          DataService.list('buildconfigs', projectContext),
          DataService.list('builds', projectContext)
        ]);
      });
    };

    this.watch = function(name, projectContext, cb) {
      return DataService.watch(name, projectContext, cb);
    };

    this.unwatchAll = function(watches) {
      DataService.unwatchAll(watches);
    };
  }
]);
