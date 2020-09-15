package runtime

import (
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
)

type runtimeDTO struct {
	InstanceID       string        `json:"instanceID"`
	RuntimeID        string        `json:"runtimeID"`
	GlobalAccountID  string        `json:"globalAccountID"`
	SubAccountID     string        `json:"subaccountID"`
	ShootName        string        `json:"shootName"`
	ServiceClassID   string        `json:"serviceClassID"`
	ServiceClassName string        `json:"serviceClassName"`
	ServicePlanID    string        `json:"servicePlanID"`
	ServicePlanName  string        `json:"servicePlanName"`
	SubaccountRegion string        `json:"subaccountRegion"`
	Status           runtimeStatus `json:"status"`
}

type runtimeStatus struct {
	CreatedAt      time.Time  `json:"createdAt"`
	ModifiedAt     time.Time  `json:"modifiedAt"`
	Provisioning   *operation `json:"provisioning"`
	Deprovisioning *operation `json:"deprovisioning"`
	UpgradingKyma  *operation `json:"upgradingKyma"`
}

type operation struct {
	State       string `json:"state"`
	Description string `json:"description"`
}

type RuntimesPage struct {
	Data       []runtimeDTO     `json:"data"`
	PageInfo   *pagination.Page `json:"pageInfo"`
	TotalCount int              `json:"totalCount"`
}