'use strict';

var expect = require('chai').expect;
var WebDriver = require('../lib/webdriver.js');

describe('dalek-internal-webdriver WebDriver', function() {

  it('should get exist', function(){
    expect(WebDriver).to.be.ok;
  });

});
