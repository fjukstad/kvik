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

/**
 * Timeout related Webdriver endpoints
 *
 * @param {Dalek.Internal.Webdriver} Driver Webdriver client instance
 * @return {Dalek.Internal.Webdriver} Driver Webdriver client instance
 */

module.exports = function (Driver) {
  'use strict';

  /**
   * Configure the amount of time that a particular type of operation can execute for before
   * they are aborted and a |Timeout| error is returned to the client.
   *
   * @method timeouts
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/timeouts
   * @param {GET} sessionId  ID of the session to route the command to
   * @param {POST} type The type of operation to set the timeout for. Valid values are: "script" for script timeouts, "implicit" for modifying the implicit wait timeout and "page load" for setting a page load timeout
   * @param {POST} ms The amount of time to wait, in milliseconds. This value has a lower bound of 0
   */

  Driver.addCommand({
    name: 'timeouts',
    url: '/session/:sessionId/timeouts',
    method: 'POST',
    params: ['type', 'ms']
  });

  /**
   * Set the amount of time, in milliseconds, that asynchronous scripts executed by
   * /session/:sessionId/execute_async are permitted to run before they are aborted
   * and a |Timeout| error is returned to the client.
   *
   * @method asyncScript
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/timeouts/async_script
   * @param {GET} sessionId  ID of the session to route the command to
   * @param {POST} ms The amount of time to wait, in milliseconds. This value has a lower bound of 0
   */

  Driver.addCommand({
    name: 'asyncScript',
    url: '/session/:sessionId/timeouts/async_script',
    method: 'POST',
    params: ['timeout']
  });

  /**
   * Set the amount of time the driver should wait when searching for elements.
   * When searching for a single element, the driver should poll the page until an element
   * is found or the timeout expires, whichever occurs first. When searching for multiple elements,
   * the driver should poll the page until at least one element is found or the timeout expires,
   * at which point it should return an empty list.
   *
   * @method implicitWait
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/timeouts/implicit_wait
   * @param {GET} sessionId  ID of the session to route the command to
   * @param {POST} ms The amount of time to wait, in milliseconds. This value has a lower bound of 0
   */

  Driver.addCommand({
    name: 'implicitWait',
    url: '/session/:sessionId/timeouts/implicit_wait',
    method: 'POST',
    params: ['timeout'],
    onRequest: function (params) {
      return {ms: params.timeout};
    }
  });

};
