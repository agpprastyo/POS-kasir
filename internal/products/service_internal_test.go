package products

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCategoriesJSON(t *testing.T) {
	t.Run("NilInput", func(t *testing.T) {
		res := parseCategoriesJSON(nil)
		assert.Empty(t, res)
	})

	t.Run("StringInput", func(t *testing.T) {
		input := `[{"id": 1, "name": "Cat 1"}]`
		res := parseCategoriesJSON(input)
		assert.Len(t, res, 1)
		assert.Equal(t, int32(1), res[0].ID)
	})

	t.Run("BytesInput", func(t *testing.T) {
		input := []byte(`[{"id": 2, "name": "Cat 2"}]`)
		res := parseCategoriesJSON(input)
		assert.Len(t, res, 1)
		assert.Equal(t, int32(2), res[0].ID)
	})

	t.Run("InterfaceSliceInput", func(t *testing.T) {
		input := []interface{}{
			map[string]interface{}{"id": float64(3), "name": "Cat 3"},
		}
		res := parseCategoriesJSON(input)
		assert.Len(t, res, 1)
		// Note: json.Unmarshal into int32 might need care if source was float64 from map
		assert.Equal(t, int32(3), res[0].ID)
	})

	t.Run("MapSliceInput", func(t *testing.T) {
		input := []map[string]interface{}{
			{"id": int32(4), "name": "Cat 4"},
		}
		res := parseCategoriesJSON(input)
		assert.Len(t, res, 1)
		assert.Equal(t, int32(4), res[0].ID)
	})

	t.Run("EmptyStringInput", func(t *testing.T) {
		res := parseCategoriesJSON("")
		assert.Empty(t, res)
	})
}
