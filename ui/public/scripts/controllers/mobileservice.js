'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:MobileServiceController
 * @description
 * # MobileServiceController
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp').controller('MobileServiceController', [
  '$scope',
  '$routeParams',
  function($scope, $routeParams) {
    $scope.alerts = {};
    $scope.breadcrumbs = [
      {
        title: 'Mobile Service',
        link: 'project/' + $routeParams.project + '/browse/mobileoverview'
      },
      {
        title: $routeParams.service
      }
    ];
  }
]);
