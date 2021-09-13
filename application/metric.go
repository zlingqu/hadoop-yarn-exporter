package application

import (
	"fmt"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
)

const appMetricsNamespace = "yarn_app"

type collector struct { //prometheus.Register()方法接收一个接口，接口需要实现Describe方法和Collect方法
	endpoint    *url.URL
	up          *prometheus.Desc
	application *prometheus.Desc
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.application
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	up := 1.0

	appsInfo, _ := fetch(c.endpoint)
	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, up)

	if up == 0.0 {
		return
	}
	for _, appItem := range *&appsInfo.Apps.App {
		ch <- prometheus.MustNewConstMetric(c.application, prometheus.CounterValue, appItem.QueueUsagePercentage,
			appItem.Id,
			appItem.User,
			appItem.Name,
			appItem.Queue,
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

func NewApplicationCollector(endpoint *url.URL) *collector {
	return &collector{
		endpoint:    endpoint,
		up:          newFuncMetric("up", "Able to contact YARN"),
		application: newAppsMetric("application", "application"),
	}
}
