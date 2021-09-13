package scheduler

import (
	"fmt"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
)

const schedulerMetricsNamespace = "yarn_schduler"
const queueMetricsNamespace = "yarn_queue"

type collector struct { //prometheus.Register()方法接收一个接口，接口需要实现Describe方法和Collect方法
	endpoint     *url.URL
	up           *prometheus.Desc
	capacity     *prometheus.Desc
	maxCapacity  *prometheus.Desc
	usedCapacity *prometheus.Desc
	queue        *prometheus.Desc
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

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.capacity
	ch <- c.maxCapacity
	ch <- c.usedCapacity
	ch <- c.queue
}

func newFuncMetric(metricName string, docString string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(schedulerMetricsNamespace, "", metricName), docString, nil, nil)
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

func NewSchedulerCollector(endpoint *url.URL) *collector {
	return &collector{
		endpoint:     endpoint,
		up:           newFuncMetric("up", "Able to contact YARN"),
		capacity:     newFuncMetric("capacity", "capacity"),
		maxCapacity:  newFuncMetric("maxCapacity", "maxCapacity"),
		usedCapacity: newFuncMetric("usedCapacity", "usedCapacity"),
		queue:        newQueuesMetric("queue", "queue"),
	}
}
