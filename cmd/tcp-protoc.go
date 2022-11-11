package main

import (
	"flag"
	plugin "github.com/sllt/tp-tcp-plugin"
	"github.com/sllt/tp-tcp-plugin/global"
)

func main() {
	var (
		listenAddr = flag.String("listen", ":8080", "listen address")
		mqttAddr   = flag.String("mqtt", "tcp://dev.thingspanel.cn:1883", "mqtt broker address")
	)

	flag.Parse()

	global.MqttAddr = *mqttAddr
	plugin.Start(*listenAddr)
}
