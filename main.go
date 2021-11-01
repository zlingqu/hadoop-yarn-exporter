package main

import (
	"flag"
	"hadoop-yarn-exporter/application"
	"hadoop-yarn-exporter/cluster"
	"hadoop-yarn-exporter/model"
	"hadoop-yarn-exporter/scheduler"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	isUseKerberos  = flag.String("isUseKerberos", "true", "是否使用kerberos,true|false")
	port           = flag.String("port", "9113", "指定监听端口")
	keytabFileName = flag.String("keytabFileName", "default.keytab", "keytab文件的名字")
)

func handler() http.Handler {
	mycluster := []model.Cluster{
		{Code: "ccgdp", Endpoint: "http://ccgdc-corenode01.i.nease.net:8088"},
		{Code: "bigdata", Endpoint: "http://ccgdc-utilitynode01.i.nease.net:8088"},
		{Code: "hdplab", Endpoint: "http://ccgdc-master01-20002.i.nease.net:8088"},
	}

	// 	//cluser metrics
	c1 := cluster.NewClusterCollector(mycluster)
	err := prometheus.Register(c1)
	if err != nil {
		log.Fatal(err)
	}

	//scheduler metrics
	c2 := scheduler.NewSchedulerCollector(mycluster)
	err = prometheus.Register(c2)
	if err != nil {
		log.Fatal(err)
	}

	// application metrics
	c3 := application.NewApplicationCollector(mycluster)
	err = prometheus.Register(c3)
	if err != nil {
		log.Fatal(err)
	}

	return promhttp.Handler()
}

func main() {

	flag.Parse()
	if *isUseKerberos != "true" && *isUseKerberos != "false" {
		flag.Usage()
		log.Fatal("isUseKerberos,参数错误")
	}
	os.Setenv("isUseKerberos", *isUseKerberos)
	os.Setenv("keytabFileName", *keytabFileName)

	http.Handle("/metrics", handler())
	http.ListenAndServe(":"+*port, nil)
}
