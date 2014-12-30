'use strict';

var expect = require('chai').expect;
var Driver = require('../index');

describe('dalek-internal-webdriver', function() {

  it('should get exist', function(){
    expect(Driver).to.be.ok;
  });

});
