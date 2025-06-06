package staple

import (
	"fmt"
	"log/slog"
)

// Event defines the interface implemented by all events Staple will
// generate.
//
// Map will return a map with the details of the event serialized into a
// map[string]interface{}, with only the values suitable for serialization.
type Event interface {
	fmt.Stringer
	slog.Leveler
	Type() EventType
	Map() map[string]interface{}
}

type (
	EventType int

	EventHook func(Event)

	// SprintFunc formats an arbitrary Go value into a string.
	// It is used by the supervisor to format the value of a call
	// to recover().
	SprintFunc func(interface{}) string
)

const (
	EventTypeStopTimeout EventType = iota
	EventTypeServicePanic
	EventTypeServiceTerminate
	EventTypeBackoff
	EventTypeResume
)

type EventStopTimeout struct {
	Supervisor     *Supervisor `json:"-"`
	SupervisorName string      `json:"supervisor_name"`
	Service        Service     `json:"-"`
	ServiceName    string      `json:"service"`
}

func (EventStopTimeout) Type() EventType {
	return EventTypeStopTimeout
}

func (e EventStopTimeout) String() string {
	return fmt.Sprintf(
		"%s: Service %s failed to terminate in a timely manner",
		e.Supervisor,
		e.ServiceName,
	)
}

func (e EventStopTimeout) Map() map[string]interface{} {
	return map[string]interface{}{
		"supervisor_name": e.SupervisorName,
		"service_name":    e.ServiceName,
	}
}

func (EventStopTimeout) Level() slog.Level {
	return slog.LevelWarn
}

type EventServicePanic struct {
	Supervisor       *Supervisor `json:"-"`
	SupervisorName   string      `json:"supervisor_name"`
	Service          Service     `json:"-"`
	ServiceName      string      `json:"service_name"`
	CurrentFailures  float64     `json:"current_failures"`
	FailureThreshold float64     `json:"failure_threshold"`
	Restarting       bool        `json:"restarting"`
	PanicMsg         string      `json:"panic_msg"`
	Stacktrace       string      `json:"stacktrace"`
}

func (EventServicePanic) Type() EventType {
	return EventTypeServicePanic
}

func (e EventServicePanic) String() string {
	return fmt.Sprintf(
		"%s, panic: %s, stacktrace: %s",
		serviceFailureString(
			e.SupervisorName,
			e.ServiceName,
			e.CurrentFailures,
			e.FailureThreshold,
			e.Restarting,
		),
		e.PanicMsg,
		e.Stacktrace,
	)
}

func (e EventServicePanic) Map() map[string]interface{} {
	return map[string]interface{}{
		"supervisor_name":   e.SupervisorName,
		"service_name":      e.ServiceName,
		"current_failures":  e.CurrentFailures,
		"failure_threshold": e.FailureThreshold,
		"restarting":        e.Restarting,
		"panic_msg":         e.PanicMsg,
		"stacktrace":        e.Stacktrace,
	}
}

func (EventServicePanic) Level() slog.Level {
	return slog.LevelError
}

type EventServiceTerminate struct {
	Supervisor       *Supervisor `json:"-"`
	SupervisorName   string      `json:"supervisor_name"`
	Service          Service     `json:"-"`
	ServiceName      string      `json:"service_name"`
	CurrentFailures  float64     `json:"current_failures"`
	FailureThreshold float64     `json:"failure_threshold"`
	Restarting       bool        `json:"restarting"`
	Err              interface{} `json:"error_msg"`
}

func (EventServiceTerminate) Type() EventType {
	return EventTypeServiceTerminate
}

func (e EventServiceTerminate) String() string {
	return fmt.Sprintf(
		"%s, error: %s",
		serviceFailureString(e.SupervisorName, e.ServiceName, e.CurrentFailures, e.FailureThreshold, e.Restarting),
		e.Err)
}

func (e EventServiceTerminate) Map() map[string]interface{} {
	return map[string]interface{}{
		"supervisor_name":   e.SupervisorName,
		"service_name":      e.ServiceName,
		"current_failures":  e.CurrentFailures,
		"failure_threshold": e.FailureThreshold,
		"restarting":        e.Restarting,
		"error":             e.Err,
	}
}

func (EventServiceTerminate) Level() slog.Level {
	return slog.LevelError
}

func serviceFailureString(supervisor, service string, currentFailures, failureThreshold float64, restarting bool) string {
	return fmt.Sprintf(
		"%s: Failed service '%s' (%f failures of %f), restarting: %#v",
		supervisor,
		service,
		currentFailures,
		failureThreshold,
		restarting,
	)
}

type EventBackoff struct {
	Supervisor     *Supervisor `json:"-"`
	SupervisorName string      `json:"supervisor_name"`
}

func (EventBackoff) Type() EventType {
	return EventTypeBackoff
}

func (e EventBackoff) String() string {
	return fmt.Sprintf("%s: Entering the backoff state.", e.Supervisor)
}

func (e EventBackoff) Map() map[string]interface{} {
	return map[string]interface{}{
		"supervisor_name": e.SupervisorName,
	}
}

func (EventBackoff) Level() slog.Level {
	return slog.LevelInfo
}

type EventResume struct {
	Supervisor     *Supervisor `json:"-"`
	SupervisorName string      `json:"supervisor_name"`
}

func (EventResume) Type() EventType {
	return EventTypeResume
}

func (e EventResume) String() string {
	return fmt.Sprintf("%s: Exiting backoff state.", e.Supervisor)
}

func (e EventResume) Map() map[string]interface{} {
	return map[string]interface{}{
		"supervisor_name": e.SupervisorName,
	}
}

func (EventResume) Level() slog.Level {
	return slog.LevelInfo
}
