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
	assert := assert.New(t)

	res := &Response[*resource]{}

	assert.Equal(&Pagination{
		Limit:  defaultPaginationLimit,
		Offset: 0,
	}, res.PaginationNext())

	assert.Equal(&Pagination{
		Limit:  defaultPaginationLimit,
		Offset: 0,
	}, res.PaginationPrevious())

	res.Next = "https://domain.com/foo?limit=5&offset=10"

	assert.Equal(&Pagination{
		Limit:  5,
		Offset: 10,
	}, res.PaginationNext())

	res.Previous = "https://domain.com/foo?limit=15&offset=20"

	assert.Equal(&Pagination{
		Limit:  15,
		Offset: 20,
	}, res.PaginationPrevious())
}
