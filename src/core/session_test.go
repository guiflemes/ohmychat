package core_test

import (
	"context"
	"testing"

	"github.com/guiflemes/ohmychat/src/core"

	"github.com/stretchr/testify/assert"
)

func TestInMemorySessionRepo(t *testing.T) {
	t.Parallel()

	t.Run("creates a new session when not found", func(t *testing.T) {
		t.Parallel()

		repo := core.NewInMemorySessionRepo()
		ctx := context.Background()

		session, err := repo.GetOrCreate(ctx, "user123")
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, "user123", session.UserID)
		assert.NotNil(t, session.Memory)
		assert.IsType(t, core.IdleState{}, session.State)
	})

	t.Run("returns the same session on second call", func(t *testing.T) {
		t.Parallel()

		repo := core.NewInMemorySessionRepo()
		ctx := context.Background()

		first, err1 := repo.GetOrCreate(ctx, "user456")
		second, err2 := repo.GetOrCreate(ctx, "user456")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, first, second)
	})

	t.Run("save does not return error", func(t *testing.T) {
		t.Parallel()

		repo := core.NewInMemorySessionRepo()
		ctx := context.Background()

		session := &core.Session{
			UserID: "user789",
			State:  core.IdleState{},
			Memory: map[string]any{"key": "value"},
		}

		err := repo.Save(ctx, session)
		assert.NoError(t, err)
	})
}
