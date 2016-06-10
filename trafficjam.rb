#!/usr/bin/env ruby

require 'cgi'
require 'dotenv'
require 'json'
require 'mail'
require 'open-uri'

API_URL = 'https://maps.googleapis.com/maps/api/distancematrix/json'

Dotenv.load

Mail.defaults do
  delivery_method :smtp,
    address:   ENV['SMTP_HOST'],
    port:      ENV['SMTP_PORT'],
    user_name: ENV['SMTP_USER'],
    password:  ENV['SMTP_PASS']
end

def send_email(subject:, body: '')
  Mail.deliver do
    to      ENV['RECIPIENT']
    from    ENV['RECIPIENT']
    subject subject
    body    body
  end
end

params = {
  origins:        ENV['ORIGIN'],
  destinations:   ENV['DESTINATION'],
  key:            ENV['API_KEY'],
  departure_time: 'now',
}

params[:avoid] = ENV['AVOID'] if ENV.has_key?('AVOID')

query_str = params.map {|k, v|
  "#{CGI.escape(k.to_s)}=#{CGI.escape(v.to_s)}"
}.join('&')

url = API_URL + '?' + query_str

begin
  resp = JSON.parse(open(url).read)
  duration = resp['rows'][0]['elements'][0]['duration_in_traffic']['value']
rescue => e
  send_email(subject: e.message, body: e.backtrace.join("\n"))
end

duration_in_minutes = (duration / 60.0).ceil

if duration_in_minutes > ENV['MAX_DURATION_MINUTES'].to_f
  send_email(subject: "#{duration_in_minutes} minutes")
end
