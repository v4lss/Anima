// Package alert - domain service stub.
// Notification logic (v2) will be implemented here.
package alert

// Notifier is the interface that any notification channel must implement.
type Notifier interface {
	Send(a *Alert) error
}
