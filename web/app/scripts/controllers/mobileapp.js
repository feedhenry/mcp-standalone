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
    $scope.installType = "";
    $scope.route = window.location.host;
    
    mcpApi.mobileApp($routeParams.id)
    .then(app=>{
      $scope.app = app;
      switch(app.clientType){
        case "cordova":
        $scope.installType = "npm";
        break;
        case "android":
        $scope.installType = "maven";
        break;
        case "iOS":
        $scope.installType = "cocoapods";
        break;
      }
    })
    .catch(e=>{
      console.error("failed to read app", e);
    });
    
    
    $scope.installationOpt = function(type){
      $scope.installType = type;
    };
    $scope.sample = "code";
    $scope.codeOpts =  function(type){
      $scope.sample = type;
    };
  }]);
