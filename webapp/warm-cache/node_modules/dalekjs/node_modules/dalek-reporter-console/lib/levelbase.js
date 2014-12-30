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
var clc = require('cli-color');

/**
 * @class LevelBase
 * @constructor
 * @param {object} opts
 */

var LevelBase = function (opts) {
  this.noColor = !!(opts && opts.noColor);
  this.noSymbols = !!(opts && opts.noSymbols);
};

/**
 * @module LevelBase
 */

module.exports = LevelBase;

/**
 * Outputs a string
 *
 * @method outputTestStarted
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.echo = function (message, options) {
  options = options || {};
  var nl = options.nl === false ? '' :'\n';
  var coloredOutput = clc;
  message = options.ec === true ? message + ' ' : message;

  // should a newline be outputted before the current message
  if (options.pnl) {
    console.log('');
  }

  // output a message without any colors
  if (this.noColor) {
    process.stdout.write(message + nl);
    return this;
  }

  // check if we have a foreground color, then apply
  if (options.foreground) {
    coloredOutput = coloredOutput[options.foreground];
  }

  // check if we have a font style, then apply
  if (options.style) {
    coloredOutput = coloredOutput[options.style];
  }

  // check if we have a background color, then apply
  if (options.background) {
    coloredOutput = coloredOutput[options.background];
  }

  // check if we need to indent the string
  if (options.indent) {
    var indention = 0;
    for (var i = 0; i <= indention; i++) {
      indention += ' ';
    }
    message = indention + message;
  }

  // output the message
  process.stdout.write(coloredOutput(message + nl));
  return this;
};

/**
 * Outputs a symbol
 *
 * @method outputTestStarted
 * @param {object} data
 * @return {string} symbol
 */

LevelBase.prototype.symbol = function (symbol) {
  if (this.noSymbols || process.platform === 'win32') {
    return symbol;
  }

  // switch to a nicer UTF-8 symbol
  switch (symbol) {
  case '*':
    symbol = '✔';
    break;
  case 'x':
    symbol = '✘';
    break;
  case '>':
    symbol = '▶';
    break;
  case '->':
    symbol = '↝';
    break;
  case '>>':
    symbol = '>>';
    break;
  case '<>':
    symbol = '☁';
    break;
  }

  return symbol;
};

/**
 * Does nothing
 *
 * @method outputRunnerStarted
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputRunnerStarted = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputTestFinished
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputTestFinished = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputRunnerFinished
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputRunnerFinished = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputAssertionResult
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputAssertionResult = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputRunBrowser
 * @param {object} browser
 * @chainable
 */

LevelBase.prototype.outputRunBrowser = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputAction
 * @chainable
 */

LevelBase.prototype.outputAction = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputAssertionExpectation
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputAssertionExpectation = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputAction
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputAction = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputTestStarted
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputTestStarted = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputLogUser
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputLogUser = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputBrowserVersion
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputBrowserVersion = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputOSVersion
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputOSVersion = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputReportWritten
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputReportWritten = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputWarning
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputWarning = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method outputError
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.outputError = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method remoteConnectionEstablished
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.remoteConnectionEstablished = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method remoteConnectionClosed
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.remoteConnectionClosed = function () {
  return this;
};

/**
 * Does nothing
 *
 * @method remoteConnectionReady
 * @param {object} data
 * @chainable
 */

LevelBase.prototype.remoteConnectionReady = function () {
  return this;
};
