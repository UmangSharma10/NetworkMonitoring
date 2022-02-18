package plugins

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
	"time"
)

type processDetails struct {
	user    string
	pid     int
	cpu     float64
	mem     float64
	VSZ     int
	RSS     int
	TTY     string
	stat    string
	start   string
	time    string
	command string
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func CollectSSH(credMaps map[string]string, mychan chan string) {

	port, _ := strconv.Atoi(credMaps["port"])
	sshHost := credMaps["host"]
	sshUser := credMaps["user"]
	sshPassword := credMaps["password"]
	sshPort := port
	// Create SSHP login configuration
	config := &ssh.ClientConfig{
		Timeout:         10 * time.Second, //ssh connection time out time is one second, if SSH validation error returns in one second
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{Ciphers: []string{
			"aes128-ctr", "aes192-ctr", "aes256-ctr",
		}},
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}

	// dial gets SSH client
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Println(err)
	}
	defer func(sshClient *ssh.Client) {
		err := sshClient.Close()
		if err != nil {

		}
	}(sshClient)

	//Disk Session
	session, err := sshClient.NewSession()
	if err != nil {
		panic(err)
	}
	disk, err := session.Output("df -h")
	if err != nil {
		panic(err)
	}

	dsstring := string(disk)

	mystringDisk := strings.Split(dsstring, "\n")
	var getsize []map[string]string
	flag := 1
	for _, v := range mystringDisk {
		if flag == 1 {
			flag = 0
			continue
		}

		split1 := strings.SplitN(standardizeSpaces(v), " ", 6)
		if len(split1) < 6 {
			break
		}

		temp1 := map[string]string{
			"FileSystem": split1[0],
			"Size":       split1[1],
			"Used":       split1[2],
			"Mounted":    split1[5],
		}

		getsize = append(getsize, temp1)
	}

	session, err = sshClient.NewSession()
	if err != nil {
		panic(err)
	}
	psaux, err := session.Output("ps aux")
	if err != nil {
		panic(err)
	}
	psauxString := string(psaux)

	myStringArray := strings.Split(psauxString, "\n")

	flag = 1
	var getsize2 []map[string]string

	for _, v := range myStringArray {

		if flag == 1 {
			flag = 0
			continue
		}
		splitN := strings.SplitN(standardizeSpaces(v), " ", 11)
		if len(splitN) <= 10 {
			break
		}
		processPID := splitN[1]
		if err != nil {
			panic(err)
		}
		processCPU := splitN[2]
		if err != nil {
			panic(err)
		}
		processMEM := splitN[3]
		processRSS := splitN[5]
		if err != nil {
			panic(err)
		}

		temp1 := map[string]string{
			"process.user": splitN[0], "process.pid": (processPID),
			"process.MEM": processMEM,
			"process.RSS": (processRSS),
			"process.cpu": processCPU,
		}

		getsize2 = append(getsize2, temp1)

	}

	result := map[string]interface{}{
		"Disk":    getsize,
		"Process": getsize2,
	}

	myJson, _ := json.MarshalIndent(result, "", "    ")
	mychan <- string(myJson)
	//for k, v := range processMap {
	//	fmt.Printf("PID: %d CPU: %g command: %s  \n", k, v.cpu, v.command)
	//}
	err = session.Close()
	if err != nil {
		return
	}

}
