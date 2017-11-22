'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-service-dashboard
 * @description
 * # mp-service-dashboard
 */
angular.module('mobileControlPanelApp').component('mpServiceDashboard', {
  template: `<div class="mp-service-dashboard">
              <mp-service-description service=$ctrl.service service-classes=$ctrl.serviceClasses></mp-service-description>
              <h2>Stats</h2>
              <mp-chart-data chart-data=$ctrl.chartData service=$ctrl.service></mp-chart-data>
            </div>`,
  controller: [
    '$routeParams',
    'ServiceDashboardService',
    function($routeParams, ServiceDashboardService, $timeout) {
      this.$onInit = function() {
        this.chartData = [];

        ServiceDashboardService.getChartData($routeParams.service)
          .then(chartData => {
            this.chartData = chartData;
          })
          .catch(err => console.error('Error retrieving chart data', err));

        Promise.all([
          ServiceDashboardService.getServiceClasses($routeParams.project),
          ServiceDashboardService.getService($routeParams.service)
        ])
          .then(dashBoardInfo => {
            const [serviceClasses = {}, service = {}] = dashBoardInfo;
            this.serviceClasses = serviceClasses;
            this.service = service;
          })
          .catch(err => console.error('Error initialising Dashboard', err));
      };
    }
  ]
});
