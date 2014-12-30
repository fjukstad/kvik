'use strict';

var expect = require('chai').expect;
var Session = require('../lib/commands/session.js');

describe('dalek-internal-webdriver Command/Session', function () {

  it('should get exist', function () {
    expect(Session).to.be.ok;
  });

});
