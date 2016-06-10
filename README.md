# trafficjam

A Ruby script that uses the excellent [Google Maps API][] to determine whether
your daily commute has unusual traffic and emails you if it does.

[Google Maps API]: https://developers.google.com/maps/

## Quick how-to

First, you need a free [Google Maps API key][] and access to an SMTP server.
[Mailgun][] has free plans.

[Google Maps API key]: https://developers.google.com/maps/documentation/javascript/get-api-key#key
[Mailgun]: https://www.mailgun.com/

Then, issue:

    $ git clone https://github.com/agorf/trafficjam.git
    $ cd trafficjam
    $ bundle install
    $ cp .env.sample .env
    $ $EDITOR .env
    $ bundle exec ruby trafficjam.rb

You can use [Cron][] to run the script at predetermined intervals (e.g. right
before heading out each morning).

[Cron]: https://en.wikipedia.org/wiki/Cron

## Configuration

The following environmental variables should be defined (all required unless
otherwise stated):

* `API_KEY` - Google Maps API key
* `MAX_DURATION_MINUTES` - If the estimated duration of your commute exceeds
  this value in minutes, you will get notified
* `ORIGIN` - Your home address
* `DESTINATION` - Your work address
* `AVOID` (optional) - Directs the Google Maps API to avoid `tolls`, `highways`,
  `ferries` or `indoor` steps
* `SMTP_HOST` - SMTP server hostname
* `SMTP_PORT` - SMTP server port
* `SMTP_USER` - SMTP server username
* `SMTP_PASS` - SMTP server password
* `RECIPIENT` - Email address to get notified at

## License

Licensed under the MIT license (see [LICENSE.txt][]).

[LICENSE.txt]: https://github.com/agorf/trafficjam/blob/master/LICENSE.txt

## Author

Angelos Orfanakos, http://agorf.gr/
