package provisioning

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/cls"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/auditlog"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/process"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/storage"
	"github.com/kyma-project/control-plane/components/provisioner/pkg/gqlschema"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type CLSAuditLogOverrides struct {
	operationManager *process.ProvisionOperationManager
	fs               afero.Fs
	auditLogConfig   auditlog.Config
	secretKey        string
}

func (alo *CLSAuditLogOverrides) Name() string {
	return "CLS_Audit_Log_Overrides"
}

func NewCLSAuditLogOverridesStep(os storage.Operations, cfg auditlog.Config, secretKey string) *CLSAuditLogOverrides {
	fileSystem := afero.NewOsFs()

	return &CLSAuditLogOverrides{
		process.NewProvisionOperationManager(os),
		fileSystem,
		cfg,
		secretKey,
	}
}

func (alo *CLSAuditLogOverrides) Run(operation internal.ProvisioningOperation, logger logrus.FieldLogger) (internal.ProvisioningOperation, time.Duration, error) {
	luaScript, err := alo.readFile("/auditlog-script/script")
	if err != nil {
		logger.Errorf("Unable to read audit config script: %v", err)
		return operation, 0, err
	}

	replaceSubAccountID := strings.Replace(string(luaScript), "sub_account_id", operation.ProvisioningParameters.ErsContext.SubAccountID, -1)
	replaceTenantID := strings.Replace(replaceSubAccountID, "tenant_id", alo.auditLogConfig.Tenant, -1)

	u, err := url.Parse(alo.auditLogConfig.URL)
	if err != nil {
		logger.Errorf("Unable to parse the URL: %v", err.Error())
		return operation, 0, err
	}
	if u.Path == "" {
		logger.Errorf("There is no Path passed in the URL")
		return operation, 0, errors.New("there is no Path passed in the URL")
	}
	auditLogHost, auditLogPort, err := net.SplitHostPort(u.Host)
	if err != nil {
		logger.Errorf("Unable to split URL: %v", err.Error())
		return operation, 0, err
	}
	if auditLogPort == "" {
		auditLogPort = "443"
		logger.Infof("There is no Port passed in the URL. Setting default to 443")
	}
	fluentbitPlugin := "http"
	if alo.auditLogConfig.EnableSeqHttp {
		fluentbitPlugin = "sequentialhttp"
	}

	clsOverrides, err := cls.DecryptOverrides(alo.secretKey, operation.Cls.Overrides)
	if err != nil {
		logger.Errorf("Unable to decrypt cls overrides")
		return operation, 0, errors.New("unable to decrypt cls overrides")
	}
	extraConfTemplate, err := auditlog.GetExtraConfTemplate(operation.RuntimeVersion.Version)
	if err != nil {
		logger.Errorf("Unable to fetch audit log config")
		return operation, 0, errors.New("unable to fetch audit log config")
	}
	auditlogOverrideParams := auditlog.Overrides{
		Host:         auditLogHost,
		Port:         auditLogPort,
		Path:         u.Path,
		HttpPlugin:   fluentbitPlugin,
		ClsOverrides: clsOverrides,
		Config:       alo.auditLogConfig,
	}

	extraConfOverride, err := alo.renderOverrides(auditlogOverrideParams, extraConfTemplate)
	if err != nil {
		logger.Errorf("Unable to generate forward plugin to push logs: %v", err)
		return operation, time.Second, nil
	}

	operation.InputCreator.AppendOverrides("logging", []*gqlschema.ConfigEntryInput{
		{Key: "fluent-bit.conf.script", Value: replaceTenantID},
		{Key: "fluent-bit.conf.extra", Value: extraConfOverride},
		{Key: "fluent-bit.config.script", Value: replaceTenantID},
		{Key: "fluent-bit.config.extra", Value: extraConfOverride},
		{Key: "fluent-bit.externalServiceEntry.resolution", Value: "DNS"},
		{Key: "fluent-bit.externalServiceEntry.hosts", Value: fmt.Sprintf(`- %s`, auditLogHost)},
		{Key: "fluent-bit.externalServiceEntry.ports", Value: fmt.Sprintf(`- number: %s
  name: https
  protocol: TLS`, auditLogPort)},
	})
	return operation, 0, nil
}

func (alo *CLSAuditLogOverrides) renderOverrides(aloOv auditlog.Overrides, tmp *template.Template) (string, error) {
	var flOutputs bytes.Buffer
	err := tmp.Execute(&flOutputs, aloOv)
	if err != nil {
		return "", errors.Wrapf(err, "Template error while injecting cls overrides: %v", err)
	}
	return flOutputs.String(), nil
}

func (alo *CLSAuditLogOverrides) readFile(fileName string) ([]byte, error) {
	return afero.ReadFile(alo.fs, fileName)
}
