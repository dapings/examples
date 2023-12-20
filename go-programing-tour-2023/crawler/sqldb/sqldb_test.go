package sqldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTable(t *testing.T) {
	sqldb, err := New(WithConnURL(ConnStrWithMySQL))
	assert.Nil(t, err)
	assert.NotNil(t, sqldb)
}
