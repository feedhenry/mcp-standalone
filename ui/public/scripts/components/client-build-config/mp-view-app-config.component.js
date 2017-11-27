'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-view-app-config
 * @description
 * # mp-view-app-config
 */
angular.module('mobileControlPanelApp').component('mpViewAppConfig', {
  template: `<div>
              <div class="row mp-app-config">
                <div class="col-lg-6">
                  <h3>Details</h3>
                  <dl class="dl-horizontal left">
                    <div>
                      <dt>Repo URL</dt>
                      <dd>{{$ctrl.config.spec.source.git.uri}}</dd>
                    </div>
                    <div>
                      <dt>Jenkins Job Name</dt>
                      <dd>{{$ctrl.config.metadata.name}}</dd>
                    </div>
                    <div>
                      <dt>Branch</dt>
                      <dd>{{$ctrl.config.spec.source.git.ref}}</dd>
                    </div>
                    <div>
                      <dt>Jenkinsfile Path</dt>
                      <dd>{{$ctrl.config.spec.strategy.jenkinsPipelineStrategy.jenkinsfilePath}}</dd>
                    </div>
                  </dl>
                </div>
              </div>
            </div>`,
  bindings: {
    config: '<'
  }
});
