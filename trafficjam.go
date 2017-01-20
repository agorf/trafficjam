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
	Origins      string `json:"origins"`
	Destinations string `json:"destinations"`
	APIKey       string `json:"api_key"`
	Mode         string `json:"mode"`
	Avoid        string `json:"avoid"`
	TrafficModel string `json:"traffic_model"`
	MaxDuration  int    `json:"max_duration"`
	SMTP         struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		User string `json:"user"`
		Pass string `json:"pass"`
	} `json:"smtp"`
	Recipient string `json:"recipient"`
}

type apiResponse struct {
	Rows []struct {
		Elements []struct {
			DurationInTraffic struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration_in_traffic"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
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
		"origins":        conf.Origins,
		"destinations":   conf.Destinations,
		"key":            conf.APIKey,
		"mode":           conf.Mode,
		"avoid":          conf.Avoid,
		"departure_time": "now",
		"traffic_model":  conf.TrafficModel,
	}

	apiResp, err := queryMapsAPI(params)
	if err != nil {
		log.Fatal(err)
	}

	duration := apiResp.Rows[0].Elements[0].DurationInTraffic.Value

	if duration > conf.MaxDuration*60 {
		sendMail(conf, apiResp.Rows[0].Elements[0].DurationInTraffic.Text)
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

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("%s: bad response status: %s\n", name, apiResp.Status)
	}
	if len(apiResp.Rows) != 1 {
		return nil, fmt.Errorf("%s: response row count is not 1\n", name)
	}
	if len(apiResp.Rows[0].Elements) != 1 {
		return nil, fmt.Errorf("%s: response first row element count is not 1\n", name)
	}
	if apiResp.Rows[0].Elements[0].Status != "OK" {
		return nil, fmt.Errorf("%s: bad response first row first element status: %s\n", name, apiResp.Rows[0].Elements[0].Status)
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

	auth := smtp.PlainAuth("", conf.SMTP.User, conf.SMTP.Pass, conf.SMTP.Host)
	sender := user.Username + "@" + hostname
	to := []string{conf.Recipient}
	msg := []byte("To: " + conf.Recipient + "\r\n" +
		"Subject: " + name + " alert\r\n" +
		"\r\n" +
		body + "\r\n")

	return smtp.SendMail(fmt.Sprintf("%s:%d", conf.SMTP.Host, conf.SMTP.Port), auth, sender, to, msg)
}
