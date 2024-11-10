package linuxopenports

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type OpenPort struct {
	Protocol string
	Port     int
}

// listOpenPorts reads /proc/net/tcp, /proc/net/udp, /proc/net/tcp6, and /proc/net/udp6 to find open ports
func GetOpenPorts() ([]OpenPort, error) {
	var openPorts []OpenPort
	uniquePorts := make(map[string]bool)

	protocolFiles := map[string][]string{
		"tcp": {"/proc/net/tcp", "/proc/net/tcp6"},
		"udp": {"/proc/net/udp", "/proc/net/udp6"},
	}

	for protocol, files := range protocolFiles {
		for _, filePath := range files {
			file, err := os.Open(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to open %s: %v", filePath, err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Scan()

			for scanner.Scan() {
				fields := strings.Fields(scanner.Text())
				if len(fields) < 2 {
					continue
				}

				localAddress := fields[1]
				addressParts := strings.Split(localAddress, ":")
				if len(addressParts) != 2 {
					continue
				}

				portHex := addressParts[1]
				port, err := strconv.ParseInt(portHex, 16, 32)
				if err != nil {
					continue
				}

				portKey := fmt.Sprintf("%s/%d", protocol, port)
				if !uniquePorts[portKey] {
					openPorts = append(openPorts, OpenPort{
						Protocol: protocol,
						Port:     int(port),
					})
					uniquePorts[portKey] = true
				}
			}
		}
	}

	return openPorts, nil
}
