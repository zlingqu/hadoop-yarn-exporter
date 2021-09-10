package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"hadoop-yarn-exporter/krb5"

	"github.com/prometheus/client_golang/prometheus"
)

type appsMetrics struct {
	Apps apps `json:"apps"`
}

type apps struct {
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

type collector struct { //prometheus.Register()方法接收一个接口，接口需要实现Describe方法和Collect方法
	endpoint    *url.URL
	up          *prometheus.Desc
	application *prometheus.Desc
}

const appMetricsNamespace = "yarn_app"

func newFuncMetric(metricName string, docString string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(appMetricsNamespace, "", metricName), docString, nil, nil)
}

func newAppsMetric(metricName string, docString string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(appMetricsNamespace, "", metricName), docString,
		[]string{"id",
			"user",
			"name",
			"queue",
			"state",
			"finalStatus",
			"progress",
			"applicationType",
			"startedTime",
			"elapsedTime",
			"allocatedMB",
			"reservedMB",
			"allocatedVCores",
			"reservedVCores",
			"clusterUsagePercentage",
		}, nil)
}
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	up := 1.0

	metrics, _ := fetch(c.endpoint)
	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, up)

	if up == 0.0 {
		return
	}
	for _, appItem := range metrics.Apps.App {
		fmt.Println(appItem)
		ch <- prometheus.MustNewConstMetric(c.application, prometheus.CounterValue, appItem.QueueUsagePercentage,
			appItem.Id,
			appItem.User,
			appItem.Name,
			appItem.State,
			appItem.FinalStatus,
			fmt.Sprintf("%f", appItem.Progress),
			appItem.ApplicationType,
			fmt.Sprintf("%d", appItem.StartedTime),
			fmt.Sprintf("%d", appItem.ElapsedTime),
			fmt.Sprintf("%f", appItem.AllocatedMB),
			fmt.Sprintf("%f", appItem.ReservedMB),
			fmt.Sprintf("%f", appItem.AllocatedVCores),
			fmt.Sprintf("%f", appItem.ReservedVCores),
			fmt.Sprintf("%f", appItem.ClusterUsagePercentage),
		)
	}
	return
}

func fetch(u *url.URL) (*appsMetrics, error) {
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

	var c appsMetrics
	err = json.NewDecoder(resp.Body).Decode(&c)
	// fmt.Printf("%#v", c)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	fmt.Print(c)
	return &c.Apps, nil
}

func NewApplicationCollector(endpoint *url.URL) *collector {
	return &collector{
		endpoint:    endpoint,
		up:          newFuncMetric("up", "Able to contact YARN"),
		application: newAppsMetric("application", "application"),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.application
}
