LaunchDarkly Redis daemon
=========================

What is it?
-----------

The LaunchDarkly Redis daemon establishes a connection to the LaunchDarkly streaming API, then pushes feature updates to a Redis store.

The daemon can be used to offload the task of maintaining a stream and writing to Redis from our SDKs. This can give platforms that do not support SSE (e.g. PHP) the benefits of LaunchDarkly's streaming model.

The daemon can be configured to synchronize multiple environments, even across multiple projects.

Quick setup
-----------

1. Copy `ld-daemon.conf` to `/etc/ld-daemon.conf`, and edit to specify your Redis host and port, as well as Redis key prefixes and LaunchDarkly API keys for each environment you wish to synchronize.

2. If building from source, have `go` 1.4+ and `godep` installed, and run `godep go build`.

3. Run `ld-daemon`.

4. Set `stream` and `use_ldd` to `true` in your application's LaunchDarkly SDK configuration. Also ensure that you specify a Redis store in your configuration. This will turn off your SDK's streaming connection, but read feature flags from the Redis store. 

Configuration file format 
-------------------------

You can configure LDD to synchronize as many environments as you want, even across different projects. Each environment should define its own unique key prefix string in Redis to avoid conflicts. We recommend using a colon-separated string including your project key and environment key. Note that if you follow this convention, your SDKs *must* be configured to use that same prefix-- so you'll need to parameterize both your API key and LDD key prefix in your SDK configuration.

Here's an example configuration file that synchronizes four environments across two different projects (called Spree and Shopnify):

        [redis]
        host = "localhost"
        port = 6379

        [main]
        streamUri = "https://stream.launchdarkly.com"
        baseUri = "https://app.launchdarkly.com"
        exitOnError = false

        [environment "Spree Project Production"]
        prefix = "ld:spree:production"
        apiKey = "SPREE_PROD_API_KEY"

        [environment "Spree Project Test"]
        prefix = "ld:spree:test"
        apiKey = "SPREE_TEST_API_KEY"

        [environment "Shopnify Project Production"]
        prefix = "ld:shopnify:production"
        apiKey = "SHOPNIFY_PROD_API_KEY"

        [environment "Shopnify Project Test"]
        prefix = "ld:shopnify:test"
        apiKey = "SHOPNIFY_TEST_API_KEY"
