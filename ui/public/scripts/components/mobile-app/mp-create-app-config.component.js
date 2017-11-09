'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-create-app-config
 * @description
 * # mp-create-app-config
 */
angular.module('mobileControlPanelApp').component('mpCreateAppConfig', {
  template: `<div>
              <div class="row create mp-app-config">
                <div class="col-lg-6">
                  <h3>Create</h3>
                  <form novalidate class="form-horizontal" name="appBuildConfig">
                    <dl class="dl-horizontal left">
                      <div ng-class="{'has-error': appBuildConfig.repoUri.$touched && appBuildConfig.repoUri.$error.required}">
                        <dt>Repo URL</dt>
                        <dd>
                          <input ng-model="config.gitRepo.uri" name="repoUri" type="text" id="repo-uri" class="form-control" required>
                            <span ng-if="appBuildConfig.repoUri.$touched && appBuildConfig.repoUri.$error.required" class="help-block error">
                              The App Repo URI is required.
                            </span>
                        </dd>
                      </div>
                      <div class="name-field">
                        <dt>Jenkins Job Name</dt>
                        <dd>
                          <div ng-class="{'has-error': appBuildConfig.buildname.$touched && (appBuildConfig.buildname.$error.pattern || appBuildConfig.buildname.$error.required)}">
                            <input placeholder="A unique name for the build config." ng-model="config.name" name="buildname" type="text" id="build-name" class="form-control" required pattern="[a-z0-9]([-a-z0-9]*[a-z0-9])?">
                            <span ng-if="appBuildConfig.buildname.$touched && (appBuildConfig.buildname.$error.pattern || appBuildConfig.buildname.$error.required)" class="help-block error">
                              Build config name is required and may only contain lower-case letters, numbers, and dashes. They may not start or end with a dash.
                            </span>
                          </div>
                        </dd>
                      </div>
                      <div>
                        <dt>Branch</dt>
                        <dd>  
                         <input placeholder="If empty defaults to master" ng-model="config.gitRepo.ref" type="text" id="branch-name" class="form-control">
                        </dd>
                      </div>
                      <div>
                        <dt>Jenkinsfile Path</dt>
                        <dd>
                          <input placeholder="If empty defaults to Jenkinsfile" ng-model="config.gitRepo.jenkinsFilePath" type="text" id="jenkins-path" class="form-control">
                        </dd>
                      </div>
                    </dl>
                    <button ng-click="create(appBuildConfig.$valid)" class="btn btn-primary">Create</button>
                  </form>
                </div>
              </div>
            </div>`,
  bindings: {
    created: '&'
  },
  controller: [
    '$scope',
    function($scope) {
      $scope.config = {
        name: '',
        gitRepo: {
          uri: '',
          ref: 'master',
          private: false,
          jenkinsFilePath: 'Jenkinsfile'
        }
      };

      $scope.create = function(isValid) {
        if (!isValid) {
          return;
        }

        $scope.$ctrl.created()($scope.config);
      };
    }
  ]
});
