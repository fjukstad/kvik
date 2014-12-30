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

/**
 * Stores the injected options as an object property
 *
 * @param {object} opts Webdriver options
 * @constructor
 */

var WebDriver = function (opts, events) {
  this.events = events;
  this.opts = opts;
};

/**
 * Webdriver base class
 *
 * @namespace Internal
 * @module DalekJS
 * @class Webdriver
 */

WebDriver.prototype = {

  /**
   * Local options, used to store data
   * like the current sessionId or similar
   */

  options: {},

  /**
   * Checks if we have a valid webdriver session
   *
   * @method hasSession
   * @return {bool} Is valid session
   */

  hasSession: function () {
    return !!this.options.sessionId;
  },

  /**
   * Closes the current webdriver session
   *
   * @method hasSession
   * @chainable
   */

  closeSession: function () {
    delete this.options.sessionId;
    return this;
  }
};

// export the module
module.exports = WebDriver;
