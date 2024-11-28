# Groat
[![codecov](https://codecov.io/gh/godepo/groat/graph/badge.svg?token=6V8BXLVZKX)](https://codecov.io/gh/godepo/groat)

Groat is lightweight toolkit for creation articulated declarative tests 

This small library help write accurate and elegant tests on golang in BDD and FP conceptions.

Groat present small and effective test case implementation, that support:

* Test Data and Dependency states.
* xUnit like structure - with support before each, before all and after each and after all phases.
* Support DSL for AAA(arrange, act, assert)/GWT (given, when, then) pattern.
* You can use any assertions libraries with groat.
* Native integrate with golang testing toolkit.
* Can make you testing experience light, opposite declarative and imperative testing.
* Have little toolkit for running blazing fast test e2e and i9n based on test containers.
* Can run all test in parallel, opposite testify/mocks.
