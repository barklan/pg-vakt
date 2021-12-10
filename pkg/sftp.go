package pkg

import (
	"fmt"
	"os"

	"github.com/pkg/sftp"
	ssh "golang.org/x/crypto/ssh"
)

func SFTP(
	data *Data,
	sshData SSHConnectionData,
	remoteFile,
	localFile string,
) error {
	clientConfig, err := GetConfig(data, sshData)
	if err != nil {
		return err
	}
	hostname := sshData.Hostname

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), clientConfig)
	if err != nil {
		return err
	}
	defer conn.Close()
	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer sftp.Close()

	// Open the source file
	srcFile, err := sftp.Open(remoteFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(localFile)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the file
	srcFile.WriteTo(dstFile)
	return nil
}
