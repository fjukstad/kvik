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
var Log = require('./level1');

/**
 * @class LogLevel2
 * @constructor
 * @param {object} opts
 */

var LogLevel2 = Log;

/**
 * @module LogLevel2
 */

module.exports = LogLevel2;

/**
 * Outputs in which browser env the tests will run
 *
 * @method outputRunBrowser
 * @param {object} data
 * @return {LogLevel2}
 */

LogLevel2.prototype.outputRunBrowser = function (browser) {
  this.echo('Running Browser:', {nl: false, ec: true, foreground: 'white'});
  this.echo(browser, {foreground: 'yellowBright'});
  return this;
};

/**
 * Outputs a message when the test starts
 *
 * @method outputTestStarted
 * @param {object} data
 * @return {LogLevel2}
 */

LogLevel2.prototype.outputTestStarted = function (data) {
  this.echo('RUNNING TEST - ', {nl: false, pnl: true, foreground: 'magentaBright'});
  this.echo('"' + data.name + '"', {foreground: 'cyanBright'});
  return this;
};

/**
 * Outputs a message when all assertions of a test has been run
 *
 * @method outputAssertionResult
 * @param {object} data
 * @return {LogLevel2}
 */

LogLevel2.prototype.outputAssertionResult = function (data) {
  if (data.success) {
    this.echo(this.symbol('*'), {nl: false, ec: true, foreground: 'greenBright'});
    this.echo(data.type.toUpperCase(), {nl: false, ec: true, foreground: 'magentaBright'});
    if (data.message) {
      this.echo(data.message, {foreground: 'greenBright'});
    } else {
      this.echo('');
    }
  } else {
    this.echo(this.symbol('x'), {nl: false, ec: true, foreground: 'redBright'});
    this.echo(data.type.toUpperCase(), {foreground: 'redBright'});
    this.echo('EXPECTED: ' + data.expected, {indent: 4, foreground: 'magentaBright'});
    this.echo('FOUND: ' + data.value, {indent: 4, foreground: 'magentaBright'});
    if (data.message) {
      this.echo('MESSAGE: ' + data.message, {indent: 4, foreground: 'magentaBright'});
    }
  }
  return this;
};

/**
 * Outputs a message when all assertions of a test has been run
 *
 * @method outputAssertionExpectation
 * @param {object} data
 * @return {LogLevel2}
 */

LogLevel2.prototype.outputAssertionExpectation = function (data) {
  if (data.status) {
    this.echo(this.symbol('*'), {nl: false, ec: true, foreground: 'greenBright'});
    this.echo(data.expected, {nl: false, ec: true, foreground: 'magentaBright'});
    this.echo('Assertions run', {foreground: 'greenBright'});
  } else {
    if (data.expected !== data.run) {
      this.echo(this.symbol('x'), {nl: false, ec: true, foreground: 'redBright'});
      this.echo(data.expected, {nl: false, ec: true, foreground: 'magentaBright'});
      this.echo('Assertions run', {foreground: 'redBright'});
      this.echo('EXPECTED: ' + data.expected, {indent: 4, foreground: 'magentaBright'});
      this.echo('RUN: ' + data.run, {indent: 4, foreground: 'magentaBright'});
    }
  }
  return this;
};

/**
 * Outputs a message when a test has been finished
 *
 * @method outputTestFinished
 * @param {object} data
 * @return {LogLevel2}
 */

LogLevel2.prototype.outputTestFinished = function (data) {
  this.echo(this.symbol((data.status ? '*' : 'x')), {nl: false, ec: true, foreground: (data.status ? 'greenBright' : 'redBright')});
  this.echo('TEST -', {nl: false, ec: true, foreground: 'magentaBright'});
  this.echo('"' + data.name + '"', {nl: false, ec: true, foreground: 'cyanBright'});
  this.echo((data.status ? 'SUCCEEDED' : 'FAILED'), {foreground: (data.status ? 'greenBright' : 'redBright')});
  return this;
};
