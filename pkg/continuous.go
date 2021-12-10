package pkg

import (
	"log"
	"time"
)

func ContinuousDBBackups(data *Data) {
	defer data.Send("Continuous database backups exited.")

	if !data.Config.Continuous {
		data.Send("Skipping continuous backups.")
		return
	}

	ticker := time.NewTicker(5 * time.Minute)

	// FIXME base backup should be performed on handle (check if continuity is enabled and wal files exist)
	// FIXME there should be handle (or periodic goroutine that cleans both basebackup and all wal files)
	for range ticker.C {
		err := PerformContinuity(data)
		if err != nil {
			data.Send("Continuous backup failed. Will retry later.")
			data.Send(err.Error())
			log.Println(err)
		}
	}
}
