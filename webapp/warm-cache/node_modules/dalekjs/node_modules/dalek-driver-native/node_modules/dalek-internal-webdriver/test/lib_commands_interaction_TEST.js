'use strict';

var expect = require('chai').expect;
var Interaction = require('../lib/commands/interaction.js');

describe('dalek-internal-webdriver Command/Interaction', function () {

  it('should get exist', function () {
    expect(Interaction).to.be.ok;
  });

});
