package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	sdm630 "github.com/gonium/gosdm630"
	"gopkg.in/urfave/cli.v1"
)

// Copied from
// https://github.com/jcuga/golongpoll/blob/master/events.go:
type lpEvent struct {
	// Timestamp is milliseconds since epoch to match javascrits Date.getTime()
	Timestamp int64  `json:"timestamp"`
	Category  string `json:"category"`
	// NOTE: Data can be anything that is able to passed to json.Marshal()
	Data sdm630.QuerySnip `json:"data"`
}

// eventResponse is the json response that carries longpoll events.
type eventResponse struct {
	Events *[]lpEvent `json:"events"`
}

func main() {
	var sdmURL string
	var sdmTimeout int
	var verbose bool
	var mqttTopic string
	var mqttBroker string
	var mqttPassword string
	var mqttUser string
	var mqttClientID string
	var mqttQos int
	var mqttCleanses bool

	app := cli.NewApp()
	app.Name = "sdm630_mqtt"
	app.Usage = "SDM630 mqtt"
	app.Version = sdm630.RELEASEVERSION
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url, u",
			Value: "localhost:8080",
			Usage: "the URL of the server we should connect to",
			//			FilePath:    "~/.sdm630_mqtt.conf",
			Destination: &sdmURL,
		},
		cli.IntFlag{
			Name:  "timeout",
			Value: 45,
			Usage: "timeout value in seconds",
			//			FilePath:    "~/.sdm630_mqtt.conf",
			Destination: &sdmTimeout,
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print verbose messages",
			//			FilePath:    "~/.sdm630_mqtt.conf",
			Destination: &verbose,
		},

		cli.StringFlag{
			Name:  "topic, t",
			Value: "gosdm630/",
			Usage: "MQTT: The topic name to/from which to publish/subscribe",
			//			FilePath:    "~/.sdm630_mqtt.conf",
			Destination: &mqttTopic,
		},
		cli.StringFlag{
			Name:  "broker, b",
			Value: "",
			Usage: "MQTT: The broker URI. ex: tcp://10.10.1.1:1883",
			//			FilePath:    "~/.sdm630_mqtt.conf",
			Destination: &mqttBroker,
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "MQTT: The password (optional)",
			//			FilePath:    "~/.sdm630_mqtt.conf",
			Destination: &mqttPassword,
		},
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "MQTT: The User (optional)",
			//			FilePath:    "~/.sdm630_mqtt.conf",
			Destination: &mqttUser,
		},
		cli.StringFlag{
			Name:  "id, i",
			Value: "gosdm630",
			Usage: "MQTT: The ClientID (optional)",
			//			FilePath:    "~/.sdm630_mqtt.conf",
			Destination: &mqttClientID,
		},
		cli.BoolFlag{
			Name:  "clean,´c",
			Usage: "MQTT: Set Clean Session (default false)",
			//			FilePath:    "~/.sdm630_mqtt.conf",
			Destination: &mqttCleanses,
		},
		cli.IntFlag{
			Name:  "qos, q",
			Value: 0,
			Usage: "MQTT: The Quality of Service 0,1,2 (default 0)",
			//			FilePath:    "~/.sdm630_mqtt.conf",
			Destination: &mqttQos,
		},
	}
	app.Action = func(c *cli.Context) {
		endpointURL := fmt.Sprintf("http://%s/firehose?timeout=%d&category=meterupdate", sdmURL, sdmTimeout)
		if verbose {
			log.Printf("Client startup - will connect to %s", endpointURL)
		}
		client := &http.Client{
			Timeout: time.Duration(sdmTimeout) * time.Second,
			Transport: &http.Transport{
				// 0 means: no limit.
				MaxIdleConns:        0,
				MaxIdleConnsPerHost: 0,
				IdleConnTimeout:     0,
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: time.Minute,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
				DisableKeepAlives:   false,
			},
		}

		mqttOpts := MQTT.NewClientOptions()
		mqttOpts.AddBroker(mqttBroker)
		mqttOpts.SetClientID(mqttClientID)
		mqttOpts.SetUsername(mqttUser)
		mqttOpts.SetPassword(mqttPassword)
		mqttOpts.SetCleanSession(mqttCleanses)
		topic := fmt.Sprintf("%s/connected", mqttTopic)
		message := fmt.Sprintf("0")
		mqttOpts.SetWill(topic, message, byte(mqttQos), true)

		if verbose {
			log.Println("Connecting to mqtt:")
			log.Printf("\tbroker:    %s\n", mqttBroker)
			log.Printf("\tclientid:  %s\n", mqttClientID)
			log.Printf("\tuser:      %s\n", mqttUser)
			if mqttPassword != "" {
				log.Printf("\tpassword:  ****\n")
			}
			log.Printf("\ttopic:     %s\n", mqttTopic)
			log.Printf("\tqos:       %d\n", mqttQos)
			log.Printf("\tcleansess: %v\n", mqttCleanses)
		}

		mqttClient := MQTT.NewClient(mqttOpts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal("Error connecting to mqtt: ", token.Error())
			panic(token.Error())
		}
		if verbose {
			log.Println("Connected to mqtt")
		}
		topic = fmt.Sprintf("%s/connected", mqttTopic)
		message = fmt.Sprintf("1")
		token := mqttClient.Publish(topic, byte(mqttQos), true, message)
		if verbose {
			log.Printf("MQTT push %s, Message: %s", topic, message)
		}
		if token.Wait() && token.Error() != nil {
			log.Fatal("Error connecting to mqtt: ", token.Error())
			panic(token.Error())
		}

		resp, err := client.Get(endpointURL)
		if err != nil {
			log.Fatal("Failed to read from endpoint: ", err.Error())
		} else {
			if verbose {
				log.Println("Connected to firehose")
			}
			topic = fmt.Sprintf("%s/connected", mqttTopic)
			message = fmt.Sprintf("2")
			token = mqttClient.Publish(topic, byte(mqttQos), true, message)
			if verbose {
				log.Printf("MQTT push %s, Message: %s", topic, message)
			}
			if token.Wait() && token.Error() != nil {
				log.Fatal("Error connecting to mqtt: ", token.Error())
				panic(token.Error())
			}
		}

		for {
			rawevents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Failed to process message: ", err.Error())
			} else {
				// handle the events.
				var events eventResponse
				err := json.Unmarshal(rawevents, &events)
				if err != nil {
					log.Fatal("Failed to decode JSON events: ", err.Error())
				}
				for _, event := range *events.Events {
					snip := event.Data
					if verbose {
						log.Printf("Device %d: Data: %s, Value: %.3f W, Desc: %s", snip.DeviceId, snip.IEC61850, snip.Value, snip.Description)
					}
					topic = fmt.Sprintf("%s/status/%d/%s", mqttTopic, snip.DeviceId, snip.IEC61850)
					message = fmt.Sprintf("%.3f", snip.Value)
					token = mqttClient.Publish(topic, byte(mqttQos), true, message)
					if verbose {
						log.Printf("MQTT: push %s, Message: %s", topic, message)
					}
					if token.Wait() && token.Error() != nil {
						log.Fatal("Error connecting to mqtt: ", token.Error())
						panic(token.Error())
					}
				}
			}
			if resp.Body != nil {
				resp.Body.Close()
			}

			resp, err = client.Get(endpointURL)
			if err != nil {
				log.Fatal("Failed to read from endpoint: ", err.Error())
				topic = fmt.Sprintf("%s/connected", mqttTopic)
				message = fmt.Sprintf("1")
				token = mqttClient.Publish(topic, byte(mqttQos), true, message)
				if verbose {
					log.Printf("MQTT push %s, Message: %s", topic, message)
				}
				if token.Wait() && token.Error() != nil {
					log.Fatal("Error connecting to mqtt: ", token.Error())
					panic(token.Error())
				}
			}
		}
	}
	app.Run(os.Args)
}
