package provisioning

import (
	"time"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/storage"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/process"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"
	"github.com/sirupsen/logrus"
)

type SkipStep struct {
	skipID           string
	step             Step
	operationManager *process.ProvisionOperationManager
}

func NewSkipStep(os storage.Operations, skipID string, step Step) *SkipStep {
	return &SkipStep{
		skipID:           skipID,
		step:             step,
		operationManager: process.NewProvisionOperationManager(os),
	}
}

func (s *SkipStep) Name() string {
	return s.step.Name()
}

func (s *SkipStep) Run(operation internal.ProvisioningOperation, log logrus.FieldLogger) (internal.ProvisioningOperation, time.Duration, error) {
	pp, err := operation.GetProvisioningParameters()
	if err != nil {
		log.Errorf("cannot fetch provisioning parameters from operation: %s", err)
		return s.operationManager.OperationFailed(operation, "invalid operation provisioning parameters")
	}
	if pp.PlanID == s.skipID {
		log.Infof("Skipping step %s for the trial version", s.Name())
		return operation, 0, nil
	}

	return s.step.Run(operation, log)
}
