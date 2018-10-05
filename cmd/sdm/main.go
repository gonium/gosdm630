package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/gonium/gosdm630"
	. "github.com/gonium/gosdm630/internal/meters"
	"gopkg.in/urfave/cli.v1"
)

const (
	DEFAULT_METER_STORE_SECONDS = 120 * time.Second
)

func main() {
	app := cli.NewApp()
	app.Name = "sdm"
	app.Usage = "SDM MODBUS daemon"
	app.Version = RELEASEVERSION
	app.HideVersion = true
	app.Flags = []cli.Flag{
		// general
		cli.StringFlag{
			Name:  "serialadapter, s",
			Value: "/dev/ttyUSB0",
			Usage: "Serial RTU device",
		},
		cli.IntFlag{
			Name:  "comset, c",
			Value: ModbusComset9600_8N1,
			Usage: `Communication parameters:
		` + strconv.Itoa(ModbusComset2400_8N1) + `:  2400 baud, 8N1
		` + strconv.Itoa(ModbusComset9600_8N1) + `:  9600 baud, 8N1
		` + strconv.Itoa(ModbusComset19200_8N1) + `: 19200 baud, 8N1
		` + strconv.Itoa(ModbusComset2400_8E1) + `:  2400 baud, 8E1
		` + strconv.Itoa(ModbusComset9600_8E1) + `:  9600 baud, 8E1
		` + strconv.Itoa(ModbusComset19200_8E1) + `: 19200 baud, 8E1
			`,
		},
		cli.StringFlag{
			Name:  "device_list, d",
			Value: "SDM:1",
			Usage: `MODBUS device type and ID to query, separated by comma.
			Valid types are:
			"SDM" for Eastron SDM meters
			"JANITZA" for Janitza B-Series meters
			"DZG" for the DZG Metering GmbH DVH4013 meters
			"SBC" for the Saia Burgess Controls ALE3 meters
			Example: -d JANITZA:1,SDM:22,DZG:23`,
		},
		cli.BoolFlag{
			Name:  "detect",
			Usage: "Detect MODBUS devices.",
		},
		cli.StringFlag{
			Name:  "unique_id_format, f",
			Value: "Meter#%d",
			Usage: `Unique ID format.
			Example: -f Meter#%d
			The %d is replaced by the device ID`,
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print verbose messages",
		},

		// http api
		cli.StringFlag{
			Name:  "url, u",
			Value: ":8080",
			Usage: "Server listening socket",
		},

		// mqtt api
		cli.StringFlag{
			Name:  "broker, b",
			Value: "",
			Usage: "MQTT: Broker URI. ex: tcp://10.10.1.1:1883",
			// Destination: &mqttBroker,
		},
		cli.StringFlag{
			Name:  "topic, t",
			Value: "sdm630",
			Usage: "MQTT: Topic name to/from which to publish/subscribe (optional)",
			// Destination: &mqttTopic,
		},
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "MQTT: User (optional)",
			// Destination: &mqttUser,
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "MQTT: Password (optional)",
			// Destination: &mqttPassword,
		},
		cli.StringFlag{
			Name:  "clientid, i",
			Value: "sdm630",
			Usage: "MQTT: ClientID (optional)",
			// Destination: &mqttClientID,
		},
		cli.IntFlag{
			Name:  "rate, r",
			Value: 0,
			Usage: "MQTT: Maximum update rate (default 0, i.e. unlimited) (after a push we will ignore more data from same device and channel for this time)",
			// Destination: &mqttRate,
		},
		cli.BoolFlag{
			Name:  "clean, l",
			Usage: "MQTT: Set Clean Session (default false)",
			// Destination: &mqttCleanSession,
		},
		cli.IntFlag{
			Name:  "qos, q",
			Value: 0,
			Usage: "MQTT: Quality of Service 0,1,2 (default 0)",
			// Destination: &mqttQos,
		},
	}

	app.Action = func(c *cli.Context) {
		// Set unique ID format
		UniqueIdFormat = c.String("unique_id_format")

		// Parse the device_list parameter
		deviceslice := strings.Split(c.String("device_list"), ",")
		meters := make(map[uint8]*Meter)
		for _, meterdef := range deviceslice {
			splitdef := strings.Split(meterdef, ":")
			if len(splitdef) != 2 {
				log.Fatalf("Cannot parse device definition %s. See -h for help.", meterdef)
			}
			metertype, devid := splitdef[0], splitdef[1]
			id, err := strconv.Atoi(devid)
			if err != nil {
				log.Fatalf("Error parsing device id %s: %s. See -h for help.", meterdef, err.Error())
			}
			meter, err := NewMeterByType(metertype, uint8(id))
			if err != nil {
				log.Fatalf("Unknown meter type %s for device %d. See -h for help.", metertype, id)
			}
			meters[uint8(id)] = meter
		}

		// create ModbusEngine with status
		status := NewStatus(meters)
		qe := NewModbusEngine(
			c.String("serialadapter"),
			c.Int("comset"),
			c.Bool("verbose"),
			status,
		)

		// detect command
		if c.Bool("detect") {
			qe.Scan()
			return
		}

		// scheduler and meter data channel
		scheduler, snips := SetupScheduler(meters, qe)
		go scheduler.Run()

		// tee that broadcasts meter messages to multiple recipients
		tee := NewQuerySnipBroadcaster(snips)
		go tee.Run()

		// websocket hub
		hub := NewSocketHub(tee.Attach(), status)
		go hub.Run()

		// MQTT client
		if c.String("broker") != "" {
			mqtt := NewMqttClient(
				tee.Attach(),
				c.String("broker"),
				c.String("topic"),
				c.String("user"),
				c.String("password"),
				c.String("clientid"),
				c.Int("qos"),
				c.Int("rate"),
				c.Bool("clean"),
				c.Bool("verbose"))
			go mqtt.Run()
		}

		// MeasurementCache for REST API
		mc := NewMeasurementCache(
			meters,
			tee.Attach(),
			scheduler,
			DEFAULT_METER_STORE_SECONDS,
			c.Bool("verbose"),
		)
		go mc.Consume()

		Run_httpd(
			mc,
			hub,
			status,
			c.String("url"),
		)
	}

	app.Run(os.Args)
}
