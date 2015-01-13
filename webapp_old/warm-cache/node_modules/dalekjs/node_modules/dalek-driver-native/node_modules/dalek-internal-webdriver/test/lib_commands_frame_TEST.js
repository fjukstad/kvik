'use strict';

var expect = require('chai').expect;
var Frame = require('../lib/commands/frame.js');

describe('dalek-internal-webdriver Command/Frame', function () {

  it('should get exist', function () {
    expect(Frame).to.be.ok;
  });

});
