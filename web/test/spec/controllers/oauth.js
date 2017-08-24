'use strict';

describe('Controller: OauthCtrl', function () {

  // load the controller's module
  beforeEach(module('mobileControlPanelApp'));

  var OauthCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    OauthCtrl = $controller('OauthCtrl', {
      $scope: scope
      // place here mocked dependencies
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(OauthCtrl.awesomeThings.length).toBe(3);
  });
});
