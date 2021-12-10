package pkg

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func constructBackupFilename(targetContainer, targetDatabase string) string {
	backupFilenameStamp := time.Now().Format("2006_01_02__15_04_05")
	backupFilename := fmt.Sprintf(
		"%s_%s_%s.sql.gz",
		targetContainer,
		targetDatabase,
		backupFilenameStamp,
	)
	return backupFilename
}

func sendBackupToTelegram(data *Data, filename string) {
	file := &tb.Document{
		File:     tb.FromDisk(filename),
		Caption:  fmt.Sprintf("Compatible backend rev: %s", "TODO"), // FIXME
		FileName: filename,
	}

	data.Send(file)
}

func BackupDatabase(data *Data) {
	containerName := data.Config.ContainerName
	databaseName := data.Config.Database

	performDump(data, containerName, databaseName)
}

func PeriodicDBBackups(data *Data) {
	defer data.Send("Periodic database backups exited.")

	interval := data.Config.IntervalMinutes
	if interval <= 0 {
		return
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Minute)

	for range ticker.C {
		BackupDatabase(data)
	}
}
