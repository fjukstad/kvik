/*!
 *
 * Copyright (c) 2013 Sebastian Golasch
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
 * OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 *
 */

'use strict';

// ext. libs
var Q = require('q');

/**
 * Interaction related methods
 *
 * @module Driver
 * @class Interaction
 * @namespace Dalek.DriverNative.Commands
 */

var Interaction = {

  moveto: function (selector, x, y, hash) {
    var result = {element: null, xoffset: x, yoffset: y};
    if (selector !== null) {
      this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
      this.actionQueue.push(function(response) {
        var defer = Q.defer();
        var data = JSON.parse(response);
        result.element = data.value.ELEMENT;
        defer.resolve(data);
        return defer.promise;
      }.bind(this));
    }
    this.actionQueue.push(this.webdriverClient.moveto.bind(this.webdriverClient, result));
    this.actionQueue.push(this._movetoCb.bind(this, selector, x, y, hash));

    
    return this;
  },

  _movetoCb: function (selector, x, y, hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'moveto', value: selector, uuid: hash, hash: hash});
    deferred.resolve();
    return deferred.promise;
  },
  
  buttonClick: function (button, hash) {
    this.actionQueue.push(this.webdriverClient.buttonClick.bind(this.webdriverClient, button));
    this.actionQueue.push(this._buttonClickCb.bind(this, button, hash));
    return this;
  },

  _buttonClickCb: function (button, hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'buttonClick', value: button, uuid: hash, hash: hash});
    deferred.resolve();
    return deferred.promise;
  }
  
  
};

/**
 * Mixes in interaction methods
 *
 * @param {Dalek.DriverNative} DalekNative Native driver base class
 * @return {Dalek.DriverNative} DalekNative Native driver base class
 */

module.exports = function (DalekNative) {
  // mixin methods
  Object.keys(Interaction).forEach(function (fn) {
    DalekNative.prototype[fn] = Interaction[fn];
  });

  return DalekNative;
};
