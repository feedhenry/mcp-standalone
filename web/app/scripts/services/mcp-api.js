'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.mcpApi
 * @description
 * # mcpApi
 * Service in the mobileControlPanelApp.
 */
angular.module('mobileControlPanelApp')
  .service('mcpApi', ['$http','AuthService',function ($http,AuthService) {
    
    let mobileAppsURL = "/mobileapp";
    let mobileServicesURL = "/mobileservice";
    // AngularJS will instantiate a singleton by calling "new" on this function
    let requestConfig = {"headers":{}};
    AuthService.addAuthToRequest(requestConfig);
    
    return{
      "mobileApps" : function(){
        return $http.get(mobileAppsURL,requestConfig)
        .then((res)=>{
          return res.data;
        })
        .catch(err=>{
          return err;
        });
      },
      "mobileApp": function(id){
        return $http.get(mobileAppsURL + "/"+id,requestConfig)
        .then((res)=>{
          return res.data;
        })
        .catch(err=>{
          return err;
        });
      },
      "createMobileApp":function(mobileApp){
        return $http.post(mobileAppsURL,mobileApp,requestConfig)
        .then((res)=>{
          return res.data;
        })
        .catch(err=>{
          return err;
        });
      },
      "mobileServices": function(){
        return $http.get(mobileServicesURL,requestConfig)
        .then((res)=>{
          return res.data;
        })
        .catch(err=>{
          return err;
        });
      },
      "mobileService": function(name, withIntegrations){
        let url = mobileServicesURL + "/"+name;
        if(withIntegrations){
          console.log("withIntegrations");
          url += "?withIntegrations=true";
        }
        console.log("calling ", url)
        return $http.get(url,requestConfig)
        .then((res)=>{
          return res.data;
        })
        .catch(err=>{
          return err;
        });
      },
      "integrateService": function(params){
        let url = mobileServicesURL + "/configure";
        return $http.post(url,params,requestConfig)
        .then((res)=>{
          return res.data;
        })
        .catch(err=>{
          return err;
        });
      }
    };
  }]);
