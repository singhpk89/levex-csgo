package main

import (
	"fmt"
	"log"
	"net"
	"os"

	csgolog "github.com/janstuemmel/csgo-log"
)

func main() {
	//  Client data

	// listen to incoming udp packets
	udpServer, err := net.ListenPacket("udp", ":1053")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	for {
		buf := make([]byte, 1024)
		_, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go response(udpServer, addr, buf)
	}

}

func response(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	fileName := "output.log"
	// time := time.Now().Format(time.ANSIC)
	// responseStr := fmt.Sprintf("time received: %v. Your message: %v!", time, string(buf))
	line := string(buf)
	fmt.Println(line)
	data := parse(line)
	// responseStr := fmt.Sprintf("time received: %v. Your message: %v!", time, data)
	err := os.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	udpServer.WriteTo(cleanByteData([]byte(data)), addr)
	// udpServer.WriteTo([]byte(data), responseStr)
}

func parse(line string) string {

	var msg csgolog.Message

	// a line from a server logfile
	// line := echo "L 11/05/2018 - 15:44:36: Player<12><STEAM_1:1:0101011><CT> purchased m4a1" | nc -w10 -u 127.0.0.1 1053
	// line := `L 11/05/2017 - 15:17:35: "Mithilf<23><STEAM_1:0:84927935><TERRORIST>" [2443 -135 1613] killed "rix #Alienware<19><STEAM_1:1:13719626><CT>" [1842 821 1833] with "ak47" (headshot)`

	// echo `L 11/12/2018 - 19:58:52: "Scott<5><BOT><TERRORIST>" [1093 -850 1614] attacked "Jon<9><BOT><CT>" [972 -539 1614] with "glock" (damage "33") (damage_armor "0") (health "67") (armor "0") (hitgroup "stomach")` | nc -w10 -u 127.0.0.1 1053
	// line = `L 11/05/2017 - 15:12:03: "rix #Alienware<6><STEAM_1:1:13719626><CT>" purchased "m4a1_silencer"`
	// parse into Message
	msg, err := csgolog.Parse(line)

	if err != nil {
		fmt.Println(err)
		return "Invalid data"
	} else {

		// get json non-htmlescaped
		jsn := csgolog.ToJSON(msg)

		fmt.Println(jsn)
		return jsn
	}

}

func cleanByteData(input []byte) []byte {
	cleanedData := make([]byte, 0, len(input))
	for _, b := range input {
		if b != 0 {
			cleanedData = append(cleanedData, b)
		}
	}
	return cleanedData
}
