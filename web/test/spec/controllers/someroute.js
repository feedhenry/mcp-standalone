'use strict';

describe('Controller: SomerouteCtrl', function () {

  // load the controller's module
  beforeEach(module('mobileControlPanelApp'));

  var SomerouteCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    SomerouteCtrl = $controller('SomerouteCtrl', {
      $scope: scope
      // place here mocked dependencies
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(SomerouteCtrl.awesomeThings.length).toBe(3);
  });
});
