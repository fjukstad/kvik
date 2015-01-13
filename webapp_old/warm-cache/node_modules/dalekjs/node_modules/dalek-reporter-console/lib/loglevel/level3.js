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

// int. libs
var Log = require('./level2');

/**
 * @class LogLevel3
 * @constructor
 * @param {object} opts
 */

var LogLevel3 = Log;

/**
 * @module LogLevel3
 */

module.exports = LogLevel3;

/**
 * Outputs a test action
 *
 * @method outputAction
 * @param {object} data
 * @return {LogLevel3}
 */

LogLevel3.prototype.outputAction = function (data) {
  // check for `nonsense` values
  if (data.value === null || typeof data.value === 'undefined') {
    data.value = '';
  }

  this.echo(this.symbol('>'), {nl: false, ec: true, foreground: 'yellow'});
  this.echo(data.type.toUpperCase() + ' ' + data.value, {foreground: 'yellow'});
  return this;
};

/**
 * Outputs the browser version
 *
 * @method outputBrowserVersion
 * @param {object} data
 * @chainable
 */

LogLevel3.prototype.outputBrowserVersion = function (data) {
  this.echo('Browser Version:', {nl: false, ec: true, foreground: 'white'});
  this.echo(data.version, {foreground: 'yellowBright'});
  return this;
};

/**
 * Outputs operating system information
 *
 * @method outputOSVersion
 * @param {object} data
 * @chainable
 */

LogLevel3.prototype.outputOSVersion = function (data) {
  this.echo('OS:', {nl: false, ec: true, foreground: 'white'});
  this.echo(data.os.name + ' ' + (data.os.version || '') + ' ' + (data.os.arch || ''), {foreground: 'yellowBright'});
  return this;
};

/**
 * Outputs operating system information
 *
 * @method outputOSVersion
 * @param {object} data
 * @chainable
 */

LogLevel3.prototype.outputReportWritten = function (data) {
  this.echo(this.symbol('->'), {nl: true, ec: true, foreground: 'yellow'});
  this.echo('Report type "' + data.type + '" written to "' + data.dest + '"', {foreground: 'yellow'});
  return this;
};
