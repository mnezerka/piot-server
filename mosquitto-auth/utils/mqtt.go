package utils

import (
    "strings"
)

func GetMqttRootTopic(topic string) string {
    if idx := strings.IndexByte(topic, '/'); idx >= 0 {
        return topic[:idx]

    }
    return topic
}

func GetMqttTopicOrg(topic string) string {
    return GetMqttRootTopic(topic)
}
