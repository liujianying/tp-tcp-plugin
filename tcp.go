package tp_tcp_plugin

import (
	"context"
	"encoding/binary"
	"log"
	"net"
	"github.com/sllt/tp-tcp-plugin/global"
	"github.com/sllt/tp-tcp-plugin/model"
	"github.com/sllt/tp-tcp-plugin/server"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	ctx := context.Background()
	for {
		buf := make([]byte, 8)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		payloadLen := binary.BigEndian.Uint32(buf[4:8])

		if !(buf[0] == 'T' && buf[1] == 'P') {
			log.Println("invalid protocol")
			// TODO : reply to client with error
			continue
		}

		// heartbeat packet
		if buf[2] == 0x1 {
			// TODO: handle heartbeat packet
			continue
		}

		payload := make([]byte, payloadLen)
		conn.Read(payload)

		packet := &model.Packet{}
		packet.Parse(append(buf[:n], payload...))

		log.Println(string(packet.Payload))

		device := &model.Device{}
		switch packet.Cmd {
		case 0x0:
			// auth
			log.Println("get auth request...")
			device = global.Devices[string(payload)]
			if device == nil || device.Online == false {
				device = &model.Device{
					AccessToken: string(payload),
				}
				err = device.Auth(global.MqttAddr)
				if err != nil {
					conn.Write([]byte("MQTT:" + err.Error()))
					continue
				}

				device.Online = true
				global.Devices[string(buf[:n])] = device
			}
			ctx = context.WithValue(ctx, "device", device)
		case 0x1:
			// publish
			log.Println("get publish request...")
			device = ctx.Value("device").(*model.Device)
			if device == nil || device.Online == false {
				conn.Write([]byte("MQTT: device auth failed"))
				continue
			}
			device.Publish(string(packet.Payload))
		case 0x2:
			// TODO publish events
		}
	}
}

func Start(addr string) {

	global.Devices = make(map[string]*model.Device)

	s := server.NewServer(addr)
	s.AddConnectionHandler(handleConnection)
	log.Println("started tcp server...")
	s.Start()
}
