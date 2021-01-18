package antilog_test

import (
	"context"
	"testing"

	. "github.com/shamaazi/antilog"
	"github.com/stretchr/testify/require"
)

func TestAttachingAndGettingFromContext(t *testing.T) {
	expectedLogger := With("some", "fields", "to", "differentiate")
	ctx := AttachToContext(context.TODO(), expectedLogger)

	logger := FromContext(ctx)

	require.Equal(t, expectedLogger, logger)
}

func TestEmptyContextGivenDefaultLogger(t *testing.T) {
	defaultLogger := AntiLog{}
	logger := FromContext(context.TODO())

	require.Equal(t, defaultLogger, logger)
}
