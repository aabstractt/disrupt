package tickable

import (
	"errors"
	"time"
)

type DTRTick struct {
	value float32

	lastUpdated time.Time
	frozenUntil time.Time
}

// Value returns the value of the DTR tick
func (m *DTRTick) Value() float32 {
	return m.value
}

// SetValue sets the value of the DTR tick
func (m *DTRTick) SetValue(value float32) {
	m.value = value
}

// UpdateRemaining updates the remaining time until the DTR tick is unfrozen
func (m *DTRTick) UpdateRemaining(seconds int64) {
	m.frozenUntil = time.Now().Add(time.Duration(seconds) * time.Second)

	m.lastUpdated = time.Now()
}

// Remaining returns the remaining time until the DTR tick is unfrozen
func (m *DTRTick) Remaining() time.Duration {
	if m.frozenUntil.UnixMilli() == 0 {
		return 0
	}

	return time.Until(m.frozenUntil)
}

func UnmarshalDTR(data map[string]interface{}) (*DTRTick, error) {
	value, ok := data["value"].(float32)
	if !ok {
		return nil, errors.New("missing DTR value")
	}

	lastUpdated, ok := data["lastUpdated"].(int64)
	if !ok {
		return nil, errors.New("missing DTR last updated time")
	}

	frozenUntil, ok := data["frozenUntil"].(int64)
	if !ok {
		return nil, errors.New("missing DTR frozen until time")
	}

	return &DTRTick{
		value:       value,
		lastUpdated: time.UnixMilli(lastUpdated),
		frozenUntil: time.UnixMilli(frozenUntil),
	}, nil
}
