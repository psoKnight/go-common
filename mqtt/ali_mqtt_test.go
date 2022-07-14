package mqtt

import (
	"testing"
)

// TODO 后续补充

func TestAliMQTT(t *testing.T) {
	aliMqtt, err := NewAliMQTT(&AliMQTTConfig{
		UpTopicPrefix:       "",
		DownTopicPrefix:     "",
		RegionId:            "",
		InstanceId:          "",
		BrokerUrl:           "",
		GroupId:             "",
		AccessKey:           "",
		SecretKey:           "",
		TokenExpireInterval: 0,
	})
	if err != nil {
		t.Errorf("New alibaba err: %v.", err)
		return
	}

	username, password, err := aliMqtt.GetAliMQTTTokenUsernameAndPassword("2633733c-8d4c-4c22-8893-28e2599fdcb2", nil)
	if err != nil {
		t.Errorf("Get alibaba mqtt token username and password err: %v.", err)
		return
	}
	t.Logf("Username: %s.", username)
	t.Logf("Password: %s.", password)

	clientId := "GID_Test@@@device_1"
	username2, password2 := aliMqtt.GetAliMQTTSignatureUsernameAndPassword(clientId)
	t.Logf("%s username: %s.", clientId, username2)
	t.Logf("%s password: %s.", clientId, password2)
}
