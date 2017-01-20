# trafficjam

Queries the [Google Maps API][] and notifies you if a trip has unusual traffic.
Useful for knowing about traffic in your daily commute before getting stuck.

[Google Maps API]: https://developers.google.com/maps/

This [used to be written in Ruby][rewrite-commit].

[rewrite-commit]: https://github.com/agorf/trafficjam/commit/e8434f742b6dec8e1644a533f769540aeebcedb2

## Quick how-to

First, you need a free [Google Maps API key][] and access to an SMTP server.
[Mailgun][] has a free plan.

[Google Maps API key]: https://developers.google.com/maps/documentation/javascript/get-api-key#key
[Mailgun]: https://www.mailgun.com/

Then install [Go](https://golang.org/), if you haven't, and issue once:

    $ git clone https://github.com/agorf/trafficjam.git
    $ cd trafficjam
    $ go build trafficjam.go

Alternatively, you can issue:

    $ go get github.com/agorf/trafficjam

Configure the program with your `$EDITOR`:

    $ cp config.json.sample config.json
    $ $EDITOR config.json

Run the program:

    $ ./trafficjam config.json

Or if you used `go get`:

    $ trafficjam config.json

You can use [Cron][] to run the script at predetermined intervals (e.g. right
before heading out each morning). Here's what I have:

    0,5,10,15,20,25,30 9 * * 1-5 agorf ./trafficjam config.json

This runs the program on working days, from 9:00 until 9:30, every five minutes.

To avoid getting spammed, you need to figure out the right threshold
(`max_duration` config option) that you want to get notified for.

[Cron]: https://en.wikipedia.org/wiki/Cron

## Configuration

The configuration file is a plain JSON file. The following keys can be defined
(all required unless otherwise stated):

* `origins` &ndash; Your home address
* `destinations` &ndash; Your work address
* `api_key` &ndash; Your Google Maps API key
* `mode` (optional) &ndash; See [here](https://developers.google.com/maps/documentation/distance-matrix/intro#travel_modes)
* `avoid` (optional) &ndash; See [here](https://developers.google.com/maps/documentation/distance-matrix/intro#Restrictions)
* `traffic_model` (optional) &ndash; See [here](https://developers.google.com/maps/documentation/distance-matrix/intro#traffic-model)
* `max_duration` (minutes) &ndash; If the estimated duration of your commute exceeds
  this value in minutes, you will get notified
* `host` (under `smtp`) &ndash; SMTP server hostname
* `port` (under `smtp`) &ndash; SMTP server port
* `user` (under `smtp`) &ndash; SMTP server username
* `pass` (under `smtp`) &ndash; SMTP server password
* `recipient` &ndash; Email address to get notified at

See `config.json.sample` for a sample configuration file.

## License

Licensed under the MIT license (see [LICENSE.txt][]).

[LICENSE.txt]: https://github.com/agorf/trafficjam/blob/master/LICENSE.txt

## Author

Angelos Orfanakos, http://agorf.gr/
