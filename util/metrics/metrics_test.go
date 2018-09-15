package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetricsSimple(t *testing.T) {
	_, err := serverVersion.GetMetricWithLabelValues("version")
	require.NoError(t, err)
}
