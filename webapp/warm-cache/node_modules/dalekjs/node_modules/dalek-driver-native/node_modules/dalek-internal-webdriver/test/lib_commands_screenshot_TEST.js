'use strict';

var expect = require('chai').expect;
var Screenshot = require('../lib/commands/screenshot.js');

describe('dalek-internal-webdriver Command/Screenshot', function () {

  it('should get exist', function () {
    expect(Screenshot).to.be.ok;
  });

});
