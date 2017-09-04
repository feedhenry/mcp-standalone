'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:ServicesCtrl
 * @description
 * # ServicesCtrl
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp')
  .controller('IntegrationsCtrl', ['$scope','mcpApi' ,'$location',function ($scope, mcpApi, $location) {
    mcpApi.mobileServices()
    .then(s=>{
      $scope.services = s;
    })
    .catch(e =>{
      console.error("error getting services ", e);
    });

    $scope.openServiceIntegration = function(name){
      $location.path("/integrations/" + name);
    };
  }]);
