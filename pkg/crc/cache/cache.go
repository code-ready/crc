package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/crc/version"
	"github.com/code-ready/crc/pkg/download"
	"github.com/code-ready/crc/pkg/embed"
	"github.com/code-ready/crc/pkg/extract"
	crcos "github.com/code-ready/crc/pkg/os"
	"github.com/pkg/errors"
)

type Cache struct {
	binaryName string
	archiveURL string
	destDir    string
	version    string
	getVersion func() (string, error)
}

type VersionMismatchError struct {
	BinaryName      string
	ExpectedVersion string
	CurrentVersion  string
}

func (e *VersionMismatchError) Error() string {
	return fmt.Sprintf("%s version mismatch: %s expected but %s found in the cache", e.BinaryName, e.ExpectedVersion, e.CurrentVersion)
}

func New(binaryName string, archiveURL string, destDir string, version string, getVersion func() (string, error)) *Cache {
	return &Cache{binaryName: binaryName, archiveURL: archiveURL, destDir: destDir, version: version, getVersion: getVersion}
}

func getCurrentOcVersion() (string, error) {
	ocPath := filepath.Join(constants.CrcOcBinDir, constants.OcBinaryName)
	stdOut, _, err := crcos.RunWithDefaultLocale(ocPath, "version", "--client")
	if len(strings.Split(stdOut, ":")) < 2 {
		return "", fmt.Errorf("Unable to parse the version information of %s", ocPath)
	}
	return strings.TrimSpace(strings.Split(stdOut, ":")[1]), err
}

func NewOcCache() *Cache {
	return New(constants.OcBinaryName, constants.GetOcURL(), constants.CrcOcBinDir, version.GetOcVersion(), getCurrentOcVersion)
}

func NewPodmanCache() *Cache {
	return New(constants.PodmanBinaryName, constants.GetPodmanURL(), constants.CrcBinDir, "", nil)
}

func NewGoodhostsCache() *Cache {
	return New(constants.GoodhostsBinaryName, constants.GetGoodhostsURL(), constants.CrcBinDir, "", nil)
}

func (c *Cache) IsCached() bool {
	if _, err := os.Stat(filepath.Join(c.destDir, c.binaryName)); os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *Cache) EnsureIsCached() error {
	if !c.IsCached() || c.CheckVersion() != nil {
		err := c.CacheBinary()
		if err != nil {
			return err
		}
	}
	return nil
}

// CacheBinary downloads and caches the requested binary into the CRC directory
func (c *Cache) CacheBinary() error {
	if c.IsCached() && c.CheckVersion() == nil {
		return nil
	}

	// Create tmp dir to download the requested tarball
	tmpDir, err := ioutil.TempDir("", "crc")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	assetTmpFile, err := c.getBinary(tmpDir)
	if err != nil {
		return err
	}

	var extractedFiles []string
	// Check the file is tarball or not
	if IsTarball(assetTmpFile) {
		// Extract the tarball and put it the cache directory.
		extractedFiles, err = extract.UncompressWithFilter(assetTmpFile, tmpDir, false,
			func(filename string) bool { return filepath.Base(filename) == c.binaryName })
		if err != nil {
			return errors.Wrapf(err, "Cannot uncompress '%s'", assetTmpFile)
		}
	} else {
		extractedFiles = append(extractedFiles, assetTmpFile)
		if filepath.Base(assetTmpFile) != c.binaryName {
			logging.Warnf("Binary name is %s but extracted file name is %s", c.binaryName, filepath.Base(assetTmpFile))
		}
	}

	// Copy the requested asset into its final destination
	err = os.MkdirAll(c.destDir, 0750)
	if err != nil && !os.IsExist(err) {
		return errors.Wrap(err, "Cannot create the target directory.")
	}

	for _, extractedFilePath := range extractedFiles {
		finalBinaryPath := filepath.Join(c.destDir, filepath.Base(extractedFilePath))
		// If the file exists then remove it (ignore error) first before copy because with `0500` permission
		// it is not possible to overwrite the file.
		os.Remove(finalBinaryPath)
		err = crcos.CopyFileContents(extractedFilePath, finalBinaryPath, 0500)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) getBinary(destDir string) (string, error) {
	logging.Debug("Trying to extract oc from crc binary")
	archiveName := filepath.Base(c.archiveURL)
	destPath := filepath.Join(destDir, archiveName)
	err := embed.Extract(archiveName, destPath)
	if err != nil {
		logging.Debug("Downloading oc")
		return download.Download(c.archiveURL, destDir, 0600)
	}

	return destPath, err
}

func (c *Cache) CheckVersion() error {
	// Check if version string is non-empty
	if c.version == "" {
		return nil
	}
	currentVersion, err := c.getVersion()
	if err != nil {
		return err
	}
	if currentVersion != c.version {
		return &VersionMismatchError{
			BinaryName:      c.binaryName,
			CurrentVersion:  currentVersion,
			ExpectedVersion: c.version,
		}
	}
	return nil
}

func IsTarball(filename string) bool {
	tarballExtensions := []string{".tar", ".tar.gz", ".tar.xz", ".zip", ".tar.bz2", ".crcbundle"}
	for _, extension := range tarballExtensions {
		if strings.HasSuffix(strings.ToLower(filename), extension) {
			return true
		}
	}
	return false
}
