package elation

import "encoding/json"

// NonNullJSONArray ensures that nil slices are marshalled to "[]" instead of "null". It also ensures that empty JSON arrays
// are marshalled to nil slices.
type NonNullJSONArray[T any] []T

func (nnja NonNullJSONArray[T]) MarshalJSON() ([]byte, error) {
	if nnja == nil {
		nnja = make([]T, 0)
	}

	return json.Marshal([]T(nnja))
}

func (nnja *NonNullJSONArray[T]) UnmarshalJSON(data []byte) error {
	type alias NonNullJSONArray[T]
	var a alias

	err := json.Unmarshal(data, &a)
	if err != nil {
		return err
	}

	if len(a) == 0 {
		return nil
	}

	*nnja = NonNullJSONArray[T](a)

	return nil
}
