// +build !windows

package notification

import (
	"fmt"
	"runtime"
	"time"

	"github.com/code-ready/crc/test/extended/os/applescript"
)

type applescriptHandler struct {
}

const (
	scriptsRelativePath           string = "applescripts"
	manageNotifications           string = "manageNotifications.applescript"
	manageNotificationActionGet   string = "get"
	manageNotificationActionClear string = "clear"

	notificationDelay   string = "5s"
	notificationRetries int    = 36
)

func NewNotification() Notification {
	if runtime.GOOS == "darwin" {
		return applescriptHandler{}

	}
	return nil
}

func RequiredResourcesPath() (string, error) {
	return applescript.GetScriptsPath(scriptsRelativePath)
}

func (a applescriptHandler) GetClusterRunning() error {
	return checkNotificationMessage(startMessage)
}

func (a applescriptHandler) GetClusterStopped() error {
	return checkNotificationMessage(stopMessage)

}

func (a applescriptHandler) GetClusterDeleted() error {
	return checkNotificationMessage(deleteMessage)
}

func (a applescriptHandler) ClearNotifications() error {
	return applescript.ExecuteApplescript(manageNotifications, manageNotificationActionClear)
}

func checkNotificationMessage(notificationMessage string) error {
	retryDelay, err := time.ParseDuration(notificationDelay)
	if err != nil {
		return err
	}
	for i := 0; i < notificationRetries; i++ {
		err := applescript.ExecuteApplescriptReturnShouldMatch(
			notificationMessage, manageNotifications, manageNotificationActionGet)
		if err == nil {
			return nil
		}
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("notification: %s. Timeout", notificationMessage)
}
