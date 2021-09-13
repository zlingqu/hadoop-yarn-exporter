package application

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

type appsInfo struct {
	Apps metrics `json:"apps"`
}

type metrics struct {
	App []appItem `json:"app"`
}

type appItem struct {
	Id                     string  `json:"id"`                     //任务id
	User                   string  `json:"user"`                   //所属用户
	Name                   string  `json:"name"`                   //任务名
	Queue                  string  `json:"queue"`                  //任务所属队列名
	State                  string  `json:"state"`                  //状态
	FinalStatus            string  `json:"finalStatus"`            //完成状态
	Progress               float64 `json:"progress"`               //任务进度
	ApplicationType        string  `json:"applicationType"`        //任务类型
	StartedTime            int64   `json:"startedTime"`            //任务开始时间
	ElapsedTime            int64   `json:"elapsedTime"`            //已经消耗的时间
	AllocatedMB            float64 `json:"allocatedMB"`            //分配内存大小
	ReservedMB             float64 `json:"reservedMB"`             //使用内存大小
	AllocatedVCores        float64 `json:"allocatedVCores"`        //分配的cpu
	ReservedVCores         float64 `json:"reservedVCores"`         //使用的cpu
	QueueUsagePercentage   float64 `json:"queueUsagePercentage"`   //资源使用占队列的百分比
	ClusterUsagePercentage float64 `json:"clusterUsagePercentage"` //资源使用占集群的百分比

}

func fetch(u *url.URL) (*appsInfo, error) {
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

	var c appsInfo
	err = json.NewDecoder(resp.Body).Decode(&c)
	// fmt.Printf("%#v", c)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	// fmt.Print(c.Apps)
	return &c, nil
}
