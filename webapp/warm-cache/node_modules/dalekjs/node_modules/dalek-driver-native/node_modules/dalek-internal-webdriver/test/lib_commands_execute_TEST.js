'use strict';

var expect = require('chai').expect;
var Execute = require('../lib/commands/execute.js');

describe('dalek-internal-webdriver Command/Execute', function () {

  it('should get exist', function () {
    expect(Execute).to.be.ok;
  });

});
