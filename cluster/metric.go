package cluster

import (
	"fmt"
	"hadoop-yarn-exporter/model"
	"log"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
)

const metricsNamespace = "yarn"

type collector struct {
	clusters              []model.Cluster
	up                    *prometheus.Desc
	applicationsSubmitted *prometheus.Desc
	applicationsCompleted *prometheus.Desc
	applicationsPending   *prometheus.Desc
	applicationsRunning   *prometheus.Desc
	applicationsFailed    *prometheus.Desc
	applicationsKilled    *prometheus.Desc
	memoryReserved        *prometheus.Desc
	memoryAvailable       *prometheus.Desc
	memoryAllocated       *prometheus.Desc
	memoryTotal           *prometheus.Desc
	virtualCoresReserved  *prometheus.Desc
	virtualCoresAvailable *prometheus.Desc
	virtualCoresAllocated *prometheus.Desc
	virtualCoresTotal     *prometheus.Desc
	containersAllocated   *prometheus.Desc
	containersReserved    *prometheus.Desc
	containersPending     *prometheus.Desc
	nodesTotal            *prometheus.Desc
	nodesLost             *prometheus.Desc
	nodesUnhealthy        *prometheus.Desc
	nodesDecommissioned   *prometheus.Desc
	nodesDecommissioning  *prometheus.Desc
	nodesRebooted         *prometheus.Desc
	nodesActive           *prometheus.Desc
	nodesShutdown         *prometheus.Desc
	scrapeFailures        *prometheus.Desc
	failureCount          int
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.applicationsSubmitted
	ch <- c.applicationsCompleted
	ch <- c.applicationsPending
	ch <- c.applicationsRunning
	ch <- c.applicationsFailed
	ch <- c.applicationsKilled
	ch <- c.memoryReserved
	ch <- c.memoryAvailable
	ch <- c.memoryAllocated
	ch <- c.memoryTotal
	ch <- c.virtualCoresReserved
	ch <- c.virtualCoresAvailable
	ch <- c.virtualCoresAllocated
	ch <- c.virtualCoresTotal
	ch <- c.containersAllocated
	ch <- c.containersReserved
	ch <- c.containersPending
	ch <- c.nodesTotal
	ch <- c.nodesLost
	ch <- c.nodesUnhealthy
	ch <- c.nodesDecommissioned
	ch <- c.nodesDecommissioning
	ch <- c.nodesRebooted
	ch <- c.nodesActive
	ch <- c.nodesShutdown
	ch <- c.scrapeFailures
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	for _, cluster := range c.clusters {
		up := 1.0
		url, _ := url.Parse(cluster.Endpoint + "/ws/v1/cluster/metrics")
		metrics, err := fetch(url)
		if err != nil {
			up = 0.0
			c.failureCount++

			log.Println("Error while collecting data from YARN: " + err.Error())
		}

		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, up, fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.scrapeFailures, prometheus.CounterValue, float64(c.failureCount), fmt.Sprintf("%s", cluster.Code))

		if up == 0.0 {
			return
		}

		ch <- prometheus.MustNewConstMetric(c.applicationsSubmitted, prometheus.CounterValue, float64(metrics.AppsSubmitted), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.applicationsCompleted, prometheus.CounterValue, float64(metrics.AppsCompleted), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.applicationsPending, prometheus.GaugeValue, float64(metrics.AppsPending), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.applicationsRunning, prometheus.GaugeValue, float64(metrics.AppsRunning), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.applicationsFailed, prometheus.CounterValue, float64(metrics.AppsFailed), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.applicationsKilled, prometheus.CounterValue, float64(metrics.AppsKilled), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.memoryReserved, prometheus.GaugeValue, float64(metrics.ReservedMB), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.memoryAvailable, prometheus.GaugeValue, float64(metrics.AvailableMB), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.memoryAllocated, prometheus.GaugeValue, float64(metrics.AllocatedMB), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.memoryTotal, prometheus.GaugeValue, float64(metrics.TotalMB), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.virtualCoresReserved, prometheus.GaugeValue, float64(metrics.ReservedVirtualCores), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.virtualCoresAvailable, prometheus.GaugeValue, float64(metrics.AvailableVirtualCores), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.virtualCoresAllocated, prometheus.GaugeValue, float64(metrics.AllocatedVirtualCores), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.virtualCoresTotal, prometheus.GaugeValue, float64(metrics.TotalVirtualCores), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.containersAllocated, prometheus.GaugeValue, float64(metrics.ContainersAllocated), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.containersReserved, prometheus.GaugeValue, float64(metrics.ContainersReserved), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.containersPending, prometheus.GaugeValue, float64(metrics.ContainersPending), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.nodesTotal, prometheus.GaugeValue, float64(metrics.TotalNodes), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.nodesLost, prometheus.GaugeValue, float64(metrics.LostNodes), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.nodesUnhealthy, prometheus.GaugeValue, float64(metrics.UnhealthyNodes), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.nodesDecommissioned, prometheus.GaugeValue, float64(metrics.DecommissionedNodes), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.nodesDecommissioning, prometheus.GaugeValue, float64(metrics.DecommissioningNodes), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.nodesRebooted, prometheus.GaugeValue, float64(metrics.RebootedNodes), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.nodesActive, prometheus.GaugeValue, float64(metrics.ActiveNodes), fmt.Sprintf("%s", cluster.Code))
		ch <- prometheus.MustNewConstMetric(c.nodesShutdown, prometheus.GaugeValue, float64(metrics.ShutdownNodes), fmt.Sprintf("%s", cluster.Code))
	}
	return
}

