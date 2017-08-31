'use strict';

describe('Service: mcpApi', function () {

  // load the service's module
  beforeEach(module('mobileControlPanelApp'));

  // instantiate service
  var mcpApi;
  beforeEach(inject(function (_mcpApi_) {
    mcpApi = _mcpApi_;
  }));

  it('should do something', function () {
    expect(!!mcpApi).toBe(true);
  });

});
