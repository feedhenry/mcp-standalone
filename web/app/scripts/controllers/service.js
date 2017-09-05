'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:ServiceCtrl
 * @description
 * # ServiceCtrl
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp')
  .controller('ServiceCtrl', ['$scope', '$routeParams', 'mcpApi', function ($scope, $routeParams, mcpApi) {
    mcpApi.mobileService($routeParams.id)
    .then(service=>{
      $scope.service = service;
    })
    .catch(e=>{
      console.error("failed to read service", e);
    });
  }]);
