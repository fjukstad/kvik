dalek-browser-phantomjs
=======================

> DalekJS browser plugin for PhantomJS

[![Build Status](https://travis-ci.org/dalekjs/dalek-browser-phantomjs.png)](https://travis-ci.org/dalekjs/dalek-browser-phantomjs)
[![Build Status](https://drone.io/github.com/dalekjs/dalek-browser-phantomjs/status.png)](https://drone.io/github.com/dalekjs/dalek-browser-phantomjs/latest)
[![Dependency Status](https://david-dm.org/dalekjs/dalek-browser-phantomjs.png)](https://david-dm.org/dalekjs/dalek-browser-phantomjs)
[![devDependency Status](https://david-dm.org/dalekjs/dalek-browser-phantomjs/dev-status.png)](https://david-dm.org/dalekjs/dalek-browser-phantomjs#info=devDependencies)
[![NPM version](https://badge.fury.io/js/dalek-browser-phantomjs.png)](http://badge.fury.io/js/dalek-browser-phantomjs)
[![Coverage](http://dalekjs.com/package/dalek-browser-phantomjs/master/coverage/coverage.png)](http://dalekjs.com/package/dalek-browser-phantomjs/master/coverage/index.html)
[![unstable](https://rawgithub.com/hughsk/stability-badges/master/dist/unstable.svg)](http://github.com/hughsk/stability-badges)
[![Stories in Ready](https://badge.waffle.io/dalekjs/dalek-browser-phantomjs.png?label=ready)](https://waffle.io/dalekjs/dalek-browser-phantomjs)
[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/dalekjs/dalek-browser-phantomjs/trend.png)](https://bitdeli.com/free "Bitdeli Badge")
[![Built with Grunt](https://cdn.gruntjs.com/builtwith.png)](http://gruntjs.com/)

[![NPM](https://nodei.co/npm/dalek-browser-phantomjs.png)](https://nodei.co/npm/dalek-browser-phantomjs/)
[![NPM](https://nodei.co/npm-dl/dalek-browser-phantomjs.png)](https://nodei.co/npm/dalek-browser-phantomjs/)

## Ressources

[API Docs](http://dalekjs.com/package/dalek-browser-phantomjs/master/api/index.html) -
[Trello](https://trello.com/b/UjcpWj7v/dalek-browser-phantomjs) -
[Code coverage](http://dalekjs.com/package/dalek-browser-phantomjs/master/coverage/index.html) -
[Code complexity](http://dalekjs.com/package/dalek-browser-phantomjs/master/complexity/index.html) -
[Contributing](https://github.com/dalekjs/dalek-browser-phantomjs/blob/master/CONTRIBUTING.md) -
[User Docs](http://dalekjs.com/docs/phantomjs.html) -
[Homepage](http://dalekjs.com) -
[Twitter](http://twitter.com/dalekjs)

## Docs
This module is a browser plugin for [DalekJS](//github.com/dalekjs/dalek).
It provides a browser launcher as well the PhantomJS browser itself.

The browser plugin comes bundled with the DalekJS base framework.

You can use the browser plugin beside others (it is the default)
by adding a config option to the your Dalekfile:

```js
"browsers": ["phantomjs", "chrome"]
```

Or you can tell Dalek that it should test in this & another browser via the command line:

```
$ dalek mytest.js -b phantomjs,chrome
```

## Help Is Just A Click Away

### #dalekjs on FreeNode.net IRC

Join the `#daleksjs` channel on [FreeNode.net](http://freenode.net) to ask questions and get help.

### [Google Group Mailing List](https://groups.google.com/forum/#!forum/dalekjs)

Get announcements for new releases, share your projects and ideas that are
using DalekJS, and join in open-ended discussion that does not fit in
to the Github issues list or StackOverflow Q&A.

**For help with syntax, specific questions on how to implement a feature
using DalekJS, and other Q&A items, use StackOverflow.**

### [StackOverflow](http://stackoverflow.com/questions/tagged/dalekjs)

Ask questions about using DalekJS in specific scenarios, with
specific features. For example, help with syntax, understanding how a feature works and
how to override that feature, browser specific problems and so on.

Questions on StackOverflow often turn in to blog posts or issues.

### [Github Issues](//github.com/dalekjs/dalek-browser-phantomjs/issues)

Report issues with DalekJS, submit pull requests to fix problems, or to
create summarized and documented feature requests (preferably with pull
requests that implement the feature).

**Please don't ask questions or seek help in the issues list.** There are
other, better channels for seeking assistance, like StackOverflow and the
Google Groups mailing list.

![DalekJS](https://raw.github.com/dalekjs/dalekjs.com/master/img/logo.png)

## Legal FooBar (MIT License)

Copyright (c) 2013 Sebastian Golasch

Distributed under [MIT license](https://github.com/dalekjs/dalek-browser-phantomjs/blob/master/LICENSE-MIT)