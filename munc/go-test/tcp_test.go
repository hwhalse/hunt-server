package go_test

import (
	"bufio"
	"encoding/json"
	"eolian/munc/constants"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func startServer() net.Listener {
	server, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Println("Error starting server", err)
	}
	fmt.Println("TCP server listening on port 3000")
	return server
}

func TestMain(m *testing.M) {
	server := startServer()
	exitCode := m.Run()
	server.Close()
	os.Exit(exitCode)
}

type InitMsg struct {
	MuncEventId int    `json:"muncEventId"`
	Id          string `json:"id"`
}

func TestConnectionEstablishment(t *testing.T) {
	client, err := net.Dial("tcp", "esp-dev.eastus2.cloudapp.azure.com:5000")
	if err != nil {
		t.Fatal("Error establishing conn", err)
	}
	initMsg := InitMsg{
		MuncEventId: constants.TCP_INIT,
		Id:          uuid.NewString(),
	}
	m, err := json.Marshal(initMsg)
	client.Write(m)

	go handleTcpResponse(client)

	defer client.Close()
}

func tcpSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {

		return 0, nil, nil
	}
	if i := strings.Index(string(data), "/n"); i >= 0 {
		return i + 2, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func handleTcpResponse(client net.Conn) {
	scanner := bufio.NewScanner(client)
	scanner.Split(tcpSplit)
	var splitMsg string
	var mg string
	for scanner.Scan() {
		splitMsg = scanner.Text()
		mg = scanner.Text()
		var m map[string]any
		err := json.Unmarshal([]byte(splitMsg), &m)
		fmt.Println(mg)
		if err != nil {
			fmt.Println("err", err)
		}
	}
}
