package main

import (
	"flag"
	"hadoop-yarn-exporter/application"
	"hadoop-yarn-exporter/cluster"
	"hadoop-yarn-exporter/scheduler"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	isUseKerberos  = flag.String("isUseKerberos", "true", "是否使用kerberos,true|false")
	port           = flag.String("port", "9113", "指定监听端口")
	keytabFileName = flag.String("keytabFileName", "default.keytab", "keytab文件的名字")
	yarnUrl        = flag.String("yarnUrl", "http://ccgdc-corenode01.i.nease.net:8088", "yarn的url地址")
)

func main() {

	flag.Parse()
	if *isUseKerberos != "true" && *isUseKerberos != "false" {
		flag.Usage()
		log.Fatal("isUseKerberos,参数错误")
	}
	os.Setenv("isUseKerberos", *isUseKerberos)
	os.Setenv("keytabFileName", *keytabFileName)

	//cluser metrics
	clusterUrl, _ := url.Parse(*yarnUrl + "/ws/v1/cluster/metrics")
	c1 := cluster.NewClusterCollector(clusterUrl)
	err := prometheus.Register(c1)
	if err != nil {
		log.Fatal(err)
	}

	//scheduler metrics
	schedulerUrl, _ := url.Parse(*yarnUrl + "/ws/v1/cluster/scheduler")
	c2 := scheduler.NewSchedulerCollector(schedulerUrl)
	err = prometheus.Register(c2)
	if err != nil {
		log.Fatal(err)
	}

	// application metrics
	applicationUrl, _ := url.Parse(*yarnUrl + "/ws/v1/cluster/apps/?state=Running")
	c3 := application.NewApplicationCollector(applicationUrl)
	err = prometheus.Register(c3)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+*port, nil)
}
