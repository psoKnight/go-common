package mqtt

import (
	"testing"
)

// TODO 后续补充

func TestAliMqtt(t *testing.T) {
	aliMqtt := NewAliMqtt(&AliConfig{
		UpTopicPrefix:       "",
		DownTopicPrefix:     "",
		RegionID:            "",
		InstanceID:          "",
		BrokerURL:           "",
		GroupID:             "",
		AccessKey:           "",
		SecretKey:           "",
		TokenExpireInterval: 0,
	})

	username, password, err := aliMqtt.GetAliMqttTokenUsernameAndPassword("2633733c-8d4c-4c22-8893-28e2599fdcb2", nil)
	if err != nil {
		t.Errorf("Get alibaba mqtt token username and password err: %v.", err)
		return
	}
	t.Logf("Username: %s.", username)
	t.Logf("Password: %s.", password)

	clientID := "GID_Test@@@device_1"
	username2, password2 := aliMqtt.GetAliMQTTSignatureUsernameAndPassword(clientID)
	t.Logf("%s username: %s.", clientID, username2)
	t.Logf("%s password: %s.", clientID, password2)
}
