package monitor

const (
	STATUS_START = "start"
	STATUS_CHECK = "check"

	STATUS_CHECK_DOWN = "check.down"

	STATUS_NOTIFY_DOWN      = "notify.down"
	STATUS_NOTIFY_SKIP      = "notify.skip"
	STATUS_NOTIFY_DOWN_SENT = "notify.down.sent"

	STATUS_NOTIFY_RUNNING = "notify.running"
	STATUS_NOTIFY_UP      = "notify.up"
	STATUS_NOTIFY_UP_SENT = "notify.up.sent"
)
