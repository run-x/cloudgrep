package util

import (
	"testing"

	"github.com/matishsiao/goInfo"
	"github.com/run-x/cloudgrep/pkg/app"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAmplitudeEvent(t *testing.T) {
	systemInfo, err := goInfo.GetInfo()
	if err != nil {
		t.Error(err)
	}
	type args struct {
		eventType       string
		eventProperties map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "TestGenerateAmplitudeEvent",
			args: args{
				eventType:       BaseEvent,
				eventProperties: nil,
			},
			want: map[string]interface{}{
				"user_id":          userId,
				"device_id":        systemInfo.Hostname,
				"event_type":       BaseEvent,
				"event_properties": map[string]interface{}{"application": application},
				"app_version":      app.Version,
				"platform":         systemInfo.Platform,
				"insert_id":        "test",
				"session_id":       sessionId,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GenerateAmplitudeEvent(tt.args.eventType, tt.args.eventProperties)
			assert.Equal(t, got["user_id"], userId)
			assert.Equal(t, got["device_id"], systemInfo.Hostname)
			assert.Equal(t, got["event_type"], tt.args.eventType)
			assert.Equal(t, got["app_version"], app.Version)
			assert.Equal(t, got["platform"], systemInfo.Platform)
			gotEventProperties := got["event_properties"].(map[string]interface{})
			assert.Equal(t, gotEventProperties["application"], application)
		})
	}
}