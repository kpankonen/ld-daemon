LaunchDarkly Redis daemon
=========================

What is it?
-----------

The LaunchDarkly Redis daemon establishes a connection to the LaunchDarkly streaming API, then pushes feature updates to a Redis store.

The daemon can be used to offload the task of maintaining a stream and writing to Redis from our SDKs. This can give platforms that do not support SSE (e.g. PHP) the benefits of LaunchDarkly's streaming model.

Quick setup
-----------

1. Edit `ld-daemon.conf` to specify your Redis host and port, key prefix, and LaunchDarkly API key.

2. If building from source, have `go` 1.4+ and `godep` installed, and run `godep go build`.

3. Run `ld-daemon`.

4. Set `stream` and `use_ldd` to `true` in your application's LaunchDarkly SDK configuration. Also ensure that you specify a Redis store in your configuration. This will turn off your SDK's streaming connection, but read feature flags from the Redis store. 
