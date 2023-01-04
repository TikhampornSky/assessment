//go:build unit
// +build unit

package repos

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://mwabtxlk:x4mWGDcSX0VqkVEugDsAkXesZOAazEwF@tiny.db.elephantsql.com/mwabtxlk")
	assert.Equal(t, InitDB(), nil)
}
