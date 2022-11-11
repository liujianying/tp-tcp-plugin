package tp_tcp_plugin_test

import (
	"fmt"
	"net"
	"time"
	"github.com/sllt/tp-tcp-plugin/model"
)

func Example_client() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	isAuthed := false
	for {
		if !isAuthed {
			packet := model.BuildAuthPacket("2072aae3-596c-3977-4a14-b127dcef41e0")
			data := packet.Serialize()

			fmt.Println(string(data))

			_, err := conn.Write(data)
			if err != nil {
				fmt.Println(err)
				break
			}
			isAuthed = true
		}

		publish := "{\"light\": 98, \"humidity\": 30.0}"

		packet := model.BuildPublishAttributesPacket([]byte(publish))
		data := packet.Serialize()

		_, err := conn.Write(data)
		if err != nil {
			fmt.Println(err)
			break
		}

		time.Sleep(time.Second * 2)
	}
}
