'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:MobileappCtrl
 * @description
 * # MobileappCtrl
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp')
  .controller('MobileappCtrl', ['mcpApi', '$routeParams','$scope','$location',function (mcpApi,$routeParams,$scope,$location) {
    $scope.installType = "";
    var url = new URL(window.location.href)
    $scope.route = url.origin;
    
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

    mcpApi.mobileServices()
    .then(services=>{
      $scope.integrations = services;
    })
    .catch(e=>{
      console.log("error getting services ", e);
    });
    
    
    $scope.installationOpt = function(type){
      $scope.installType = type;
    };
    $scope.sample = "code";
    $scope.codeOpts =  function(type){
      $scope.sample = type;
    };

    $scope.openServiceIntegration = function(name){
      $location.path("/integrations/"+ name);
    };

  }]);
