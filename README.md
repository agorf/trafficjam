# trafficjam

A Ruby script that uses the [Google Maps API][] to notify you if your daily
commute has unusual traffic.

[Google Maps API]: https://developers.google.com/maps/

## Quick how-to

First, you need a free [Google Maps API key][] and access to an SMTP server.
[Mailgun][] has free plans.

[Google Maps API key]: https://developers.google.com/maps/documentation/javascript/get-api-key#key
[Mailgun]: https://www.mailgun.com/

Then, issue once:

    $ git clone https://github.com/agorf/trafficjam.git
    $ cd trafficjam
    $ bundle install
    $ cp .env.sample .env

Configure the script with your `$EDITOR`:

    $ $EDITOR .env

Run the script:

    $ bundle exec ruby trafficjam.rb

You can use [Cron][] to run the script at predetermined intervals (e.g. right
before heading out each morning).

[Cron]: https://en.wikipedia.org/wiki/Cron

## Configuration

The following environment variables should be defined (all required unless
otherwise stated):

* `ORIGINS` &ndash; Your home address
* `DESTINATIONS` &ndash; Your work address
* `KEY` &ndash; Your Google Maps API key
* `MODE` (optional) &ndash; See [here](https://developers.google.com/maps/documentation/distance-matrix/intro#travel_modes)
* `AVOID` (optional) &ndash; See [here](https://developers.google.com/maps/documentation/distance-matrix/intro#Restrictions)
* `TRAFFIC_MODEL` (optional) &ndash; See [here](https://developers.google.com/maps/documentation/distance-matrix/intro#traffic-model)
* `MAX_DURATION_MINUTES` &ndash; If the estimated duration of your commute exceeds
  this value in minutes, you will get notified
* `SMTP_HOST` &ndash; SMTP server hostname
* `SMTP_PORT` &ndash; SMTP server port
* `SMTP_USER` &ndash; SMTP server username
* `SMTP_PASS` &ndash; SMTP server password
* `RECIPIENT` &ndash; Email address to get notified at

## License

Licensed under the MIT license (see [LICENSE.txt][]).

[LICENSE.txt]: https://github.com/agorf/trafficjam/blob/master/LICENSE.txt

## Author

Angelos Orfanakos, http://agorf.gr/
