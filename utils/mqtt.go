package utils

import (
    "strings"
)

func GetMqttTopicOrg(topic string) string {
    if idx := strings.IndexByte(topic, '/'); idx >= 0 {
        return topic[:idx]

    }
    return topic
}
