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

// Unmarshal unmarshals the DTR tick from a map
func (m *DTRTick) Unmarshal(body map[string]interface{}) error {
    value, ok := body["value"].(float32)
    if !ok {
        return errors.New("missing DTR value")
    }

    m.value = value

    lastUpdated, ok := body["lastUpdated"].(int64)
    if !ok {
        return errors.New("missing DTR last updated time")
    }

    m.lastUpdated = time.UnixMilli(lastUpdated)

    frozenUntil, ok := body["frozenUntil"].(int64)
    if !ok {
        return errors.New("missing DTR frozen until time")
    }

    m.frozenUntil = time.UnixMilli(frozenUntil)

    return nil
}

func (m *DTRTick) Marshal() (map[string]interface{}, error) {
    return map[string]interface{}{
        "value":       m.value,
        "lastUpdated": m.lastUpdated.UnixMilli(),
        "frozenUntil": m.frozenUntil.UnixMilli(),
    }, nil
}
