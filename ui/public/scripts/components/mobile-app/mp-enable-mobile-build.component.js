'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-enable-mobile-build
 * @description
 * # mp-enable-mobile-build
 */
angular.module('mobileControlPanelApp').component('mpEnableMobileBuild', {
  template: `<div class="blank-slate-pf" id="">
              <div class="blank-slate-pf-icon">
                <span class="pficon pficon pficon-add-circle-o"></span>
              </div>
              <h1>
                Enable Mobile CI/CD Service
              </h1>
              <p>
                To enable mobile application builds, please provision the Mobile CI/CD service.
              </p>
              <p>
                This can be provisioned via the Service Catalog.
              </p>
              <p>
                Learn more about this <a href="http://feedhenry.org/docs/">in the documentation</a>.
              </p>
              <div class="blank-slate-pf-main-action">
                <a ng-href="/" class="btn btn-primary btn-lg">Provision Mobile CI/CD Service</a>
              </div>
            </div>`
});
