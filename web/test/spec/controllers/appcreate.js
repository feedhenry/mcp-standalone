'use strict';

describe('Controller: AppcreateCtrl', function () {

  // load the controller's module
  beforeEach(module('mobileControlPanelApp'));

  var AppcreateCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    AppcreateCtrl = $controller('AppcreateCtrl', {
      $scope: scope
      // place here mocked dependencies
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(AppcreateCtrl.awesomeThings.length).toBe(3);
  });
});
