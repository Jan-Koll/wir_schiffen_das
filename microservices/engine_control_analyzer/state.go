package main

import "errors"

type State int

const (
	Queued State = iota
	Running
	Failed
	Ready
	Canceled
)

func (state State) String() string {
	switch state {
	case Queued:
		return "queued"
	case Running:
		return "running"
	case Failed:
		return "failed"
	case Ready:
		return "ready"
	case Canceled:
		return "canceled"
	}
	return "invalid"
}

func (state State) MarshalJSON() ([]byte, error) {
	switch state {
	case Queued:
		return []byte(`"queued"`), nil
	case Running:
		return []byte(`"running"`), nil
	case Failed:
		return []byte(`"failed"`), nil
	case Ready:
		return []byte(`"ready"`), nil
	case Canceled:
		return []byte(`"canceled"`), nil
	}
	return nil, errors.New("invalid state")
}
