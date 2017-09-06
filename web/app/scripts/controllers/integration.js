'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:IntegrationsServiceCtrl
 * @description
 * # IntegrationsServiceCtrl
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp')
  .controller('IntegrationCtrl', ['$scope', 'mcpApi', '$routeParams', function ($scope, mcpApi, $routeParams) {
    mcpApi.mobileService($routeParams.service, "true")
      .then(s => {
        $scope.service = s;
      })
      .catch(e => {
        console.error("failed to read service ", e);
      });
      mcpApi.mobileApps()
      .then((apps) => {
        console.log(apps);
        $scope.mobileapps =  {};
        for(var i=0; i < apps.length; i++){
          let app = apps[i];
          $scope.mobileapps[app.clientType] = "true";
        }
        $scope.clients = Object.keys($scope.mobileapps);
        console.log("clients", $scope.clients, $scope.mobileapps);
      })
      .catch(e => {
        console.error(e);
      });  
      $scope.enabled = function(service){
        if(!service){
          return false;
        }
        return service.enabled == true;
      };
      $scope.enableIntegration = function(service){
        console.log("enableing integration",service);
        mcpApi.integrateService(service)
        .then((res)=>{
          console.log("Service integrated");
        })
        .catch(e=>{
          console.log("error integrating service ", e);
        })
        return true;
      };

      
      $scope.installationOpt = function(type){
        $scope.clientType = type;
      };

  }]);
