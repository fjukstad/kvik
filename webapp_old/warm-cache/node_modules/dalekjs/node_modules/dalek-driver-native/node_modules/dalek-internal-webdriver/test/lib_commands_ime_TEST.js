'use strict';

var expect = require('chai').expect;
var Ime = require('../lib/commands/ime.js');

describe('dalek-internal-webdriver Command/Ime', function () {

  it('should get exist', function () {
    expect(Ime).to.be.ok;
  });

});
