package config

type MachineConfig struct {
	// CRC system bundle
	BundleName string

	// Virtual machine configuration
	Name            string
	Memory          int
	CPUs            int
	DiskPath        string
	ImageSourcePath string
	SSHKeyPath      string

	// HyperKit specific configuration
	KernelCmdLine string
	Initramfs     string
	Kernel        string
}
