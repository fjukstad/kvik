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
 * Element related methods
 *
 * @module Driver
 * @class Element
 * @namespace Dalek.DriverNative.Commands
 */

var Element = {

  /**
   * Checks if an element exists
   *
   * @method exists
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  exists: function (selector, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this._existsCb.bind(this, selector, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `exists` call
   *
   * @method _existsCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the exists call
   * @return {object} promise Exists promise
   * @private
   */

  _existsCb: function (selector, hash, result) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'exists', selector: selector, hash: hash, value: (JSON.parse(result).value === -1 || JSON.parse(result).status === 7 ? 'false' : 'true')});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Checks if an element is visible
   *
   * @method visible
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  visible: function (selector, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.displayed.bind(this.webdriverClient, selector));
    this.actionQueue.push(this._visibleCb.bind(this, selector, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `visible` call
   *
   * @method _visibleCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the visible call
   * @return {object} promise Visible promise
   * @private
   */

  _visibleCb: function (selector, hash, result) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'visible', selector: selector, hash: hash, value: JSON.parse(result).value});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Checks if an element has the expected text
   *
   * @method text
   * @param {string} selector Selector expression to find the element
   * @param {string} expected The expected text content
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  text: function (selector, expected, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.text.bind(this.webdriverClient, selector));
    this.actionQueue.push(this._textCb.bind(this, selector, hash, expected));
    return this;
  },

  /**
   * Sends out an event with the results of the `text` call
   *
   * @method _textCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} expected The expected text content
   * @param {string} result Serialized JSON with the reuslts of the text call
   * @return {object} Promise
   * @private
   */

  _textCb: function (selector, hash, expected, result) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'text', hash: hash, expected: expected, selector: selector, value: JSON.parse(result).value});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Checks if an element has the expected css property/value pairs
   *
   * @method css
   * @param {string} selector Selector expression to find the element
   * @param {string} property The css property to check
   * @param {string} expected The expected css value
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  css: function (selector, property, expected, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.cssProperty.bind(this.webdriverClient, property));
    this.actionQueue.push(this._cssCb.bind(this, selector, property, hash, expected));
    return this;
  },

  /**
   * Sends out an event with the results of the `css` call
   *
   * @method _cssCb
   * @param {string} selector Selector expression to find the element
   * @param {string} property The css property to check
   * @param {string} expected The expected css value
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the text call
   * @return {object} Promise
   * @private
   */
  
  _cssCb: function (selector, property, hash, expected, result) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'css', hash: hash, property: property, expected: expected, selector: selector, value: JSON.parse(result).value});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Checks th width of an element
   *
   * @method width
   * @param {string} selector Selector expression to find the element
   * @param {string} expected The expected width value
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  width: function (selector, expected, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.size.bind(this.webdriverClient));
    this.actionQueue.push(this._widthCb.bind(this, selector, hash, expected));
    return this;
  },

  /**
   * Sends out an event with the results of the `width` call
   *
   * @method _widthCb
   * @param {string} selector Selector expression to find the element
   * @param {string} expected The expected width value
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the text call
   * @return {object} Promise
   * @private
   */

  _widthCb: function (selector, hash, expected, result) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'width', hash: hash, expected: expected, selector: selector, value: JSON.parse(result).value.width});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Checks th height of an element
   *
   * @method height
   * @param {string} selector Selector expression to find the element
   * @param {string} expected The expected height value
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  height: function (selector, expected, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.size.bind(this.webdriverClient));
    this.actionQueue.push(this._heightCb.bind(this, selector, hash, expected));
    return this;
  },

  /**
   * Sends out an event with the results of the `height` call
   *
   * @method _heightCb
   * @param {string} selector Selector expression to find the element
   * @param {string} expected The expected height value
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the text call
   * @return {object} Promise
   * @private
   */

  _heightCb: function (selector, hash, expected, result) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'height', hash: hash, expected: expected, selector: selector, value: JSON.parse(result).value.height});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Checks if an element has an attribute with the expected value
   *
   * @method attribute
   * @param {string} selector Selector expression to find the element
   * @param {string} attribute The attribute that should be checked
   * @param {string} expected The expected text content
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  attribute: function (selector, attribute, expected, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.getAttribute.bind(this.webdriverClient, attribute));
    this.actionQueue.push(this._attributeCb.bind(this, selector, hash, attribute, expected));
    return this;
  },

  /**
   * Sends out an event with the results of the `attribute` call
   *
   * @method _attributeCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} attribute The attribute that should be checked
   * @param {string} expected The expected attribute content
   * @param {string} result Serialized JSON with the results of the attribute call
   * @return {object} promise Attribute promise
   * @private
   */

  _attributeCb: function (selector, hash, attribute, expected, result) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'attribute', selector: selector, hash: hash, expected: expected, value: JSON.parse(result).value });
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Checks if an element has the expected value
   *
   * @method val
   * @param {string} selector Selector expression to find the element
   * @param {string} expected The expected content
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  val: function (selector, expected, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.getAttribute.bind(this.webdriverClient, 'value'));
    this.actionQueue.push(this._valCb.bind(this, selector, hash, expected));
    return this;
  },

  /**
   * Sends out an event with the results of the `val` call
   *
   * @method _valCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} expected The expected content
   * @param {string} result Serialized JSON with the results of the val call
   * @return {object} Promise
   * @private
   */

  _valCb: function (selector, hash, expected, result) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'val', selector: selector, hash: hash, expected: expected, value: JSON.parse(result).value });
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Checks if an element is selected
   *
   * @method selected
   * @param {string} selector Selector expression to find the element
   * @param {string} expected The expected content
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  selected: function (selector, expected, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.selected.bind(this.webdriverClient));
    this.actionQueue.push(this._selectedCb.bind(this, selector, hash, expected));
    return this;
  },

  /**
   * Sends out an event with the results of the `selected` call
   *
   * @method _selectedCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} expected The expected content
   * @param {string} result Serialized JSON with the results of the selected call
   * @return {object} Promise
   * @private
   */

  _selectedCb: function (selector, hash, expected, result) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'selected', selector: selector, hash: hash, expected: expected, value: JSON.parse(result).value });
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Checks if an element is enabled
   *
   * @method enabled
   * @param {string} selector Selector expression to find the element
   * @param {string} expected The expected content
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  enabled: function (selector, expected, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.enabled.bind(this.webdriverClient));
    this.actionQueue.push(this._enabledCb.bind(this, selector, hash, expected));
    return this;
  },

  /**
   * Sends out an event with the results of the `enabled` call
   *
   * @method _enabledCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} expected The expected content
   * @param {string} result Serialized JSON with the results of the selected call
   * @return {object} Promise
   * @private
   */

  _enabledCb: function (selector, hash, expected, result) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'enabled', selector: selector, hash: hash, expected: expected, value: JSON.parse(result).value });
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Submits a form
   *
   * @method submit
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @chainable
   */

  submit: function (selector, hash, uuid) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.submit.bind(this.webdriverClient));
    this.actionQueue.push(this._submitCb.bind(this, selector, hash, uuid));
    return this;
  },

  /**
   * Sends out an event with the results of the `submit` call
   *
   * @method _submitCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @return {object} promise Click promise
   * @private
   */

  _submitCb: function (selector, hash, uuid) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'submit', value: selector, uuid: uuid, hash: hash});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Clicks an element
   *
   * @method click
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @chainable
   */

  click: function (selector, hash, uuid) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.click.bind(this.webdriverClient));
    this.actionQueue.push(this._clickCb.bind(this, selector, hash, uuid));
    return this;
  },

  /**
   * Sends out an event with the results of the `click` call
   *
   * @method _clickCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @return {object} promise Click promise
   * @private
   */

  _clickCb: function (selector, hash, uuid) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'click', value: selector, uuid: uuid, hash: hash});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Scrolls from an element to a location defined in pixels
   *
   * @method scroll
   * @param {string} selector Selector expression to find the element
   * @param {object} options X offset, Y offset, Speed
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @chainable
   */

  scroll: function (selector, options, hash, uuid) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector, options));
    this.actionQueue.push(this.webdriverClient.scroll.bind(this.webdriverClient));
    this.actionQueue.push(this._clickCb.bind(this, selector, options, hash, uuid));
    return this;
  },

  /**
   * Sends out an event with the results of the `scroll` call
   *
   * @method _scrollCb
   * @param {string} selector Selector expression to find the element
   * @param {object} options X offset, Y offset, Speed
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @return {object} promise Scroll promise
   * @private
   */

  _scrollCb: function (selector, options, hash, uuid) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'scroll', value: selector, uuid: uuid, hash: hash});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Clicks an element
   *
   * @method click
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  type: function (selector, keystrokes, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.val.bind(this.webdriverClient, keystrokes));
    this.actionQueue.push(this._typeCb.bind(this, selector, keystrokes, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `type` call
   *
   * @method _typeCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @return {object} promise Type promise
   * @private
   */

  _typeCb: function (selector, keystrokes, hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'type', value: selector, keystrokes: keystrokes, uuid: hash, hash: hash});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Set a value to an Element
   *
   * @method setValue
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  setValue: function (selector, keystrokes, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.clear.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.val.bind(this.webdriverClient, keystrokes));
    this.actionQueue.push(this._setValueCb.bind(this, selector, keystrokes, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `setvalue` call
   *
   * @method _setValueCb
   * @param {string} selector Selector expression to find the element
   * @param {string} keystrokes Value to set
   * @param {string} hash Unique hash of that fn call
   * @return {object} promise Type promise
   * @private
   */

  _setValueCb: function (selector, keystrokes, hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'setValue', value: selector, keystrokes: keystrokes, uuid: hash, hash: hash});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Sends keys to an element (whether or not it is an input)
   *
   * @method sendKeys
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  sendKeys: function (selector, keystrokes, hash) {
    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    this.actionQueue.push(this.webdriverClient.sendKeys.bind(this.webdriverClient, keystrokes));
    this.actionQueue.push(this._sendKeysCb.bind(this, selector, keystrokes, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `sendKeys` call
   *
   * @method _sendKeysCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @return {object} promise Type promise
   * @private
   */

  _sendKeysCb: function (selector, keystrokes, hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'sendKeys', value: selector, keystrokes: keystrokes, uuid: hash, hash: hash});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Wait for an element for a specific amount of time
   *
   * @method waitForElement
   * @param {string} selector Selector expression to find the element
   * @param {integer} timeout Time to wait in ms
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @chainable
   */

  waitForElement: function (selector, timeout, hash, uuid) {
    this.actionQueue.push(this.webdriverClient.implicitWait.bind(this.webdriverClient, timeout));
    this.actionQueue.push(this._waitForElementCb.bind(this, selector, hash, uuid));
    return this;
  },

  /**
   * Sends out an event with the results of the `waitForElement` call
   *
   * @method _waitForElementCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @return {object} promise WaitForElement promise
   * @private
   */

  _waitForElementCb: function (selector, hash, uuid) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'waitForElement', selector: selector, uuid: uuid, hash: hash});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Returns the number of elements matched by the selector
   *
   * @method getNumberOfElements
   * @param {string} selector Selector expression to find the elements
   * @param {integer} expected Expected number of matched elements
   * @param {string} uuid Unique hash of that fn call
   * @chainable
   */

  getNumberOfElements: function (selector, expected, hash) {
    this.actionQueue.push(this.webdriverClient.elements.bind(this.webdriverClient, selector));
    this.actionQueue.push(this._getNumberOfElementsCb.bind(this, selector, hash, expected));
    return this;
  },

  /**
   * Sends out an event with the results of the `getNumberOfElements` call
   *
   * @method _getNumberOfElementsCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {integer} expected Expected number of matched elements
   * @param {string} res Serialized JSON with the results of the getNumberOfElements call
   * @return {object} promise GetNumberOfElements promise
   * @private
   */

  _getNumberOfElementsCb: function (selector, hash, expected, res) {
    var deferred = Q.defer();
    var result = JSON.parse(res);
    // check if the expression matched any element
    if (result.value === -1) {
      this.events.emit('driver:message', {key: 'numberOfElements', hash: hash, selector: selector, expected: expected, value: 0});
    } else {
      this.events.emit('driver:message', {key: 'numberOfElements', selector: selector, expected: expected, hash: hash, value: result.value.length});
    }

    deferred.resolve();
    return deferred.promise;
  },


  /**
   * Returns the number of visible elements matched by the selector
   *
   * @method getNumberOfVisibleElements
   * @param {string} selector Selector expression to find the elements
   * @param {integer} expected Expected number of matched elements
   * @param {string} uuid Unique hash of that fn call
   * @chainable
   */

  getNumberOfVisibleElements: function (selector, expected, hash) {
    this.actionQueue.push(this.webdriverClient.elements.bind(this.webdriverClient, selector));
    this.actionQueue.push(function (result) {
      var deferred = Q.defer();
      var res = JSON.parse(result);
      var resLength = res.value.length;
      var curLength = 0;
      var visibleElement = [];

      res.value.forEach(function (element) {
        var fakeResponse = JSON.stringify({sessionId: res.sessionId, status: 0, value: element.ELEMENT});
        this.webdriverClient.options.id = element.ELEMENT;
        this.webdriverClient.displayed.bind(this.webdriverClient, selector)(fakeResponse).then(function (visRes) {
          curLength++;
          if (JSON.parse(visRes).value === true) {
            visibleElement.push(element);
          }
          if (curLength === resLength) {
            deferred.resolve(JSON.stringify({sessionId: res.sessionId, status: 0, value: visibleElement}));
          }
        });
      }.bind(this));

      return deferred.promise;
    }.bind(this));
    this.actionQueue.push(this._getNumberOfVisibleElementsCb.bind(this, selector, hash, expected));
    return this;
  },

  /**
   * Sends out an event with the results of the `getNumberOfVisibleElements` call
   *
   * @method _getNumberOfElementsCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {integer} expected Expected number of matched elements
   * @param {string} res Serialized JSON with the results of the getNumberOfVisibleElements call
   * @return {object} promise GetNumberOfElements promise
   * @private
   */

  _getNumberOfVisibleElementsCb: function (selector, hash, expected, res) {
    var deferred = Q.defer();
    var result = JSON.parse(res);

    // check if the expression matched any element
    if (result.value === -1) {
      this.events.emit('driver:message', {key: 'numberOfVisibleElements', hash: hash, selector: selector, expected: expected, value: 0});
    } else {
      this.events.emit('driver:message', {key: 'numberOfVisibleElements', selector: selector, expected: expected, hash: hash, value: result.value.length});
    }

    deferred.resolve();
    return deferred.promise;
  }

};

/**
 * Mixes in element methods
 *
 * @param {Dalek.DriverNative} DalekNative Native driver base class
 * @return {Dalek.DriverNative} DalekNative Native driver base class
 */

module.exports = function (DalekNative) {
  // mixin methods
  Object.keys(Element).forEach(function (fn) {
    DalekNative.prototype[fn] = Element[fn];
  });

  return DalekNative;
};
