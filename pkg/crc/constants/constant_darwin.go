package constants

const (
	DefaultVMDriver = "hyperkit"
	OcBinaryName    = "oc"
	DefaultOcURL    = "https://mirror.openshift.com/pub/openshift-v4/clients/oc/latest/macosx/oc.tar.gz"
)

var (
	SupportedVMDrivers = [...]string{
		"virtualbox",
	}
)
