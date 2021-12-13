/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mq

type MQEnvConfig struct {
	QueueManager   string `envconfig:"QUEUE_MANAGER" default:"QM1"`
	ChannelName    string `envconfig:"CHANNEL_NAME" default:"DEV.APP.SVRCONN"`
	ConnectionName string `envconfig:"CONNECTION_NAME" default:"localhost(1414)"`
	User           string `envconfig:"USER" default:"app"`
	Password       string `envconfig:"PASSWORD" default:"password"`
	QueueName      string `envconfig:"QUEUE_NAME" default:"DEV.QUEUE.1"`
}

func (e *MQEnvConfig) Config() *ConnConfig {
	return &ConnConfig{
		ChannelName:    e.ChannelName,
		ConnectionName: e.ConnectionName,
		User:           e.User,
		Password:       e.Password,
		QueueManager:   e.QueueManager,
		QueueName:      e.QueueName,
	}
}
