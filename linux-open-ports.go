package linuxopenports

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type OpenPort struct {
	Protocol string
	Port     int
	PID      string
}

func GetOpenPorts() ([]OpenPort, error) {
	var openPorts []OpenPort
	uniquePorts := make(map[string]bool)

	protocolFiles := map[string][]string{
		"tcp": {"/proc/net/tcp", "/proc/net/tcp6"},
		"udp": {"/proc/net/udp", "/proc/net/udp6"},
	}

	cachedInodePIDMap := inodePIDMap()

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
				if len(fields) < 10 {
					continue
				}

				localAddress := fields[1]
				inode := fields[9]
				addressParts := strings.Split(localAddress, ":")
				if len(addressParts) != 2 {
					continue
				}

				portHex := addressParts[1]
				port, err := strconv.ParseInt(portHex, 16, 32)
				if err != nil {
					continue
				}

				pid, ok := cachedInodePIDMap[inode]
				if !ok {
					continue
				}

				portKey := fmt.Sprintf("%s/%d", protocol, port)
				if !uniquePorts[portKey] {
					openPorts = append(openPorts, OpenPort{
						Protocol: protocol,
						Port:     int(port),
						PID:      pid,
					})
					uniquePorts[portKey] = true
				}
			}
		}
	}

	return openPorts, nil
}

func inodePIDMap() map[string]string {
	m := map[string]string{}
	procDirs, _ := os.ReadDir("/proc")
	for _, procDir := range procDirs {
		pid := procDir.Name()
		if !procDir.IsDir() && !unicode.IsDigit(rune(pid[0])) {
			continue
		}

		fdDir := filepath.Join("/proc", pid, "fd")
		fdFiles, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fdFile := range fdFiles {
			path := filepath.Join(fdDir, fdFile.Name())
			linkName, err := os.Readlink(path)
			if err != nil {
				continue
			}
			if strings.Contains(linkName, "socket") {
				inode := linkName[8 : len(linkName)-1]
				m[inode] = pid
			}
		}
	}
	return m
}
