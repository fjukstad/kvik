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
var http = require('http');
var Q = require('q');

// export the driver base function
module.exports = function exportDriverBase(WebDriver) {

  /**
   * Native Webdriver base class
   *
   * @module DalekJS
   * @class Driver
   * @namespace Internal
   */

  var Driver = {

    /**
     * Parses an JSON Wire protocol dummy url
     *
     * @method parseUrl
     * @param {string} url URL with placeholders
     * @param {object} options List of url options
     * @return {string} url Parsed URL
     */

    parseUrl: function (url, options) {
      return Object.keys(options).reduce(this._replacePlaceholderInUrl.bind(this, options), url);
    },

    /**
     * Replaces placeholders in urls
     *
     * @method _replacePlaceholderInUrl
     * @param {object} options List of url options
     * @param {string} url URL with placeholders
     * @param {string} option Option to process
     * @return {string} url Parsed URL
     * @private
     */

    _replacePlaceholderInUrl: function (options, url, option) {
      return url.replace(':' + option, options[option]);
    },

    /**
     * Generates a set of params for the message body of the request
     *
     * @method generateParamset
     * @param {object|null} requestedParams Keys & placeholders for the paramset
     * @param {object} providedParams Values for the paramset
     * @return {object} params Params for the message body
     */

    generateParamset: function (requestedParams, providedParams) {
      return !requestedParams ? {} : requestedParams.reduce(this._mapParams.bind(this, providedParams), {});
    },

    /**
     * Mpas object values & keys of two objects
     *
     * @method _mapParams
     * @param {object} providedParams Values for the paramset
     * @param {object} params The object to be filled
     * @param {string} param The key of the output object
     * @param {integer} idx Index of the iteration
     * @return {object} params Params for the message body
     * @private
     */

    _mapParams: function (providedParams, params, param, idx) {
      params[param] = providedParams[idx];
      return params;
    },

    /**
     * Generates the message body for webdriver client requests of type POST
     *
     * @method generateBody
     * @param {object} options Browser options (name, bin path, etc.)
     * @param {function|undefined} cb Callback function that should be invoked to generate the message body
     * @param {Dalek.Internal.Webdriver} wd Webdriver base object
     * @param {object} params Parameters that should be part of the message body
     * @return {string} body Serialized JSON of body request data
     */

    generateBody: function (options, cb, wd, params) {
      // if no cb is given, generate a body with the `desiredCapabilities` object
      if (!cb) {
        // check if we have parameters set up
        if (Object.keys(params).length > 0) {
          return JSON.stringify(params);
        }
        return '';
      }

      // invoke the given callback & stringify
      var data = cb.call(wd, params);
      return data === null ? '{}' : JSON.stringify(data);
    },

    /**
     * Generates the request options for a webdriver client request
     *
     * @method generateRequestOptions
     * @param {string} hostname Hostname of the webdriver server
     * @param {integer} port Port of the webdriver server
     * @param {string} prefix Url address prefix of the webdriver endpoint
     * @param {string} url Url of the webdriver method
     * @param {string} method Request method e.g. (GET, POST, DELETE, PUT)
     * @param {string} body The message body of the request
     * @return {object} options Request options
     */

    generateRequestOptions: function (hostname, port, prefix, url, method, body, auth) {
      var options = {
        hostname: hostname,
        port: port,
        path: prefix + url,
        method: method,
        headers: {
          'Content-Type': 'application/json;charset=utf-8',
          'Content-Length': Buffer.byteLength(body, 'utf8')
        }
      };

      // check if auth information is available
      if (auth) {
        options.auth = auth;
      }

      return options;
    },

    /**
     * Generates a new webdriver client command
     * Takes a skeleton of obtions that will be converted
     * into a new function that can be invoked & will issue
     * a webdriver command to the webdriver server
     *
     * @method addCommand
     * @param {object} remote Object skeleton that will be turned into a webdriver client method
     * @chainable
     */

    addCommand: function (remote) {
      // assign the generated function to the webdriver prototype
      WebDriver.prototype[remote.name] = this._generateWebdriverCommand(remote, this);
      return this;
    },

    /**
     * Generates the webdriver callback function
     *
     * @method _generateWebdriverCommand
     * @param {object} remote Dummy request body (function name, url, method)
     * @param {DalekJs.Internal.Driver} driver Driver instance
     * @return {function} webdriverCommand Generated webdriver command function
     * @private
     */

    _generateWebdriverCommand: function (remote, driver) {
      return function webdriverCommand() {
        var deferred = Q.defer();
        // the request meta data
        var params = Driver.generateParamset(remote.params, arguments);
        var body = Driver.generateBody({}, remote.onRequest, this, params);
        var options = Driver.generateRequestOptions(this.opts.host, this.opts.port, this.opts.path, Driver.parseUrl(remote.url, this.options), remote.method, body, this.opts.auth);

        // generate the request, wait for response & fire the request
        var req = new http.ClientRequest(options);
        req.on('response', driver._onResponse.bind(this, driver, remote, options, deferred));
        req.end(body);

        return deferred.promise;
      };
    },

    /**
     * Response callback function
     *
     * @method _onResponse
     * @param {DalekJs.Internal.Driver} driver Driver instance
     * @param {object} remote Dummy request body (function name, url, method)
     * @param {object} options Request options (method, port, path, headers, etc.)
     * @param {object} deferred Webdriver command deferred
     * @param {object} response Response from the webdriver server
     * @chainable
     * @private
     */

    _onResponse: function (driver, remote, options, deferred, response) {
      this.data = '';
      response.on('data', driver._concatDataChunks.bind(this));
      response.on('end', driver._onResponseEnd.bind(this, driver, response, remote, options, deferred));
      return this;
    },

    /**
     * Concatenates chunks of strings
     *
     * @method _concatDataChunks
     * @param {string} chunk String to add
     * @return {string} data Concatenated string
     * @private
     */

    _concatDataChunks: function (chunk) {
      return this.data += chunk;
    },

    /**
     * Response end callback function
     *
     * @method _onResponseEnd
     * @param {DalekJs.Internal.Driver} driver Driver instance
     * @param {object} response Response from the webdriver server
     * @param {object} remote Dummy request body (function name, url, method)
     * @param {object} options Request options (method, port, path, headers, etc.)
     * @param {object} deferred Webdriver command deferred
     * @chainable
     * @priavte
     */

    _onResponseEnd: function (driver, response, remote, options, deferred) {
      return driver[(response.statusCode === 500 ? '_onError' : '_onSuccess')].bind(this)(driver, response, remote, options, deferred);
    },

    /**
     * On error callback function
     *
     * @method _onError
     * @param {DalekJs.Internal.Driver} driver Driver instance
     * @param {object} response Response from the webdriver server
     * @param {object} remote Dummy request body (function name, url, method)
     * @param {object} options Request options (method, port, path, headers, etc.)
     * @param {object} deferred Webdriver command deferred
     * @chainable
     * @private
     */

    _onError: function (driver, response, remote, options, deferred) {
      // Provide a default error handler to prevent hangs.
      if (!remote.onError) {
        remote.onError = function (request, remote, options, deferred, data) {
          data = JSON.parse(data);
          var value = -1;
          if (typeof data.value.message === 'string') {
            var msg = JSON.parse(data.value.message);
            value = msg.errorMessage;
          }
          deferred.resolve(JSON.stringify({'sessionId': data.sessionId, value: value}));
        };
      }
      remote.onError.call(this, response, remote, options, deferred, this.data);
      return this;
    },

    /**
     * On success callback function
     *
     * @method _onSuccess
     * @param {DalekJs.Internal.Driver} driver Driver instance
     * @param {object} response Response from the webdriver server
     * @param {object} remote Dummy request body (function name, url, method)
     * @param {object} options Request options (method, port, path, headers, etc.)
     * @param {object} deferred Webdriver command deferred
     * @chainable
     * @private
     */

    _onSuccess: function (driver, response, remote, options, deferred) {
      // log response data
      this.events.emit('driver:webdriver:response', {
        statusCode: response.statusCode,
        method: response.req.method,
        path: response.req.path,
        data: this.data
      });
      
      if (remote.onResponse) {
        remote.onResponse.call(this, response, remote, options, deferred, this.data);
      } else {
        deferred.resolve(this.data);
      }
      return this;
    }
  };

  return Driver;
};
