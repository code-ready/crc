package config

import (
	validations "github.com/code-ready/crc/pkg/crc/config"
	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/spf13/cobra"
	"strings"
)

// validationFnType takes the key, value as args and checks if valid
type validationFnType func(interface{}) (bool, string)

type setting struct {
	Name          string
	DefaultValue  interface{}
	ValidationFns []validationFnType
}

// SettingsList holds all the config settings
var SettingsList = make(map[string]*setting)

var (
	// Start command settings in config
	VMDriver = createSetting("vm-driver", constants.DefaultVMDriver, []validationFnType{validations.ValidateDriver})
	Bundle   = createSetting("bundle", nil, []validationFnType{validations.ValidateBundle})
	CPUs     = createSetting("cpus", constants.DefaultCPUs, []validationFnType{validations.ValidateCPUs})
	Memory   = createSetting("memory", constants.DefaultMemory, []validationFnType{validations.ValidateMemory})

	// Preflight checks
	SkipCheckVirtEnabled             = createSetting("skip-check-virt-enabled", nil, []validationFnType{validations.ValidateBool})
	WarnCheckVirtEnabled             = createSetting("warn-check-virt-enabled", nil, []validationFnType{validations.ValidateBool})
	SkipCheckKvmEnabled              = createSetting("skip-check-kvm-enabled", nil, []validationFnType{validations.ValidateBool})
	WarnCheckKvmEnabled              = createSetting("warn-check-kvm-enabled", nil, []validationFnType{validations.ValidateBool})
	SkipCheckLibvirtInstalled        = createSetting("skip-check-libvirt-installed", nil, []validationFnType{validations.ValidateBool})
	WarnCheckLibvirtInstalled        = createSetting("warn-check-libvirt-installed", nil, []validationFnType{validations.ValidateBool})
	SkipCheckLibvirtEnabled          = createSetting("skip-check-libvirt-enabled", nil, []validationFnType{validations.ValidateBool})
	WarnCheckLibvirtEnabled          = createSetting("warn-check-libvirt-enabled", nil, []validationFnType{validations.ValidateBool})
	SkipCheckLibvirtRunning          = createSetting("skip-check-libvirt-running", nil, []validationFnType{validations.ValidateBool})
	WarnCheckLibvirtRunning          = createSetting("warn-check-libvirt-running", nil, []validationFnType{validations.ValidateBool})
	SkipCheckUserInLibvirtGroup      = createSetting("skip-check-user-in-libvirt-group", nil, []validationFnType{validations.ValidateBool})
	WarnCheckUserInLibvirtGroup      = createSetting("warn-check-user-in-libvirt-group", nil, []validationFnType{validations.ValidateBool})
	SkipCheckLibvirtDriver           = createSetting("skip-check-libvirt-driver", nil, []validationFnType{validations.ValidateBool})
	WarnCheckLibvirtDriver           = createSetting("warn-check-libvirt-driver", nil, []validationFnType{validations.ValidateBool})
	SkipCheckCrcNetwork              = createSetting("skip-check-crc-network", nil, []validationFnType{validations.ValidateBool})
	WarnCheckCrcNetwork              = createSetting("warn-check-crc-network", nil, []validationFnType{validations.ValidateBool})
	SkipCheckCrcNetworkActive        = createSetting("skip-check-crc-network-active", nil, []validationFnType{validations.ValidateBool})
	WarnCheckCrcNetworkActive        = createSetting("warn-check-crc-network-active", nil, []validationFnType{validations.ValidateBool})
	SkipCheckCrcDnsmasqFile          = createSetting("skip-check-crc-dnsmasq-file", nil, []validationFnType{validations.ValidateBool})
	WarnCheckCrcDnsmasqFile          = createSetting("warn-check-crc-dnsmasq-file", nil, []validationFnType{validations.ValidateBool})
	SkipCheckVirtualBoxInstalled     = createSetting("skip-check-virtualbox-installed", nil, []validationFnType{validations.ValidateBool})
	WarnCheckVirtualBoxInstalled     = createSetting("warn-check-virtualbox-installed", nil, []validationFnType{validations.ValidateBool})
	SkipCheckCrcNetworkManagerConfig = createSetting("skip-check-network-manager-config", nil, []validationFnType{validations.ValidateBool})
	WarnCheckCrcNetworkManagerConfig = createSetting("warn-check-network-manager-config", nil, []validationFnType{validations.ValidateBool})
)

// CreateSetting returns a filled struct of ConfigSetting
// takes the config name and default value as arguments
func createSetting(name string, defValue interface{}, validationFn []validationFnType) *setting {
	s := setting{Name: name, DefaultValue: defValue, ValidationFns: validationFn}
	SettingsList[name] = &s
	return &s
}

var (
	ConfigCmd = &cobra.Command{
		Use:   "config SUBCOMMAND [flags]",
		Short: "Modifies crc configuration properties.",
		Long: `Modifies crc configuration properties. Some of the configuration properties are equivalent
to the options that you set when you run the 'crc start' command.
Configurable properties (enter as SUBCOMMAND): ` + "\n\n" + configurableFields(),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func configurableFields() string {
	var fields []string
	for _, s := range SettingsList {
		fields = append(fields, " * "+s.Name)
	}
	return strings.Join(fields, "\n")
}
