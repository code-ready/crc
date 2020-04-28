// +build !windows

package cmd

import (
	"os"

	"github.com/code-ready/crc/pkg/crc/api"
	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/logging"
)

func runDaemon() {
	//setup separate logfile for daemon
	logging.CloseLogging()
	logging.InitLogrus(logging.LogLevel)
	logging.SetupFileHook(constants.DaemonLogFilePath)

	// Remove if an old socket is present
	os.Remove(constants.DaemonSocketPath)
	crcApiServer, err := api.CreateApiServer(constants.DaemonSocketPath)
	if err != nil {
		logging.Fatal("Failed to launch daemon", err)
	}
	crcApiServer.Serve()
}
