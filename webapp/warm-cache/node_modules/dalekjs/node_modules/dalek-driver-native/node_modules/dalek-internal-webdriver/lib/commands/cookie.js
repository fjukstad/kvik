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
 * Cookie related Webdriver endpoints
 *
 * @param {Dalek.Internal.Webdriver} Driver Webdriver client instance
 * @return {Dalek.Internal.Webdriver} Driver Webdriver client instance
 */

module.exports = function (Driver) {
  'use strict';

  /**
   * Retrieve all cookies visible to the current page.
   *
   * @method getCookies
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/cookie
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'getCookies',
    url: '/session/:sessionId/cookie',
    method: 'GET'
  });

  /**
   * Retrieve a cookies by its name.
   *
   * @method getCookie
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/cookie/:name
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} name Name of the cookie
   */

  Driver.addCommand({
    name: 'getCookie',
    url: '/session/:sessionId/cookie',
    method: 'GET'
  });

  /**
   * Set a cookie.
   * If the cookie path is not specified, it should be set to "/".
   * Likewise, if the domain is omitted,
   * it should default to the current page's domain.
   *
   * @method setCookie
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/cookie
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} cookie The cookie object
   */

  Driver.addCommand({
    name: 'setCookie',
    url: '/session/:sessionId/cookie',
    method: 'POST',
    params: ['cookie']
  });

  return Driver;
};
