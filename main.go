package main

import (
	"flag"
	"hadoop-yarn-exporter/collector"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	isUseKerberos  = flag.String("is_use_krb", "true", "是否使用kerberos,ture|false,默认true")
	keytabFileName = flag.String("keytab_file_name", "default.keytab", "keytab文件的名字，默认是default.keytab")
	kerberosName   = flag.String("kerberos_name", "", "kerberos的用户名，必填")
	yarnUrl        = flag.String("yarn_url", "http://ccgdc-utilitynode01.i.nease.net:8088", "yarn的url地址，默认是http://ccgdc-utilitynode01.i.nease.net:8088") //老的hdp2集群
)

func main() {

	flag.Parse()
	if *isUseKerberos != "true" && *isUseKerberos != "false" {
		log.Fatal("isUseKerberos,参数错误")
	}
	os.Setenv("isUseKerberos", *isUseKerberos)

	clusterUrl, _ := url.Parse(*yarnUrl + "/ws/v1/cluster/metrics")
	c1 := collector.NewClusterCollector(clusterUrl)
	err := prometheus.Register(c1)
	if err != nil {
		log.Fatal(err)
	}

	// schedulerUrl, _ := url.Parse(*yarnUrl + "/ws/v1/cluster/scheduler")
	// c2 := scheduler.NewSchedulerCollector(schedulerUrl)
	// err := prometheus.Register(c2)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9113", nil)
}
