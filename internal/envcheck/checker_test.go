package envcheck

import (
	"testing"
)

func TestDetect(t *testing.T) {
	rep, err := Detect()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Platform: %s", rep.Platform)
	t.Logf("DockerReady: %v", rep.DockerReady)
	t.Logf("Issues: %v", rep.Issues)
	t.Logf("FixAdvice: %v", rep.FixAdvice)
}