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
	app.Name = "sdm630_httpd"
	app.Usage = "SDM630 power measurements via HTTP."
	app.Version = RELEASEVERSION
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "serialadapter, s",
			Value: "/dev/ttyUSB0",
			Usage: "path to serial RTU device",
		},
		cli.IntFlag{
			Name:  "comset, c",
			Value: ModbusComset9600_8N1,
			Usage: `which communication parameter set to use. Valid sets are
		` + strconv.Itoa(ModbusComset2400_8N1) + `:  2400 baud, 8N1
		` + strconv.Itoa(ModbusComset9600_8N1) + `:  9600 baud, 8N1
		` + strconv.Itoa(ModbusComset19200_8N1) + `: 19200 baud, 8N1
		` + strconv.Itoa(ModbusComset2400_8E1) + `:  2400 baud, 8E1
		` + strconv.Itoa(ModbusComset9600_8E1) + `:  9600 baud, 8E1
		` + strconv.Itoa(ModbusComset19200_8E1) + `: 19200 baud, 8E1
			`,
		},
		cli.StringFlag{
			Name:  "url, u",
			Value: ":8080",
			Usage: "the URL the server should respond on",
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print verbose messages",
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
		cli.StringFlag{
			Name:  "unique_id_format, f",
			Value: "Instrument%d",
			Usage: `Unique ID format.
			Example: -f Instrument%d
			The %d is replaced by the device ID`,
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

		// scheduler and meter data channel
		scheduler, snips := SetupScheduler(meters, qe)
		go scheduler.Run()

		// tee that broadcasts meter messages to multiple recipients
		tee := NewQuerySnipBroadcaster(snips)
		go tee.Run()

		// Longpoll firehose
		firehose := NewFirehose(
			tee.Attach(),
			status,
			c.Bool("verbose"))
		go firehose.Run()

		// websocket hub
		hub := NewSocketHub(tee.Attach(), status)
		go hub.Run()

		// MeasurementCache for REST API
		mc := NewMeasurementCache(
			meters,
			tee.Attach(),
			scheduler,
			DEFAULT_METER_STORE_SECONDS,
			c.Bool("verbose"),
		)
		go mc.Consume()

		log.Printf("Starting API httpd at %s", c.String("url"))
		Run_httpd(
			mc,
			firehose,
			hub,
			status,
			c.String("url"),
		)
	}

	app.Run(os.Args)
}
