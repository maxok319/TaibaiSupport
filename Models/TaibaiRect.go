package Models

import "encoding/json"

type TaibaiRect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (m *TaibaiRect) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *TaibaiRect) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}
