'use strict';

var expect = require('chai').expect;
var Screenshot = require('../lib/commands/screenshot.js');

describe('dalek-driver-native Commands Screenshot', function() {
  var MockDate = function(config) {
    this.getHours = function() {
      return config.hours;
    };
    this.getMinutes = function() {
      return config.minutes;
    };
    this.getSeconds = function() {
      return config.seconds;
    };
    this.getTimezoneOffset = function() {
      return config.timezoneOffset;
    };
  };

  var screenshot;
  beforeEach(function(){
    screenshot = new Screenshot();
  });

  describe('_parseDate', function() {
    it('should take a js date object (2014,11,30) and return the ISO 8601 formatted date', function() {
      var date = screenshot._parseDate(new Date(2014, 11, 30));
      expect(date).to.equal('2014-12-30');
    });

    it('should fill up the month to two digits (2015,2,30)', function() {
      var date = screenshot._parseDate(new Date(2015, 2, 30));
      expect(date).to.equal('2015-03-30');
    });

    it('should fill up the day to two digits (2020,11,1)', function() {
      var date = screenshot._parseDate(new Date(2020, 11, 1));
      expect(date).to.equal('2020-12-01');
    });

    it('should use the current date if no parameter has been given', function() {
      var date = screenshot._parseDate();
      expect(date).to.match(/[\d]{4}-[\d]{2}-[\d]{2}/);
    });
  });

  describe('_parseDateTime', function() {
    it('should return the ISO 8601 formatted date time for non UTC dates with negative offset hours (YYYY-MM-DDThh:mm:ss-hh:mm)', function() {
      var dateTime = screenshot._parseDatetime(new MockDate({hours: 10, minutes: 20, seconds: 30, timezoneOffset: -60}));
      expect(dateTime).to.match(/^[\d]{4}-[\d]{2}-[\d]{2}T{1}[\d]{2}:[\d]{2}:[\d]{2}(-|\+)[\d]{2}:[\d]{2}$/);
    });
    it('should return the ISO 8601 formatted date time for non UTC dates with negative offset minutes (YYYY-MM-DDThh:mm:ss-hh:mm)', function() {
      var dateTime = screenshot._parseDatetime(new MockDate({hours: 10, minutes: 20, seconds: 30, timezoneOffset: -30}));
      expect(dateTime).to.match(/^[\d]{4}-[\d]{2}-[\d]{2}T{1}[\d]{2}:[\d]{2}:[\d]{2}(-|\+)[\d]{2}:[\d]{2}$/);
    });
    it('should return the ISO 8601 formatted date time for non UTC dates with negative offset hours and minutes (YYYY-MM-DDThh:mm:ss-hh:mm)', function() {
      var dateTime = screenshot._parseDatetime(new MockDate({hours: 10, minutes: 20, seconds: 30, timezoneOffset: -90}));
      expect(dateTime).to.match(/^[\d]{4}-[\d]{2}-[\d]{2}T{1}[\d]{2}:[\d]{2}:[\d]{2}(-|\+)[\d]{2}:[\d]{2}$/);
    });
    it('should return the ISO 8601 formatted date time for non UTC dates with positive offset hours (YYYY-MM-DDThh:mm:ss+hh:mm)', function() {
      var dateTime = screenshot._parseDatetime(new MockDate({hours: 10, minutes: 20, seconds: 30, timezoneOffset: +60}));
      expect(dateTime).to.match(/^[\d]{4}-[\d]{2}-[\d]{2}T{1}[\d]{2}:[\d]{2}:[\d]{2}(-|\+)[\d]{2}:[\d]{2}$/);
    });
    it('should return the ISO 8601 formatted date time for non UTC dates with positive offset minutes (YYYY-MM-DDThh:mm:ss+hh:mm)', function() {
      var dateTime = screenshot._parseDatetime(new MockDate({hours: 10, minutes: 20, seconds: 30, timezoneOffset: +20}));
      expect(dateTime).to.match(/^[\d]{4}-[\d]{2}-[\d]{2}T{1}[\d]{2}:[\d]{2}:[\d]{2}(-|\+)[\d]{2}:[\d]{2}$/);
    });
    it('should return the ISO 8601 formatted date time for non UTC dates with positive offset hours and minutes (YYYY-MM-DDThh:mm:ss-hh:mm)', function() {
      var dateTime = screenshot._parseDatetime(new MockDate({hours: 10, minutes: 20, seconds: 30, timezoneOffset: 90}));
      expect(dateTime).to.match(/^[\d]{4}-[\d]{2}-[\d]{2}T{1}[\d]{2}:[\d]{2}:[\d]{2}(-|\+)[\d]{2}:[\d]{2}$/);
    });
    it('should return the ISO 8601 formatted date time for UTC dates (YYYY-MM-DDThh:mm:ssZ)', function() {
      var dateTime = screenshot._parseDatetime(new MockDate({hours: 10, minutes: 20, seconds: 30, timezoneOffset: 0}));
      expect(dateTime).to.match(/^[\d]{4}-[\d]{2}-[\d]{2}T{1}[\d]{2}:[\d]{2}:[\d]{2}Z$/);
    });
    // nyone in UTC +/-0 tell me if this works for you! :)
    it('should use the current date if no parameter has been given and return the ISO 8601 formatted date time', function() {
      var dateTime = screenshot._parseDatetime();
      expect(dateTime).to.match(/[\d]{4}-[\d]{2}-[\d]{2}T{1}[\d]{2}:[\d]{2}:[\d]{2}(Z|(-|\+)[\d]{2}:[\d]{2})/);
    });
  });
});