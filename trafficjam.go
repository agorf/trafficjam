package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"os/user"
)

const (
	name   = "trafficjam"
	apiURL = "https://maps.googleapis.com/maps/api/distancematrix/json"
)

type config struct {
	origins      string `json:"origins"`
	destinations string `json:"destinations"`
	apiKey       string `json:"api_key"`
	mode         string `json:"mode"`
	avoid        string `json:"avoid"`
	trafficModel string `json:"traffic_model"`
	maxDuration  int    `json:"max_duration"`
	smtp         struct {
		host string `json:"host"`
		port int    `json:"port"`
		user string `json:"user"`
		pass string `json:"pass"`
	} `json:"smtp"`
	recipient string `json:"recipient"`
}

type apiResponse struct {
	rows []struct {
		elements []struct {
			durationInTraffic struct {
				text  string `json:"text"`
				value int    `json:"value"`
			} `json:"duration_in_traffic"`
			status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	status string `json:"status"`
}

func main() {
	log.SetPrefix(name + ": ")
	log.SetFlags(0)

	if len(os.Args) != 2 {
		log.Fatalf("usage: %s config.json", name)
	}

	conf, err := readConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	params := map[string]string{
		"origins":        conf.origins,
		"destinations":   conf.destinations,
		"key":            conf.apiKey,
		"mode":           conf.mode,
		"avoid":          conf.avoid,
		"departure_time": "now",
		"traffic_model":  conf.trafficModel,
	}
	apiResp, err := queryMapsAPI(params)
	if err != nil {
		log.Fatal(err)
	}

	duration := apiResp.rows[0].elements[0].durationInTraffic.value
	if duration > conf.maxDuration*60 {
		if err := sendMail(conf, apiResp.rows[0].elements[0].durationInTraffic.text); err != nil {
			log.Fatal(err)
		}
	}
}

func readConfig(filename string) (*config, error) {
	var conf config

	confData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(confData, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func queryMapsAPI(params map[string]string) (*apiResponse, error) {
	query := url.Values{}
	for key, val := range params {
		if val != "" {
			query.Set(key, val)
		}
	}

	uri, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	uri.RawQuery = query.Encode()

	resp, err := http.Get(uri.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}
	if apiResp.status != "OK" {
		return nil, fmt.Errorf("%s: bad response status: %s\n", name, apiResp.status)
	}
	if len(apiResp.rows) != 1 {
		return nil, fmt.Errorf("%s: response row count is not 1\n", name)
	}
	if len(apiResp.rows[0].elements) != 1 {
		return nil, fmt.Errorf("%s: response first row element count is not 1\n", name)
	}
	if apiResp.rows[0].elements[0].status != "OK" {
		return nil, fmt.Errorf("%s: bad response first row first element status: %s\n", name, apiResp.rows[0].elements[0].status)
	}

	return &apiResp, nil
}

func sendMail(conf *config, body string) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", conf.smtp.user, conf.smtp.pass, conf.smtp.host)
	sender := user.Username + "@" + hostname
	to := []string{conf.recipient}
	msg := []byte("To: " + conf.recipient + "\r\n" +
		"Subject: " + name + " alert\r\n" +
		"\r\n" +
		body + "\r\n")

	return smtp.SendMail(fmt.Sprintf("%s:%d", conf.smtp.host, conf.smtp.port), auth, sender, to, msg)
}
