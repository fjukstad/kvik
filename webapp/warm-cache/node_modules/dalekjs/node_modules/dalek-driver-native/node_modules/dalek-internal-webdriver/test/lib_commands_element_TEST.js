'use strict';

var expect = require('chai').expect;
var Element = require('../lib/commands/element.js');

describe('dalek-internal-webdriver Command/Element', function () {

  it('should get exist', function () {
    expect(Element).to.be.ok;
  });

});
