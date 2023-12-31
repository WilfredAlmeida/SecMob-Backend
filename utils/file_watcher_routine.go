package utils

import (
	"github.com/fsnotify/fsnotify"
)

func WatchFile(fileName string, callBack func()) {

	var watcher, err = fsnotify.NewWatcher()

	if err != nil {
		MLogger.ErrorLog(err.Error())
	}

	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {

		}
	}(watcher)

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				MLogger.InfoLog("event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					MLogger.InfoLog("modified file:", event.Name)
					callBack()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				MLogger.ErrorLog(err.Error())
			}
		}
	}()

	err = watcher.Add(fileName)
	if err != nil {
		MLogger.ErrorLog(err.Error())
	}
	<-done
}
