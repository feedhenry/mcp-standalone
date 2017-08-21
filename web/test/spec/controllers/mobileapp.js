'use strict';

describe('Controller: MobileappCtrl', function () {

  // load the controller's module
  beforeEach(module('mobileControlPanelApp'));

  var MobileappCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    MobileappCtrl = $controller('MobileappCtrl', {
      $scope: scope
      // place here mocked dependencies
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(MobileappCtrl.awesomeThings.length).toBe(3);
  });
});
