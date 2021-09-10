package scheduler

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

type collector struct { //prometheus.Register()方法接收一个接口，接口需要实现Describe方法和Collect方法
	endpoint     *url.URL
	up           *prometheus.Desc
	capacity     *prometheus.Desc
	maxCapacity  *prometheus.Desc
	usedCapacity *prometheus.Desc
	queue        *prometheus.Desc
}

const schedulerMetricsNamespace = "yarn_schduler"
const queueMetricsNamespace = "yarn_queue"

func newFuncMetric(metricName string, docString string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(schedulerMetricsNamespace, "", metricName), docString, nil, nil)
}

func newCapacityMetric(metricName string, docString string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(schedulerMetricsNamespace, "", metricName), docString, []string{"name"}, nil)
}

func newQueuesMetric(metricName string, docString string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(queueMetricsNamespace, "", metricName), docString,
		[]string{"queueName",
			"state",
			"capacity",
			"maxCapacity",
			"absoluteCapacity",
			"absoluteMaxCapacity",
			"absoluteUsedCapacity",
			"numApplications",
			"maxApplications",
			"maxApplicationsPerUser",
			"numContainers",
			"userLimit",
		}, nil)
}
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	up := 1.0

	metrics, _ := fetch(c.endpoint)
	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, up)

	if up == 0.0 {
		return
	}

	ch <- prometheus.MustNewConstMetric(c.capacity, prometheus.CounterValue, float64(metrics.Capacity), "haha")
	ch <- prometheus.MustNewConstMetric(c.maxCapacity, prometheus.CounterValue, float64(metrics.MaxCapacity))
	ch <- prometheus.MustNewConstMetric(c.usedCapacity, prometheus.CounterValue, float64(metrics.UsedCapacity))
	for _, queueItem := range metrics.Queues.Queue {
		ch <- prometheus.MustNewConstMetric(c.queue, prometheus.CounterValue, queueItem.UsedCapacity,
			queueItem.QueueName,
			queueItem.State,
			fmt.Sprintf("%f", queueItem.Capacity),
			fmt.Sprintf("%f", queueItem.MaxCapacity),
			fmt.Sprintf("%f", queueItem.AbsoluteCapacity),
			fmt.Sprintf("%f", queueItem.AbsoluteMaxCapacity),
			fmt.Sprintf("%f", queueItem.AbsoluteUsedCapacity),
			fmt.Sprintf("%d", queueItem.NumApplications),
			fmt.Sprintf("%d", queueItem.MaxApplications),
			fmt.Sprintf("%d", queueItem.MaxApplicationsPerUser),
			fmt.Sprintf("%d", queueItem.NumContainers),
			fmt.Sprintf("%d", queueItem.UserLimit),
		)

	}
	return
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

func NewSchedulerCollector(endpoint *url.URL) *collector {
	return &collector{
		endpoint:     endpoint,
		up:           newFuncMetric("up", "Able to contact YARN"),
		capacity:     newCapacityMetric("capacity", "capacity"),
		maxCapacity:  newFuncMetric("maxCapacity", "maxCapacity"),
		usedCapacity: newFuncMetric("usedCapacity", "usedCapacity"),
		queue:        newQueuesMetric("queue", "queue"),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.capacity
	ch <- c.maxCapacity
	ch <- c.usedCapacity
	ch <- c.queue
}
