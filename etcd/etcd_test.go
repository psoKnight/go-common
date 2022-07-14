package etcd

import (
	"fmt"
	"testing"
	"time"
)

func TestEtcd(t *testing.T) {
	// 获取etcd
	etcdCfg := &EtcdConfig{
		Endpoints:   []string{"10.117.49.69:12379", "10.117.49.69:22379", "10.117.49.69:32379"},
		DialTimeout: time.Duration(15) * time.Second,
		Username:    "",
		Password:    "",
		TTL:         5, // 续约时间为：time.Now().Add((time.Duration(karesp.TTL) * time.Second) / 3.0)
	}

	etcdClinet, err := NewEtcdClient(etcdCfg)
	if err != nil {
		t.Errorf("New etcd client err: %v.", err)
		return
	}

	go etcdClinet.WatchDiscoverService("/host") // 启动监控服务

	time.Sleep(time.Duration(3) * time.Second) // 保证下文服务正常注册

	err = etcdClinet.PutRegisterService("/host/mysql_host", "10.171.5.216:3306")
	if err != nil {
		t.Errorf("Mysql put err: %v.", err)
		return
	}
	err = etcdClinet.PutRegisterService("/host/redis_host", "10.171.5.216:6382")
	if err != nil {
		t.Errorf("Redis put err: %v.", err)
		return
	}
	err = etcdClinet.PutRegisterService("/host/mqtt_host", "10.171.5.216:1883")
	if err != nil {
		t.Errorf("Mqtt put err: %v.", err)
		return
	}

	count := 0
	for {
		count++
		service, err := etcdClinet.GetDiscoverService("/host")
		if err != nil {
			t.Errorf("Get discover service err: %+v.", err)
			time.Sleep(time.Duration(3) * time.Second)
			continue
		}
		t.Logf(fmt.Sprintf("Current list: %v.", service))

		time.Sleep(time.Duration(1) * time.Second)

		if count == 5 {
			err := etcdClinet.DeleteRegisterService("/host/mqtt_host")
			if err != nil {
				t.Errorf("Delete register service err: %v.", err)
			}
		}
		if count == 15 {
			t.Log("It's time to close service.Bye bye.")
			etcdClinet.CloseDiscoverService()
			etcdClinet.CloceRegisterService()
			return
		}
	}
}
