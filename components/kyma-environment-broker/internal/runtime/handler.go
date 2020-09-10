package runtime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/storage/dberr"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"

	"github.com/kyma-incubator/compass/components/director/pkg/pagination"

	"github.com/pkg/errors"

	"github.com/gorilla/mux"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/storage"
	"github.com/sirupsen/logrus"
)

const (
	limitParam  = "limit"
	cursorParam = "cursor"
)

type Converter interface {
	InstancesAndOperationsToDTO(internal.Instance, *internal.ProvisioningOperation, *internal.DeprovisioningOperation, *internal.UpgradeKymaOperation) dto
}

type Handler struct {
	instancesDb  storage.Instances
	operationsDb storage.Operations
	log          logrus.FieldLogger
	converter    Converter

	defaultMaxPage int
}

func NewHandler(instanceDb storage.Instances, operationDb storage.Operations, log logrus.FieldLogger, defaultMaxPage int, converter Converter) *Handler {
	return &Handler{
		instancesDb:    instanceDb,
		operationsDb:   operationDb,
		log:            log,
		converter:      converter,
		defaultMaxPage: defaultMaxPage,
	}
}

type dto struct {
	InstanceID          string `json:"instanceId"`
	RuntimeID           string `json:"runtimeId"`
	GlobalAccountID     string `json:"globalAccountId"`
	SubAccountID        string `json:"subaccountId"`
	ShootName           string `json:"shootName"`
	ProvisioningState   string `json:"provisioningState"`
	DeprovisioningState string `json:"deprovisioningState"`
	UpgradeState        string `json:"upgradeState"`
}

type RuntimesPage struct {
	Data       []dto            `json:"Data"`
	PageInfo   *pagination.Page `json:"PageInfo"`
	TotalCount int              `json:"TotalCount"`
}

func (h *Handler) AttachRoutes(router *mux.Router) {
	router.HandleFunc("/runtimes", h.getRuntimes)
}

func (h *Handler) getRuntimes(w http.ResponseWriter, req *http.Request) {
	var toReturn []dto
	limit, cursor, err := h.getParams(req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, errors.Wrap(err, "while getting query parameters"))
		return
	}

	instances, pageInfo, totalCount, err := h.instancesDb.List(limit, cursor)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while fetching instances"))
		return
	}
	for _, instance := range instances {
		pOpr, dOpr, ukOpr, err := h.getOperationsForInstance(instance)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while fetching operations for instance"))
			return
		}
		dto := h.converter.InstancesAndOperationsToDTO(instance, pOpr, dOpr, ukOpr)
		toReturn = append(toReturn, dto)
	}

	page := RuntimesPage{
		Data:       toReturn,
		PageInfo:   pageInfo,
		TotalCount: totalCount,
	}
	writeResponse(w, http.StatusOK, page)
}

func (h *Handler) getOperationsForInstance(instance internal.Instance) (*internal.ProvisioningOperation, *internal.DeprovisioningOperation, *internal.UpgradeKymaOperation, error) {
	pOpr, err := h.operationsDb.GetProvisioningOperationByInstanceID(instance.InstanceID)
	if err != nil && !dberr.IsNotFound(err) {
		return nil, nil, nil, err
	}
	dOpr, err := h.operationsDb.GetDeprovisioningOperationByInstanceID(instance.InstanceID)
	if err != nil && !dberr.IsNotFound(err) {
		return nil, nil, nil, err
	}
	ukOpr, err := h.operationsDb.GetUpgradeKymaOperationByInstanceID(instance.InstanceID)
	if err != nil && !dberr.IsNotFound(err) {
		return nil, nil, nil, err
	}
	return pOpr, dOpr, ukOpr, nil
}

func (h *Handler) getParams(req *http.Request) (int, string, error) {
	var limit int
	var cursor string
	var err error

	params := req.URL.Query()
	limitArr, ok := params[limitParam]
	if len(limitArr) > 1 {
		return 0, "", errors.New("limit has to be one parameter")
	}

	if !ok {
		limit = h.defaultMaxPage
	} else {
		limit, err = strconv.Atoi(limitArr[0])
		if err != nil {
			return 0, "", errors.New("limit has to be an integer")
		}
	}

	if limit > h.defaultMaxPage {
		return 0, "", errors.New(fmt.Sprintf("limit is bigger than maxPage(%d)", h.defaultMaxPage))
	}

	cursorArr, ok := params[cursorParam]
	if len(cursorArr) > 1 {
		return 0, "", errors.New("cursor has to be one parameter")
	}
	if !ok {
		cursor = ""
	} else {
		cursor = cursorArr[0]
	}

	return limit, cursor, nil
}

func writeResponse(w http.ResponseWriter, code int, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		logrus.Warnf("could not write response %s", string(data))
	}
}

type errObj struct {
	Error string `json:"error"`
}

func writeErrorResponse(w http.ResponseWriter, code int, err error) {
	writeResponse(w, code, errObj{Error: err.Error()})
}
