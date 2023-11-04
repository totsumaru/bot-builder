package message

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

const (
	MatchTypeComplete = "complete" // 完全一致
	MatchTypePartial  = "partial"  // 部分一致
)

// 一致条件です
type MatchType struct {
	value string
}

// 一致条件を生成します
func NewMatchType(value string) (MatchType, error) {
	d := MatchType{value: value}

	if err := d.validate(); err != nil {
		return d, errors.NewError("検証に失敗しました", err)
	}

	return d, nil
}

// 一致条件を返します
func (m MatchType) String() string {
	return m.value
}

// 一致条件が存在しているか確認します
func (m MatchType) IsEmpty() bool {
	return m.value == ""
}

// 一致条件を検証します
func (m MatchType) validate() error {
	switch m.value {
	case MatchTypeComplete, MatchTypePartial:
		break
	default:
		return errors.NewError("一致条件が不正です")
	}

	return nil
}

// 一致条件をJSONに変換します
func (m MatchType) MarshalJSON() ([]byte, error) {
	data := struct {
		MatchType string `json:"match_type"`
	}{
		MatchType: m.value,
	}

	return json.Marshal(data)
}

// JSONから一致条件を復元します
func (m *MatchType) UnmarshalJSON(b []byte) error {
	data := struct {
		MatchType string `json:"match_type"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからMatchTypeの復元に失敗しました", err)
	}

	m.value = data.MatchType

	return nil
}
