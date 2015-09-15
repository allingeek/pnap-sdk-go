package pnap

import (
	"errors"
	"fmt"
	"time"
)

const documentType = "application/vnd.pncp.v.1.0+json"

// Interfaces

type API interface {
	CreateVirtualMachine(props MachineProperties) (r Future, err error)
}

// API implementation

type PNAP struct {
	Endpoint       string
	AccountID      string
	ApplicationKey string
	SharedSecret   string
	NodeID         string
	Debug          bool
	Backoff        time.Duration
}

func NewPNAP() PNAP {
	return PNAP{}
}

func (r *PNAP) CreateVirtualMachine(props MachineProperties) (res Future, err error) {
	path := fmt.Sprintf(`/account/%s/node/%s/device/virtualmachine`, r.AccountID, r.NodeID)
	var (
		emsg     string
		retry    bool
		attempts int
	)
	for retry = true; retry && attempts < 5; attempts = Backoff(attempts) {
		res, emsg, retry, _ = r.call(`POST`, path, ``, props)
	}

	return res, errors.New(emsg)
}

func Backoff(i int) int {
	return i + 1
}

//
// API Request/Response Structures
//

type OSTemplate struct {
	ResourceURL string `json:"resourceURL"`
}

type Task struct {
	PercentageComplete    int
	RequestStateEnum      string
	ProcessDescription    string
	LatestTaskDescription string
	Result                interface{}
	ErrorCode             string
	ErrorMessage          string
	LastUpdatedTimestamp  string
	CreatedTimestamp      string
}

type MachineProperties struct {
	Name                    string     `json:"name"`
	Description             string     `json:"description"`
	StorageInGB             uint16     `json:"storageGB"`
	MemoryInMB              uint32     `json:"memoryMB"`
	VCpuCount               uint8      `json:"vCPUs"`
	StorageType             string     `json:"storageType"`
	PowerStatus             string     `json:"powerStatus"`
	OperatingSystemTemplate OSTemplate `json:"operatingSystemTemplate"`
	ImageResource           string     `json:"imageResource"`
	Password                string     `json:"newOperatingSystemAdminPassword"`
}
