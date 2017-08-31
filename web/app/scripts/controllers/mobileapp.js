'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:MobileappCtrl
 * @description
 * # MobileappCtrl
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp')
  .controller('MobileappCtrl', ['mcpApi', '$routeParams','$scope',function (mcpApi,$routeParams,$scope) {
    console.log("app",$routeParams);
    $scope.route = window.location.host
    mcpApi.mobileApp($routeParams.id)
    .then(app=>{
      console.log("got app",app);
      $scope.app = app;
    })
    .catch(e=>{
      console.error("failed to read app", e);
    });
  }]);
