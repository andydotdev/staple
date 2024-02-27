Staple
======

[![Go Reference](https://pkg.go.dev/badge/andy.dev/staple.svg)](https://pkg.go.dev/andy.dev/staple)

    import "andy.dev/staple"

Staple provides Erlang-ish supervisor trees for Go.

It is intended to deal gracefully with the real failure cases that can
occur with supervision trees (such as burning all your CPU time endlessly
restarting dead services), while also making no unnecessary demands on the
"service" code, and providing hooks to perform adequate logging with in a
production environment.

[A blog post describing the design decisions](http://www.jerf.org/iri/post/2930)
is available.

This module is fairly fully covered
with [godoc](https://pkg.go.dev/andy.dev/staple)
including an example, usage, and everything else you might expect from a
README.md on GitHub. (DRY.)

Special Thanks
--------------

Special thanks to [thejerf](https://github.com/thejerf) for a great package.
When in doubt, use that instead, as this package primariliy exists so I can
tweak it ever so slightly to suit my needs.


Code Signing
------------

Starting with the commit after ac7cf8591b, I will be signing this repository
with the ["jerf" keybase account](https://keybase.io/jerf). If you are viewing
this repository through GitHub, you should see the commits as showing as
"verified" in the commit view.

(Bear in mind that due to the nature of how git commit signing works, there
may be runs of unverified commits; what matters is that the top one is signed.)

Aspiration
----------

One of the big wins the Erlang community has with their pervasive OTP
support is that it makes it easy for them to distribute libraries that
easily fit into the OTP paradigm. It ought to someday be considered a good
idea to distribute libraries that provide some sort of supervisor tree
functionality out of the box. It is possible to provide this functionality
without explicitly depending on the Staple library.

Changelog
---------

staple uses semantic versioning and go modules.

* 1.0.0:
  * Initial fork from [suture](https://github.com/thejerf/suture)
  * Minor lint-based cleanup, nothing to write home about.
  * Add [slog](https://pkg.go.dev/log/slog) as the default logger and allow all
    event types to conform to `slog.LogValuer` so that they decide their own
    logging levels.
* suture-4.0.2:
  * Add the ability to specify a handler for non-string panics to format
    them.
  * Fixed an issue where trying to close a currently-panicked service was
    having problems. (This may have leaked goroutines in other ways too.)
  * Merged a PR that addresses race conditions in the test suite. (These
    seem to have been isolated to the test suite and not have affected the
    core code.)
* suture-4.0.1:
  * Add a channel returned from ServeBackground that can be used to
    examine any error coming out of the supervisor once it is stopped.
  * Tweak up the docs to try to make it more clear suture's special
    error returns are checked via errors.Is when possible, addressing
    issue #51.
* suture-4.0:
  * Switched the entire API to be context based.
  * Switched how logging works to take a single closure that will be
    presented with a defined set of structs, rather than a set of closures
    for each event.
  * Consequently, "Stop" removed from the Service interface. A wrapper for
    old-style code is provided.
  * Services can now return errors. Errors will be included in the log
    message. Two special errors control restarting behavior:
      * ErrDoNotRestart indicates the service should not be restarted,
        but other services should be unaffected.
      * ErrTerminateTree indicates the parent service tree should be
        terminated. Supervisor trees can be configured to either continue
        terminating upwards, or terminate themselves but not continue
        propagating the termination upwards.
  * UnstoppedServiceReport calling semantics modified to allow correctly
    retrieving reports from entire trees. (Prior to 4.0, a report was
    only on the supervisor it was called on.)