'use strict';

var expect = require('chai').expect;
var Timeout = require('../lib/commands/timeout.js');

describe('dalek-internal-webdriver Command/Timeout', function () {

  it('should exist', function () {
    expect(Timeout).to.be.ok;
  });

});
