/* jshint strict: false, unused: false */
/* global __SCRIPT */
var args = Array.prototype.slice.call(arguments, 0);
var dalek = {
	store: args.shift(),
	test: [],
	assert: {},
	_assert: {
		ok: function (isOk, message) {
			this.test.push({ok: isOk, message: message});
		}
	},
	data: function (key, value) {
		if(value) {
			this.store[key] = value;
			return this;
		}
		return this.store[key] || null;
	}
};

dalek.assert.ok = function () {
	return dalek._assert.ok.apply(dalek, Array.prototype.slice.call(arguments));
};

var userRet = function(__USERARGUMENTS) {
	__SCRIPT;
}.apply(dalek, args);

return {dalek: dalek.store, test: dalek.test, userRet: userRet};