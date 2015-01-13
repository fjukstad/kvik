'use strict';

var expect = require('chai').expect;
var PhantomJSDriver = require('../index');

describe('dalek-browser-phantomjs', function() {

  it('should get default webdriver port', function(){
    expect(PhantomJSDriver.port).to.equal(9001);
  });

});
