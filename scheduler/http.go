package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"hadoop-yarn-exporter/krb5"
	"log"
	"net/http"
	"net/url"
	"os"
)

type schedulerMetrics struct {
	Scheduler scheduler `json:"scheduler"`
}
type scheduler struct {
	SchedulerInfo metrics `json:"schedulerInfo"`
}

type metrics struct {
	Capacity     float64 `json:"capacity"`
	MaxCapacity  float64 `json:"maxCapacity"`
	UsedCapacity float64 `json:"usedCapacity"`
	Queues       queues  `json:"queues"`
}

type queues struct {
	Queue []queueItem `json:"queue"`
}
type queueItem struct {
	QueueName              string  `json:"queueName"`              //队列名
	State                  string  `json:"state"`                  //队列状态
	Capacity               float64 `json:"capacity"`               //容量占比
	UsedCapacity           float64 `json:"usedCapacity"`           //容量使用量
	MaxCapacity            float64 `json:"maxCapacity"`            //最大容量占比
	AbsoluteCapacity       float64 `json:"absoluteCapacity"`       //绝对容量占比
	AbsoluteMaxCapacity    float64 `json:"absoluteMaxCapacity"`    //绝对容量最大占比
	AbsoluteUsedCapacity   float64 `json:"absoluteUsedCapacity"`   //绝对使用的容量占比
	NumApplications        int64   `json:"numApplications"`        //任务数量
	MaxApplications        int64   `json:"maxApplications"`        //最大任务数量
	MaxApplicationsPerUser int64   `json:"maxApplicationsPerUser"` //每个用户最大用户数量
	NumContainers          int64   `json:"numContainers"`          //容器任务数量
	UserLimit              int64   `json:"userLimit"`              //用户数量限制

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

	var c schedulerMetrics
	err = json.NewDecoder(resp.Body).Decode(&c)
	// fmt.Printf("%#v", c)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	return &c.Scheduler.SchedulerInfo, nil
}
