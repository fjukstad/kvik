'use strict';

var expect = require('chai').expect;
var Page = require('../lib/commands/page.js');

describe('dalek-internal-webdriver Command/Page', function () {

  it('should get exist', function () {
    expect(Page).to.be.ok;
  });

});
