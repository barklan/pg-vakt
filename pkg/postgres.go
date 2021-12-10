package pkg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func backupPostgresCmdString(containerName, databaseName string) string {
	cmd := fmt.Sprintf(
		"docker exec $(docker ps -q -f name=%s) pg_dump -U postgres %s | gzip",
		containerName,
		databaseName,
	)
	return cmd
}

func performDump(data *Data, containerName, databaseName string) {
	cmd := backupPostgresCmdString(containerName, databaseName)

	sshData := SSHConnectionData{
		Hostname: data.Config.SSHHostname,
		Username: data.Config.SSHUser,
		Keypath:  data.Config.SSHKeyFilename,
	}

	buffer, err := ConnectAndExecute(data, sshData, cmd)
	if err != nil {
		data.Send("Error when making backup.")
		data.Send(err.Error())
		return
	}

	baseFileName := constructBackupFilename(containerName, databaseName)
	data.CreateMediaDirIfNotExists("project/sqldumps")
	fullFilename := data.MediaPath + "/project/sqldumps/" + baseFileName
	err = ioutil.WriteFile(fullFilename, buffer.Bytes(), 0777)
	if err != nil {
		data.Send("Failed to save backup.")
		log.Println(err)
	}

	// TODO consider if this is usefull at all
	// sendBackupToTelegram(data, fullFilename)
}

func baseBackupPostgresCmdString(containerName string) string {
	cmd := fmt.Sprintf(
		"docker exec $(docker ps -q -f name=%s) pg_basebackup -U postgres -D /pgbackups -Ft -z",
		containerName,
	)
	return cmd
}

func PerformBaseBackup(data *Data) {
	// FIXME
}

func PerformContinuity(data *Data) error {
	sshData := SSHConnectionData{
		Hostname: data.Config.SSHHostname,
		Username: data.Config.SSHUser,
		Keypath:  data.Config.SSHKeyFilename,
	}

	remoteFolder := data.Config.ContinuousPath
	output, err := ConnectAndExecute(
		data,
		sshData,
		fmt.Sprintf("ls -p %s | grep -v /", remoteFolder),
	)
	if err != nil {
		return err
	}

	remoteFiles := strings.Fields(output.String())
	log.Printf("Remote files: %s", remoteFiles)

	localFiles := make([]string, 0)

	data.CreateMediaDirIfNotExists("project/pgbackups")
	localFolder := data.MediaPath + "/project/pgbackups"
	files, err := ioutil.ReadDir(localFolder)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			localFiles = append(localFiles, file.Name())
		}
	}

	log.Printf("Local files: %s", localFiles)

	filesToCopy := difference(remoteFiles, localFiles)
	log.Printf("Files to copy: %s", filesToCopy)

	if len(filesToCopy) == 0 {
		return nil
	}

	containerName := data.Config.ContainerName

	for _, file := range filesToCopy {
		remoteFilePath := remoteFolder + "/" + file
		localFilePath := localFolder + "/" + file

		log.Printf("Transferring %s", file)

		err = SFTP(data, sshData, remoteFilePath, localFilePath)
		if err != nil {
			_, e := os.Stat(localFilePath)
			if errors.Is(e, os.ErrNotExist) {
				return err
			}

			e = os.Remove(localFilePath)
			if e != nil {
				data.Send("Continous backup is corrupted! Please intervene. Panic in 10 seconds.")
				time.Sleep(10 * time.Second)
				log.Panic(e)
			}

			return err
		}

		log.Printf("Removing remote  %s", file)

		_, err = ConnectAndExecute(
			data,
			sshData,
			fmt.Sprintf("rm %s", remoteFilePath),
		)
		if err != nil {
			data.Send(
				fmt.Sprintf("%s. Failed to remove copied WAL files from server.", containerName),
			)
		}
	}

	return nil
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
