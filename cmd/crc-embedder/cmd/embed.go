package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/constants/darwin"
	"github.com/code-ready/crc/pkg/crc/constants/linux"
	"github.com/code-ready/crc/pkg/crc/constants/windows"
	"github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/download"

	"github.com/code-ready/crc/pkg/crc/machine/hyperkit"
	"github.com/code-ready/crc/pkg/crc/machine/libvirt"

	"github.com/YourFin/binappend"
	"github.com/spf13/cobra"
)

var (
	bundleDir string
	goos      string
)

func init() {
	embedCmd.Flags().StringVar(&bundleDir, "bundle-dir", constants.MachineCacheDir, "Directory where the OpenShift bundle can be found")
	embedCmd.Flags().StringVar(&goos, "goos", runtime.GOOS, "Target platform (darwin, linux or windows)")
	rootCmd.AddCommand(embedCmd)
}

var embedCmd = &cobra.Command{
	Use:   "embed",
	Short: "Embed data files in crc binary",
	Long:  `Embed the OpenShift bundle and the binaries needed at runtime in the crc binary`,
	Run: func(cmd *cobra.Command, args []string) {
		runEmbed(args)
	},
}

func runEmbed(args []string) {
	if len(args) != 1 {
		logging.Fatalf("embed takes exactly one argument")
	}
	binaryPath := args[0]
	destDir, err := ioutil.TempDir("", "crc-embedder")
	if err != nil {
		logging.Errorf(fmt.Sprintf("Failed to create temporary directory: %v", err))
	}
	defer os.RemoveAll(destDir)
	downloadedFiles, err := downloadDataFiles(goos, destDir)
	if err != nil {
		logging.Errorf(fmt.Sprintf("Failed to download data files: %v", err))
	}

	bundlePath := path.Join(bundleDir, getDefaultBundleForOs(goos))
	downloadedFiles = append(downloadedFiles, bundlePath)
	err = embedFiles(binaryPath, downloadedFiles)
	if err != nil {
		logging.Errorf(fmt.Sprintf("Failed to embed data files: %v", err))
	}
}

func embedFiles(binary string, filenames []string) error {
	appender, err := binappend.MakeAppender(binary)
	if err != nil {
		return err
	}
	defer appender.Close()
	for _, filename := range filenames {
		logging.Debugf("Embedding %s in %s", filename, binary)
		f, err := os.Open(filename) // #nosec G304
		if err != nil {
			return fmt.Errorf("Failed to open %s: %v", filename, err)
		}
		defer f.Close()

		err = appender.AppendStreamReader(path.Base(filename), f, false)
		if err != nil {
			return fmt.Errorf("Failed to append %s to %s: %v", filename, binary, err)
		}
	}

	return nil
}

var defaultBundleForOs = map[string]string{
	"darwin":  darwin.DefaultBundle,
	"linux":   linux.DefaultBundle,
	"windows": windows.DefaultBundle,
}

func getDefaultBundleForOs(os string) string {
	return defaultBundleForOs[os]
}

var (
	dataFileUrls = map[string][]string{
		"darwin": []string{
			hyperkit.MachineDriverDownloadUrl,
			hyperkit.HyperkitDownloadUrl,
			darwin.OcUrl,
			darwin.GetCrcTrayUrl(),
			darwin.PodmanUrl,
		},
		"linux": []string{
			libvirt.MachineDriverDownloadUrl,
			linux.OcUrl,
			linux.PodmanUrl,
		},
		"windows": []string{
			windows.OcUrl,
		},
	}
)

func downloadDataFiles(goos string, destDir string) ([]string, error) {
	downloadedFiles := []string{}
	downloads := dataFileUrls[goos]
	for _, url := range downloads {
		filename, err := download.Download(url, destDir, 0644)
		if err != nil {
			return nil, err
		}
		downloadedFiles = append(downloadedFiles, filename)
	}

	return downloadedFiles, nil
}
