'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.ServiceDashboardService
 * @description
 * # ServiceDashboardService
 * ServiceDashboardService
 */
angular.module('mobileControlPanelApp').service('ServiceDashboardService', [
  'DataService',
  'ProjectsService',
  'McpService',
  function(DataService, ProjectsService, McpService) {
    this.getChartData = function(serviceId) {
      return McpService.mobileServiceMetrics(serviceId).then(chartData => {
        return chartData;
      });
    };

    this.getServiceClasses = function(projectId) {
      return ProjectsService.get(projectId)
        .then(projectInfo => {
          const [project = {}, projectContext = {}] = projectInfo;

          return DataService.list(
            {
              group: 'servicecatalog.k8s.io',
              resource: 'clusterserviceclasses'
            },
            projectContext
          );
        })
        .then(serviceClasses => {
          return serviceClasses['_data'];
        });
    };

    this.getService = function(serviceId) {
      return McpService.mobileService(serviceId, 'true');
    };
  }
]);
