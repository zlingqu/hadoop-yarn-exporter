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

type scheduler struct {
	SchedulerInfo metric `json:"schedulerInfo"`
}

type metric struct {
	Capacity    float64 `json:"capacity"`
	MaxCapacity float64 `json:"maxCapacity"`
	// QueueName    float64 `json:"queueName"`
	UsedCapacity float64 `json:"usedCapacity"`
}

type collector struct {
	endpoint    *url.URL
	up          *prometheus.Desc
	capacity    *prometheus.Desc
	maxCapacity *prometheus.Desc
	// queueName    *prometheus.Desc
	usedCapacity *prometheus.Desc
}

const metricsNamespace = "yarn"

func newFuncMetric(metricName string, docString string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(metricsNamespace, "", metricName), docString, nil, nil)
}

func NewSchedulerCollector(endpoint *url.URL) *collector {
	return &collector{
		endpoint:    endpoint,
		up:          newFuncMetric("up", "Able to contact YARN"),
		capacity:    newFuncMetric("capacity", "capacity"),
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

	ch <- prometheus.MustNewConstMetric(c.capacity, prometheus.CounterValue, float64(metrics.Capacity))
	ch <- prometheus.MustNewConstMetric(c.maxCapacity, prometheus.CounterValue, float64(metrics.MaxCapacity))
	// ch <- prometheus.MustNewConstMetric(c.queueName, prometheus.CounterValue, float64(metrics.QueueName))
	ch <- prometheus.MustNewConstMetric(c.usedCapacity, prometheus.CounterValue, float64(metrics.UsedCapacity))
	return
}

func fetch(u *url.URL) (*metric, error) {
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
		resp, err = krb5.GetSpnegoHttpClient().Client.Do(&req)
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

	var c scheduler
	err = json.NewDecoder(resp.Body).Decode(&c)
	fmt.Printf("%#v", c)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v", c.SchedulerInfo)
	return &c.SchedulerInfo, nil
}
