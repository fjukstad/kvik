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
 * Element related Webdriver endpoints
 *
 * @param {Dalek.Internal.Webdriver} Driver Webdriver client instance
 * @return {Dalek.Internal.Webdriver} Driver Webdriver client instance
 */

module.exports = function (Driver) {
  'use strict';

  /**
   * Search for an element on the page, starting from the document root.
   * The located element will be returned as a WebElement JSON object.
   *
   * @method element
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/element
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} selector The The search target.
   */

  Driver.addCommand({
    name: 'element',
    url: '/session/:sessionId/element',
    method: 'POST',
    params: ['selector'],
    onRequest: function (params) {
      var type = 'css selector';
      return {using: type, value: params.selector};
    },
    onResponse: function (request, remote, options, deferred, data) {
      this.options.id = JSON.parse(data).value.ELEMENT;
      deferred.resolve(data);
    },
    onError: function (request, remote, options, deferred, data) {
      data = JSON.parse(data);
      deferred.resolve(JSON.stringify({'sessionId': data.sessionId, value: -1}));
    }
  });

  /**
   * Search for multiple elements on the page, starting from the document root.
   * The located element will be returned as a WebElement JSON object.
   *
   * @method elements
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/elements
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} selector The The search target.
   */

  Driver.addCommand({
    name: 'elements',
    url: '/session/:sessionId/elements',
    method: 'POST',
    params: ['selector'],
    onRequest: function (params) {
      var type = 'css selector';
      return {using: type, value: params.selector};
    },
    onResponse: function (request, remote, options, deferred, data) {
      deferred.resolve(data);
    },
    onError: function (request, remote, options, deferred, data) {
      data = JSON.parse(data);
      deferred.resolve(JSON.stringify({'sessionId': data.sessionId, value: -1}));
    }
  });

  /**
   * Get the element on the page that currently has focus.
   * The element will be returned as a WebElement JSON object.
   *
   * @method active
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/element/:id/active
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} id ID of the element to route the command to
   */

  Driver.addCommand({
    url: '/session/:sessionId/element/:id/active',
    method: 'GET'
  });

  /**
   * Get the element on the page that currently has focus.
   * The element will be returned as a WebElement JSON object.
   *
   * @method elementInfo
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/element/:id
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   */

  Driver.addCommand({
    name: 'elementInfo',
    url: '/session/:sessionId/element/:id',
    method: 'GET'
  });

  /**
   * Search for an element on the page, starting from the identified element.
   * The located element will be returned as a WebElement JSON object.
   * The table below lists the locator strategies that each server should support.
   * Each locator must return the first matching element located in the DOM.
   *
   * @method childElement
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/:id/element
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   * @param {POST} using The locator strategy to use. // Not yet supported
   * @param {POST} value The The search target.
   */

  Driver.addCommand({
    name: 'childElement',
    url: '/session/:sessionId/element/:id/element',
    method: 'GET'
  });

  /**
   * Search for an element on the page, starting from the identified element.
   * The located element will be returned as a WebElement JSON object.
   * The table below lists the locator strategies that each server should support.
   * Each locator must return the first matching element located in the DOM.
   *
   * @method element
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/:id/element
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   * @param {POST} using The locator strategy to use. // Not yet supported
   * @param {POST} value The The search target.
   */

  Driver.addCommand({
    url: '/session/:sessionId/element/:id/elements',
    method: 'GET'
  });

  /**
   * Click on an element.
   *
   * @method click
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/element/:id/click
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} id ID of the element to route the command to
   */

  Driver.addCommand({
    name: 'click',
    url: '/session/:sessionId/element/:id/click',
    method: 'POST'
  });

  /**
   * Submit a FORM element.
   * The submit command may also be applied to any element that is a descendant of a FORM element.
   *
   * @method submit
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/element/:id/submit
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} id ID of the element to route the command to
   */

  Driver.addCommand({
    name: 'submit',
    url: '/session/:sessionId/element/:id/submit',
    method: 'POST'
  });

  /**
   * Returns the visible text for the element.
   *
   * @method text
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/element/:id/text
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} id ID of the element to route the command to
   */

  Driver.addCommand({
    name: 'text',
    url: '/session/:sessionId/element/:id/text',
    method: 'GET'
  });

  /**
   * Send a sequence of key strokes to an element
   *
   * @method val
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/element/:id/value
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} id ID of the element to route the command to
   * @param {POST} text The keys sequence to be sent
   */

  Driver.addCommand({
    name: 'val',
    url: '/session/:sessionId/element/:id/value',
    method: 'POST',
    params: ['text'],
    onRequest: function (params) {
      return {value: params.text.split('')};
    }
  });

  /**
   * Query for an element's tag name
   *
   * @method name
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/element/:id/name
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   */

  Driver.addCommand({
    name: 'name',
    url: '/session/:sessionId/element/:id/name',
    method: 'GET'
  });

  /**
   * Clear a TEXTAREA or text INPUT element's value
   *
   * @method clear
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/element/:id/clear
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   */

  Driver.addCommand({
    name: 'clear',
    url: '/session/:sessionId/element/:id/clear',
    method: 'POST'
  });

  /**
   * Determine if an OPTION element, or an INPUT element of type checkbox or radiobutton is currently selected
   *
   * @method selected
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/element/:id/selected
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   */

  Driver.addCommand({
    name: 'selected',
    url: '/session/:sessionId/element/:id/selected',
    method: 'GET'
  });

  /**
   * Determine if an element is currently enabled
   *
   * @method enabled
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/element/:id/enabled
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   */

  Driver.addCommand({
    name: 'enabled',
    url: '/session/:sessionId/element/:id/enabled',
    method: 'GET'
  });

  /**
   * Get the value of an element's attribute.
   *
   * @method getAttribute
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/element/:id/attribute/:name
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   * @param {GET} attr Attribute that should be fetched
   */

  Driver.addCommand({
    name: 'getAttribute',
    url: '/session/:sessionId/element/:id/attribute/:name',
    method: 'GET',
    params: ['name'],
    onRequest: function (params) {
      this.options.name = params.name;
      return null;
    }
  });

  /**
   * Test if two element IDs refer to the same DOM element
   *
   * @method equals
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/element/:id/equals/:other
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   * @param {GET} other ID of the element to compare
   */

  Driver.addCommand({
    url: '/session/:sessionId/element/:id/equals/:other',
    method: 'GET'
  });

  /**
   * Determine an element's location on the page.
   * The point (0, 0) refers to the upper-left corner of the page.
   * The element's coordinates are returned as a JSON object with x and y properties.
   *
   * @method displayed
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/element/:id/displayed
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   */

  Driver.addCommand({
    name: 'displayed',
    url: '/session/:sessionId/element/:id/displayed',
    method: 'GET'
  });

  /**
   * Determine an element's location on the page. The point (0, 0) refers to the upper-left corner of the page.
   * The element's coordinates are returned as a JSON object with x and y properties.
   *
   * @method location
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/element/:id/location
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   */

  Driver.addCommand({
    name: 'location',
    url: '/session/:sessionId/element/:id/location',
    method: 'GET'
  });

  /**
   * Determine an element's location on the screen once it has been scrolled into view.
   * Note: This is considered an internal command and should only be used to determine an element's location for correctly generating native events.
   *
   * @method locationInView
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/element/:id/location_in_view
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   */

  Driver.addCommand({
    name: 'locationInView',
    url: '/session/:sessionId/element/:id/location_in_view',
    method: 'GET'
  });

  /**
   * Determine an element's size in pixels.
   * The size will be returned as a JSON object with width and height properties.
   *
   * @method size
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/element/:id/size
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   */

  Driver.addCommand({
    name: 'size',
    url: '/session/:sessionId/element/:id/size',
    method: 'GET'
  });

  /**
   * Query the value of an element's computed CSS property.
   * The CSS property to query should be specified using the CSS property name,
   * not the JavaScript property name (e.g. background-color instead of backgroundColor).
   *
   * @method cssProperty
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/element/:id/size
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} elementId ID of the element to route the command to
   * @param {GET} propertyName Name of the css property to fetch
   */

  Driver.addCommand({
    name: 'cssProperty',
    url: '/session/:sessionId/element/:id/css/:propertyName',
    method: 'GET',
    params: ['propertyName'],
    onRequest: function (params) {
      this.options.propertyName = params.propertyName;
      return null;
    }
  });

  /**
   * Get the log for a given log type.
   * Log buffer is reset after each request.
   *
   * @method sendKeys
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/keys
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} id ID of the element to route the command to
   * @param {POST} text The keys sequence to be sent
   */

  Driver.addCommand({
    name: 'sendKeys',
    url: '/session/:sessionId/keys',
    method: 'POST',
    params: ['text'],
    onRequest: function (params) {
      return {value: params.text.split('')};
    }
  });

};
