package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
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
	SchedulerInfo schedulerInfo `json:"schedulerInfo"`
}

type schedulerInfo struct {
	Capacity    float64 `json:"capacity"`
	MaxCapacity float64 `json:"maxCapacity"`
	// QueueName    float64 `json:"queueName"`
	UsedCapacity float64 `json:"usedCapacity"`
	// Queues       queue   `json:"queues"`
}

type queue struct {
	QueueItem []queueItem `json:"queue"`
}
type queueItem struct {
	Capacity             float64 `json:"capacity"`
	UsedCapacity         float64 `json:"usedCapacity"`
	MaxCapacity          float64 `json:"maxCapacity"`
	AbsoluteCapacity     float64 `json:"absoluteCapacity"`
	AbsoluteMaxCapacity  float64 `json:"absoluteMaxCapacity"`
	AbsoluteUsedCapacity float64 `json:"absoluteUsedCapacity"`
	NumApplications      float64 `json:"numApplications"`
	QueueName            float64 `json:"queueName"`
	State                string  `json:"state"`
}

type collector struct {
	endpoint    *url.URL
	up          *prometheus.Desc
	capacity    *prometheus.Desc
	maxCapacity *prometheus.Desc
	// queueName    *prometheus.Desc
	usedCapacity *prometheus.Desc
}

const metricsNamespace = "yarn_scheduler"

func newFuncMetric(metricName string, docString string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(metricsNamespace, "", metricName), docString, nil, nil)
}

func newCapacityMetric(metricName string, docString string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(metricsNamespace, "", metricName), docString, []string{"name"}, nil)
}

func newQueuesMetric(metricName string, docString string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(metricsNamespace, "", metricName), docString,
		[]string{"name", "capacity", "usedCapacity", "maxCapacity", "absoluteCapacity", "absoluteMaxCapacity", "absoluteUsedCapacity", "numApplications", "queueName", "state"}, nil)
}

func NewSchedulerCollector(endpoint *url.URL) *collector {
	return &collector{
		endpoint:    endpoint,
		up:          newFuncMetric("up", "Able to contact YARN"),
		capacity:    newCapacityMetric("capacity", "capacity"),
		maxCapacity: newFuncMetric("maxCapacity", "maxCapacity"),
		// queueName:    newFuncMetric("queueName", "usedCapacity"),
		usedCapacity: newFuncMetric("usedCapacity", "usedCapacity"),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.capacity
	ch <- c.maxCapacity
	// ch <- c.queueName
	ch <- c.usedCapacity
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
	// ch <- prometheus.MustNewConstMetric(c.queueName, prometheus.CounterValue, float64(metrics.QueueName))
	ch <- prometheus.MustNewConstMetric(c.usedCapacity, prometheus.CounterValue, float64(metrics.UsedCapacity))
	return
}

func fetch(u *url.URL) (*schedulerInfo, error) {
	req := http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}
	// resp,err:=krb5.GetSpnegoHttpClient().Client.Do(&req)
	var resp *http.Response
	var err error
	// resp, err = krb5.GetSpnegoHttpClient().Client.Do(&req)
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
	fmt.Printf("%#v", c)
	if err != nil {
		return nil, err
	}
	return &c.Scheduler.SchedulerInfo, nil
}
