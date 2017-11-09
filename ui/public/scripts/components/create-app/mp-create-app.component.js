'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-create-app
 * @description
 * # mp-create-app
 */
angular.module('mobileControlPanelApp').component('mpCreateApp', {
  template: `<form name="createAppForm" novalidate="" class="ng-pristine ng-invalid ng-invalid-required ng-valid-pattern ng-valid-minlength ng-valid-maxlength">
              <fieldset ng-disabled="disableInputs">
                <div class="form-group">
                  <label for="name" class="required">Name</label>
                  <span ng-class="{'has-error': (createAppForm.name.$error.pattern &amp;&amp; createProjectForm.name.$touched) || nameTaken}">
                    <input class="form-control input-lg ng-pristine ng-empty ng-invalid ng-invalid-required ng-valid-pattern ng-valid-minlength ng-valid-maxlength ng-touched" name="name" id="name" placeholder="my-app" type="text" required="" take-focus="" minlength="2" maxlength="63" pattern="[a-z0-9]([-a-z0-9]*[a-z0-9])?" aria-describedby="nameHelp" ng-model="$ctrl.app.name" ng-model-options="{ updateOn: 'default blur' }" ng-change="nameTaken = false" autocorrect="off" autocapitalize="off" spellcheck="false" style="">
                  </span>
                  <div>
                    <span class="help-block">A unique name for your app.</span>
                  </div>
                  <div class="has-error">
                    <!-- ngIf: createProjectForm.name.$error.required && createProjectForm.name.$dirty -->
                  </div>
                  <div class="has-error">
                    <!-- ngIf: createProjectForm.name.$error.minlength && createProjectForm.name.$touched -->
                  </div>
                  <div class="has-error">
                    <!-- ngIf: createProjectForm.name.$error.pattern && createProjectForm.name.$touched -->
                  </div>
                  <div class="has-error">
                    <!-- ngIf: nameTaken -->
                  </div>
                </div>
                <div class="form-group">
                  <label for="description">Description</label>
                  <input class="form-control input-lg ng-pristine ng-empty" name="description" id="description" placeholder="description" type="text" ng-model="$ctrl.app.description" autocorrect="off" autocapitalize="off" spellcheck="false" style="">
                </div>  
                <div class="form-group">
                  <label for="clientType">Client Type</label>
                  <select class="form-control" name="clientType" ng-model="$ctrl.app.clientType">
                      <option value="">---Please select---</option>
                      <option value="android">Android</option>
                      <option value="iOS">iOS</option>
                      <option value="cordova">Cordova</option>
                  </select>                
                </div>
            
                <!--<div class="form-group">
                  <label for="description">Description</label>
                  <textarea class="form-control input-lg ng-pristine ng-untouched ng-valid ng-empty" name="description" id="description" placeholder="A short description." ng-model="description"></textarea>
                </div>-->
            
                <div class="button-group">
                  <button type="submit" class="btn btn-primary btn-lg" ng-class="{'dialog-btn': isDialog}" ng-click="$ctrl.created()($ctrl.app)" ng-disabled="createAppForm.$invalid || nameTaken || disableInputs" value="" disabled="disabled">
                    Create
                  </button>
                  <button class="btn btn-default btn-lg" ng-class="{'dialog-btn': isDialog}" ng-click="$ctrl.cancelled()()">
                    Cancel
                  </button>
                </div>
              </fieldset>
            </form>`,
  bindings: {
    created: '&',
    cancelled: '&'
  },
  controller: [
    function() {
      this.app = { clientType: '' };
    }
  ]
});
