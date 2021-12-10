package pkg

import (
	"bytes"
	"fmt"
	"io/ioutil"

	ssh "golang.org/x/crypto/ssh"
)

type SSHConnectionData struct {
	Hostname string
	Username string
	Keypath  string
}

var port = "22"

func executeCmd(command, hostname string, config *ssh.ClientConfig) (bytes.Buffer, error) {
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), config)
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(command)

	return stdoutBuf, nil
}

func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

func GetConfig(data *Data, sshData SSHConnectionData) (*ssh.ClientConfig, error) {
	// TODO this may cause runtime panic
	baseFileName := sshData.Keypath

	keyFilename := fmt.Sprintf("/app/config/%s", baseFileName)

	publicKey, err := PublicKeyFile(keyFilename)
	if err != nil {
		return nil, err
	}

	// TODO this may cause runtime panic
	user := sshData.Username
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// ssh.Password(p),
			publicKey,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config, nil
}

func ConnectAndExecute(
	data *Data,
	sshData SSHConnectionData,
	cmd string,
) (bytes.Buffer, error) {
	config, err := GetConfig(data, sshData)
	if err != nil {
		return bytes.Buffer{}, err
	}

	// TODO this may cause runtime panic
	output, err := executeCmd(cmd, sshData.Hostname, config)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return output, nil
}
