package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"hadoop-yarn-exporter/krb5"
	"net/http"
	"net/url"
	"os"
)

type clusterMetrics struct {
	ClusterMetrics metrics `json:"clusterMetrics"`
}

type metrics struct {
	AppsSubmitted         int `json:"appsSubmitted"`
	AppsCompleted         int `json:"appsCompleted"`
	AppsPending           int `json:"appsPending"`
	AppsRunning           int `json:"appsRunning"`
	AppsFailed            int `json:"appsFailed"`
	AppsKilled            int `json:"appsKilled"`
	ReservedMB            int `json:"reservedMB"`
	AvailableMB           int `json:"availableMB"`
	AllocatedMB           int `json:"allocatedMB"`
	ReservedVirtualCores  int `json:"reservedVirtualCores"`
	AvailableVirtualCores int `json:"availableVirtualCores"`
	AllocatedVirtualCores int `json:"allocatedVirtualCores"`
	ContainersAllocated   int `json:"containersAllocated"`
	ContainersReserved    int `json:"containersReserved"`
	ContainersPending     int `json:"containersPending"`
	TotalMB               int `json:"totalMB"`
	TotalVirtualCores     int `json:"totalVirtualCores"`
	TotalNodes            int `json:"totalNodes"`
	LostNodes             int `json:"lostNodes"`
	UnhealthyNodes        int `json:"unhealthyNodes"`
	DecommissioningNodes  int `json:"decommissioningNodes"`
	DecommissionedNodes   int `json:"decommissionedNodes"`
	RebootedNodes         int `json:"rebootedNodes"`
	ActiveNodes           int `json:"activeNodes"`
	ShutdownNodes         int `json:"shutdownNodes"`
}

func fetch(u *url.URL) (*metrics, error) {
	req := http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}
	var resp *http.Response
	var err error
	isUseKerberos := os.Getenv("isUseKerberos")
	if isUseKerberos == "true" {
		spnegoClient := krb5.GetSpnegoHttpClient()
		resp, err = spnegoClient.Do(&req)
		// fmt.Printf("Use krb")
	} else {
		// fmt.Printf("No use krb")
		resp, err = http.DefaultClient.Do(&req)
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected HTTP status: %v", resp.StatusCode))
	}

	var c clusterMetrics
	err = json.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("%#v", c.ClusterMetrics)
	return &c.ClusterMetrics, nil
}
