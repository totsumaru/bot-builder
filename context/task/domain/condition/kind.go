package condition

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

const (
	KindCreatedMessageIs = "CreatedMessageIs" // 送られたメッセージがxだったら
	KindClickedButtonIs  = "ClickedButtonIs"  // 押されたボタンがxだったら
	KindOperatorIs       = "OperatorIs"       // 操作を実行したユーザーがxだったら
	KindOperatorRoleHas  = "OperatorRoleHas"  // 操作を実行したユーザーがxのロールを持っていたら
)

// 条件の種類です
//
// ifブロックの条件となる種別です。
type Kind struct {
	value string
}

// 条件の種類を作成します
func NewKind(value string) (Kind, error) {
	k := Kind{value: value}

	if err := k.validate(); err != nil {
		return k, errors.NewError("検証に失敗しました", err)
	}

	return k, nil
}

// 条件の種類を返します
func (k Kind) String() string {
	return k.value
}

// 条件の種類が存在しているか確認します
func (k Kind) IsEmpty() bool {
	return k.value == ""
}

// 条件の種類を検証します
func (k Kind) validate() error {
	switch k.value {
	case KindCreatedMessageIs:
	case KindClickedButtonIs:
	case KindOperatorIs:
	case KindOperatorRoleHas:
	default:
		return errors.NewError("条件の種類が不正です")
	}

	return nil
}

// 条件の種類をJSONに変換します
func (k Kind) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: k.value,
	}

	return json.Marshal(data)
}

// 条件の種類をJSONから復元します
func (k *Kind) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	k.value = data.Value

	return nil
}
