//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package types

import (
	"encoding/json"
	"time"
)

type EventList []Event

type Event struct {
	// Mark event as initial
	Initial bool `json:"initial"`
	// Mark event as initial
	Ticker bool `json:"ticker"`
	// Event meta
	Meta NodeMeta `json:"meta"`
	// Event meta
	State NodeState `json:"state"`
	// Activity created time
	Pods []*Pod `json:"pods"`
	// Event created time
	Timestamp time.Time `json:"timestamp"`
}

type PodEvent struct {
	// Event type
	Event string
	// Event pod
	Pod *Pod
}

type ContainerEvent struct {
	// Event type
	Event string
	// Activity container
	Container *Container
}

type HostEvent struct {
	Event
}

func (s *Event) ToJson() ([]byte, error) {
	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *EventList) ToJson() ([]byte, error) {

	if s == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
