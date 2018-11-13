package api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/nickysemenza/hyperion/util/clock"
)

func TestRunCommands(t *testing.T) {
	tests := []struct {
		jsonBody     string
		responseCode int
		numCues      int
	}{
		{`["set(par1:blue:0)"]`, http.StatusOK, 1},
		{`["foo(a"]`, http.StatusBadRequest, 0},
		{`[a`, http.StatusBadRequest, 0},
	}
	for _, tt := range tests {
		t.Run(tt.jsonBody, func(t *testing.T) {
			m := cue.InitializeMaster(clock.RealClock{})
			s := &config.Server{}
			s.Inputs.HTTP.Enabled = true
			ctx := s.InjectIntoContext(context.Background())

			w := httptest.NewRecorder()
			router := getRouter(ctx, m, true)

			var jsonStr = []byte(tt.jsonBody)

			req, _ := http.NewRequest("POST", "/commands", bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			require.Len(t, m.GetDefaultCueStack().Cues, tt.numCues)
			spew.Dump(m)
			require.Equal(t, tt.responseCode, w.Code)
		})
	}

}
