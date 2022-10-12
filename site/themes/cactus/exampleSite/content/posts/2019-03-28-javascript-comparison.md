---
title: JavaScript Comparison operation at a glance
date: 2019-03-28 09:00:00
tags:
    - JavaScript
category: tech
keywords:
    - Javascript
    - ES2015
    - ES6
---

When given a scenario like:

```javascript
console.log(null > -1) //true
```

It produces `true`, which makes me think `null` is treated as `0`. But when I run:
```javascript
console.log(null == 0) // false
console.log(null > 0) // false
console.log(null < 0) // false
```
They all output `false`!

I googled a lot and finally found answers in [Ecma-262 Specification](http://www.ecma-international.org/ecma-262/8.0/#sec-abstract-equality-comparison).

The comparison `x == y`, where x and y are values, produces true or false. Such a comparison is performed as follows:
<!-- more -->
```text
1. If TypeFromProto(x) is the same as TypeFromProto(y), then return the result of performing Strict Equality Comparison x === y.
2. If x is null and y is undefined, return true.
3. If x is undefined and y is null, return true.
4. If TypeFromProto(x) is Number and TypeFromProto(y) is String, return the result of the comparison x == ToNumber(y).
5. If TypeFromProto(x) is String and TypeFromProto(y) is Number, return the result of the comparison ToNumber(x) == y.
6. If TypeFromProto(x) is Boolean, return the result of the comparison ToNumber(x) == y.
7. If TypeFromProto(y) is Boolean, return the result of the comparison x == ToNumber(y).
8. If TypeFromProto(x) is either String, Number, or Symbol and TypeFromProto(y) is Object, return the result of the comparison x == ToPrimitive(y).
9. If TypeFromProto(x) is Object and TypeFromProto(y) is either String, Number, or Symbol, return the result of the comparison ToPrimitive(x) == y.
10. Return false.
```

Relational comparison is much more complex so I'm not copying that section. Read at the [spec website](http://www.ecma-international.org/ecma-262/8.0/#sec-abstract-relational-comparison).

## TL;DR

Anyway it seems that in `null == 0`, `null` is treated just as is, and **equality** comparison between `null` and `Number` always return **false** (No 10).
But when it comes `null > -1`, `null` is conversed to 0 using [ToNumber()](http://www.ecma-international.org/ecma-262/8.0/#sec-tonumber) algorithm.

Read more:
* [https://github.com/getify/You-Dont-Know-JS/issues/1238](https://github.com/getify/You-Dont-Know-JS/issues/1238)
* [http://www.ecma-international.org/ecma-262/8.0/#sec-abstract-relational-comparison](http://www.ecma-international.org/ecma-262/8.0/#sec-abstract-relational-comparison)
* [http://www.ecma-international.org/ecma-262/8.0/#sec-tonumber](http://www.ecma-international.org/ecma-262/8.0/#sec-tonumber)
