package krb5

import (
	"log"
	"os"

	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
	"gopkg.in/jcmturner/gokrb5.v7/spnego"
)

func GetSpnegoHttpClient() spnego.Client {
	l := log.New(os.Stderr, "GOKRB5 Client: ", log.LstdFlags)

	// 加载keytab文件
	keytabFileName := os.Getenv("keytabFileName")
	ktFromFile, err := keytab.Load(keytabFileName)
	if err != nil {
		l.Fatalf("加载keytab文件错误")
	}

	// fmt.Printf("%#v", ktFromFile.Entries[0])
	// fmt.Printf("%#v", ktFromFile.Entries[0].Principal.Components[0])

	// 加载config文件
	conf, err := config.Load("krb5.conf")
	if err != nil {
		l.Fatalf("could not load krb5.conf: %v", err)
	}

	// fmt.Printf("%#v", conf.Realms)

	//初始化client，keytab方式
	cl := client.NewClientWithKeytab(ktFromFile.Entries[0].Principal.Components[0],
		conf.Realms[0].DefaultDomain,
		ktFromFile,
		conf,
		client.Logger(l),
		client.DisablePAFXFAST(true))

	err = cl.Login()
	if err != nil { //测试客户端是否正确
		l.Fatalf("could not login client: %v", err)
	}

	spnegoCl := spnego.NewClient(cl, nil, "")
	// spnegoCl.Client()
	// 返回一个http client
	return *spnegoCl
}
