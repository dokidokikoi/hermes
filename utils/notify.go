package utils

func ConcatEvent(topic, event string) string {
	return topic + ":" + event
}
