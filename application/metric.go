package application

import (
	"fmt"
	"hadoop-yarn-exporter/model"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
)

const appMetricsNamespace = "yarn_app"

type collector struct { //prometheus.Register()方法接收一个接口，接口需要实现Describe方法和Collect方法
	clusters    []model.Cluster
	up          *prometheus.Desc
	application *prometheus.Desc
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.application
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	for _, cluster := range c.clusters {
		url, _ := url.Parse(cluster.Endpoint + "/ws/v1/cluster/apps/?state=Running")
		appsInfo, _ := fetch(url)
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
				fmt.Sprintf("%s", cluster.Code),
			)
		}
	}
	return
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
			"clusterCode",
		}, nil)
}

func NewApplicationCollector(clusters []model.Cluster) *collector {
	return &collector{
		clusters:    clusters,
		application: newAppsMetric("application", "application"),
	}
}
