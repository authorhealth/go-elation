package elation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type resource struct {
	ID int64
}

func TestResponse_HasPrevious(t *testing.T) {
	assert := assert.New(t)

	res := &Response[*resource]{}
	assert.False(res.HasPrevious())

	res.Previous = "foo"
	assert.True(res.HasPrevious())
}

func TestResponse_HasNext(t *testing.T) {
	assert := assert.New(t)

	res := &Response[*resource]{}
	assert.False(res.HasNext())

	res.Next = "foo"
	assert.True(res.HasNext())
}

func TestResponse_PaginationNext_PaginationPrevious(t *testing.T) {
	res := &Response[*resource]{}

	t.Run("it returns defaults when there is no next or previous URL", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(&Pagination{
			Cursor: "",
			Limit:  defaultPaginationLimit,
			Offset: 0,
		}, res.PaginationNext())

		assert.Equal(&Pagination{
			Cursor: "",
			Limit:  defaultPaginationLimit,
			Offset: 0,
		}, res.PaginationPrevious())
	})

	t.Run("it returns the limit and offset from the next and previous URLs", func(t *testing.T) {
		assert := assert.New(t)
		res.Next = "https://domain.com/foo?limit=5&offset=10"

		assert.Equal(&Pagination{
			Cursor: "",
			Limit:  5,
			Offset: 10,
		}, res.PaginationNext())

		res.Previous = "https://domain.com/foo?limit=15&offset=20"

		assert.Equal(&Pagination{
			Cursor: "",
			Limit:  15,
			Offset: 20,
		}, res.PaginationPrevious())
	})

	t.Run("it returns the cursor and limit from the next and previous URLs", func(t *testing.T) {
		assert := assert.New(t)
		res.Next = "https://domain.com/foo?cursor=abc123&limit=5"

		assert.Equal(&Pagination{
			Cursor: "abc123",
			Limit:  5,
			Offset: 0,
		}, res.PaginationNext())

		res.Previous = "https://domain.com/foo?cursor=def456&limit=10"

		assert.Equal(&Pagination{
			Cursor: "def456",
			Limit:  10,
			Offset: 0,
		}, res.PaginationPrevious())
	})
}

func TestResponse_PaginationNextWithLimit(t *testing.T) {
	res := &Response[*resource]{}

	limit := 11

	t.Run("it returns the passed-in limit and the defaults when there is no next URL", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(&Pagination{
			Cursor: "",
			Limit:  limit,
			Offset: 0,
		}, res.PaginationNextWithLimit(limit))
	})

	t.Run("it returns the passed-in limit and the offset from the next URL", func(t *testing.T) {
		assert := assert.New(t)
		res.Next = "https://domain.com/foo?limit=5&offset=10"

		assert.Equal(&Pagination{
			Cursor: "",
			Limit:  limit,
			Offset: 10,
		}, res.PaginationNextWithLimit(limit))
	})

	t.Run("it returns the passed-in limit and the cursor from the next URL", func(t *testing.T) {
		assert := assert.New(t)
		res.Next = "https://domain.com/foo?cursor=abc123&limit=5"

		assert.Equal(&Pagination{
			Cursor: "abc123",
			Limit:  limit,
			Offset: 0,
		}, res.PaginationNextWithLimit(limit))
	})
}
