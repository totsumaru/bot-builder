package condition

import "encoding/json"

// 起動するための条件(if)です
type Condition struct {
	kind     Kind     // 条件の種類
	expected Expected // 期待する値
}

// 条件を作成します
func NewCondition(kind Kind, expected Expected) (Condition, error) {
	c := Condition{
		kind:     kind,
		expected: expected,
	}

	if err := c.validate(); err != nil {
		return c, err
	}

	return c, nil
}

// 条件の種類を返します
func (c Condition) Kind() Kind {
	return c.kind
}

// 期待する値を返します
func (c Condition) Expected() Expected {
	return c.expected
}

// 検証します
func (c Condition) validate() error {
	return nil
}

// 条件をJSONに変換します
func (c Condition) MarshalJSON() ([]byte, error) {
	data := struct {
		Kind     Kind     `json:"kind"`
		Expected Expected `json:"expected"`
	}{
		Kind:     c.kind,
		Expected: c.expected,
	}

	return json.Marshal(data)
}

// JSONから条件を復元します
func (c *Condition) UnmarshalJSON(b []byte) error {
	data := struct {
		Kind     Kind     `json:"kind"`
		Expected Expected `json:"expected"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	c.kind = data.Kind
	c.expected = data.Expected

	return nil
}
