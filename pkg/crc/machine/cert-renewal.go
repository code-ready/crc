package machine

import (
	"fmt"
	"regexp"
	"time"

	"github.com/code-ready/crc/pkg/crc/errors"
	"github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/crc/oc"
	"github.com/code-ready/crc/pkg/crc/ssh"
	"github.com/code-ready/crc/pkg/crc/systemd"
)

func waitForPendingCsrs(oc oc.OcConfig) error {
	waitForPendingCsr := func() error {
		output, _, err := oc.RunOcCommand("get", "csr")
		if err != nil {
			return &errors.RetriableError{Err: err}
		}
		matched, err := regexp.MatchString("Pending", output)
		if err != nil {
			return &errors.RetriableError{Err: err}
		}
		if !matched {
			return &errors.RetriableError{Err: fmt.Errorf("No Pending CSR")}
		}
		return nil
	}

	return errors.RetryAfter(60, waitForPendingCsr, time.Second)
}

func RegenerateCertificates(sshRunner *ssh.SSHRunner, machineName string) error {
	sd := systemd.NewInstanceSystemdCommander(sshRunner)
	startedKubelet, err := sd.Start("kubelet")
	if err != nil {
		logging.Debugf("Error starting kubelet service: %v", err)
		return err
	}
	if startedKubelet {
		defer sd.Stop("kubelet") //nolint:errcheck
	}
	oc := oc.UseOCWithConfig(machineName)
	/* 2 CSRs to approve, one right after kubelet restart, the other one a few dozen seconds after
	approving the first one
	- First one is requested by system:serviceaccount:openshift-machine-config-operator:node-bootstrapper
	- Second one is requested by system:node:<node_name> */
	err = waitForPendingCsrs(oc)
	if err != nil {
		logging.Debugf("Error waiting for first pending (node-bootstrapper) CSR: %v", err)
		return err
	}
	err = oc.ApproveNodeCSR()
	if err != nil {
		logging.Debugf("Error approving first pending (node-bootstrapper) CSR: %v", err)
		return err
	}

	err = waitForPendingCsrs(oc)
	if err != nil {
		logging.Debugf("Error waiting for second pending (system:node) CSR: %v", err)
		return err
	}
	err = oc.ApproveNodeCSR()
	if err != nil {
		logging.Debugf("Error approving second pending (system:node) CSR: %v", err)
		return err
	}

	return nil
}
