'use strict';

/**
 * @ngdoc component
 * @name mcp.component:create-app-config
 * @description
 * # create-app-config
 */
angular.module('mobileControlPanelApp').component('createAppConfig', {
  template: `<div>
              <div class="row edit app-config">
                <div class="col-lg-6">
                  <h3>Create</h3>
                  <form novalidate class="form-horizontal" name="appBuildConfig">
                    <dl class="dl-horizontal left">
                      <div>
                        <dt>Name</dt>
                        <dd>
                          <div ng-class="{'has-error': appBuildConfig.buildname.$touched && (appBuildConfig.buildname.$error.pattern || appBuildConfig.buildname.$error.required)}">
                            <input ng-model="config.name" name="buildname" type="text" id="build-name" class="form-control" required pattern="[a-z0-9]([-a-z0-9]*[a-z0-9])?">
                          </div>
                        </dd>
                      </div>
                      <div ng-class="{'has-error': appBuildConfig.repoUri.$touched && appBuildConfig.repoUri.$error.required}">
                        <dt>Repo URL</dt>
                        <dd>
                          <input ng-model="config.gitRepo.uri" name="repoUri" type="text" id="repo-uri" class="form-control" required>
                        </dd>
                      </div>
                      <div>
                        <dt>Branch</dt>
                        <dd>
                         <input ng-model="config.gitRepo.ref" type="text" id="branch-name" class="form-control">
                        </dd>
                      </div>
                      <div>
                        <dt>Jenkinsfile Path</dt>
                        <dd>
                          <input ng-model="config.gitRepo.jenkinsFilePath" type="text" id="jenkins-path" class="form-control">
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
