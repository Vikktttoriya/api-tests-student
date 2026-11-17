package utils

import "time"

// WaitUntil повторяет check с интервалом interval, пока он не вернёт true или не истечёт timeout.
func WaitUntil(timeout, interval time.Duration, check func() bool) bool {
	deadline := time.Now().Add(timeout)
	for {
		if check() {
			return true
		}
		if time.Now().Add(interval).After(deadline) {
			return false
		}
		time.Sleep(interval)
	}
}
