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
    address:   ENV.fetch('SMTP_HOST'),
    port:      ENV.fetch('SMTP_PORT'),
    user_name: ENV.fetch('SMTP_USER'),
    password:  ENV.fetch('SMTP_PASS')
end

def send_email(subject, body = '')
  Mail.deliver do
    to      ENV.fetch('RECIPIENT')
    from    ENV.fetch('RECIPIENT')
    subject subject
    body    body
  end
end

params = {
  origins:        ENV.fetch('ORIGINS'),
  destinations:   ENV.fetch('DESTINATIONS'),
  key:            ENV.fetch('KEY'),
  mode:           ENV['MODE'],
  avoid:          ENV['AVOID'],
  departure_time: 'now',
  traffic_model:  ENV['TRAFFIC_MODEL'],
}

query_str = params.map {|k, v|
  "#{CGI.escape(k.to_s)}=#{CGI.escape(v.to_s)}" unless v.nil?
}.compact.join('&')

url = API_URL + '?' + query_str

begin
  resp = JSON.parse(open(url).read)
  duration = resp['rows'][0]['elements'][0]['duration_in_traffic']['value']
rescue => e
  send_email(e.message, e.backtrace.join("\n"))
end

duration_in_minutes = (duration / 60.0).ceil

if duration_in_minutes > ENV.fetch('MAX_DURATION_MINUTES').to_f
  send_email("#{duration_in_minutes} minutes")
end
