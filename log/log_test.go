package log

import (
	"context"
	"testing"

	cs "github.com/cloudtrust/common-service"
	"github.com/stretchr/testify/assert"
)

func TestExtractInfoFromContext(t *testing.T) {
	t.Run("Nil context", func(t *testing.T) {
		assert.Len(t, extractInfoFromContext(nil), 0)
	})

	var ctx = context.TODO()
	t.Run("Empty context", func(t *testing.T) {
		assert.Len(t, extractInfoFromContext(ctx), 0)
	})

	ctx = context.WithValue(ctx, cs.CtContextAccessToken, "the-access-token")

	t.Run("Added AccessToken", func(t *testing.T) {
		assert.Len(t, extractInfoFromContext(ctx), 0)
	})

	ctx = context.WithValue(ctx, cs.CtContextUserID, "the-user-id")
	t.Run("Added UserID", func(t *testing.T) {
		assert.Len(t, extractInfoFromContext(ctx), 1*2)
	})

	ctx = context.WithValue(ctx, cs.CtContextRealmID, "the-realm-id")
	t.Run("Added RealmID", func(t *testing.T) {
		assert.Len(t, extractInfoFromContext(ctx), 2*2)
	})

	ctx = context.WithValue(ctx, cs.CtContextCorrelationID, "the-correlation-id")
	t.Run("Added CorrelationID", func(t *testing.T) {
		assert.Len(t, extractInfoFromContext(ctx), 3*2)
	})
}