func newFuncMetric(metricName string, docString string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(metricsNamespace, "", metricName), docString, []string{"clusterCode"}, nil)
}

func NewClusterCollector(clusters []model.Cluster) *collector {
	return &collector{
		clusters:              clusters,
		up:                    newFuncMetric("up", "Able to contact YARN"),
		applicationsSubmitted: newFuncMetric("applications_submitted", "Total applications submitted"),
		applicationsCompleted: newFuncMetric("applications_completed", "Total applications completed"),
		applicationsPending:   newFuncMetric("applications_pending", "Applications pending"),
		applicationsRunning:   newFuncMetric("applications_running", "Applications running"),
		applicationsFailed:    newFuncMetric("applications_failed", "Total application failed"),
		applicationsKilled:    newFuncMetric("applications_killed", "Total application killed"),
		memoryReserved:        newFuncMetric("memory_reserved", "Memory reserved"),
		memoryAvailable:       newFuncMetric("memory_available", "Memory available"),
		memoryAllocated:       newFuncMetric("memory_allocated", "Memory allocated"),
		memoryTotal:           newFuncMetric("memory_total", "Total memory"),
		virtualCoresReserved:  newFuncMetric("virtual_cores_reserved", "Virtual cores reserved"),
		virtualCoresAvailable: newFuncMetric("virtual_cores_available", "Virtual cores available"),
		virtualCoresAllocated: newFuncMetric("virtual_cores_allocated", "Virtual cores allocated"),
		virtualCoresTotal:     newFuncMetric("virtual_cores_total", "Total virtual cores"),
		containersAllocated:   newFuncMetric("containers_allocated", "Containers allocated"),
		containersReserved:    newFuncMetric("containers_reserved", "Containers reserved"),
		containersPending:     newFuncMetric("containers_pending", "Containers pending"),
		nodesTotal:            newFuncMetric("nodes_total", "Nodes total"),
		nodesLost:             newFuncMetric("nodes_lost", "Nodes lost"),
		nodesUnhealthy:        newFuncMetric("nodes_unhealthy", "Nodes unhealthy"),
		nodesDecommissioned:   newFuncMetric("nodes_decommissioned", "Nodes decommissioned"),
		nodesDecommissioning:  newFuncMetric("nodes_decommissioning", "Nodes decommissioning"),
		nodesRebooted:         newFuncMetric("nodes_rebooted", "Nodes rebooted"),
		nodesActive:           newFuncMetric("nodes_active", "Nodes active"),
		nodesShutdown:         newFuncMetric("nodes_shutdown", "Nodes shutdown"),
		scrapeFailures:        newFuncMetric("scrape_failures_total", "Number of errors while scraping YARN metrics"),
	}
}
