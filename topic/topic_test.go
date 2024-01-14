package topic

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/alainrk/flemq/config"
)

func TestRestoreDefaultTopics(t *testing.T) {
	testFolder := fmt.Sprintf("/tmp/flemq_test_%d", rand.Int())
	defer os.RemoveAll(testFolder)

	c := config.Config{
		Store: config.StoreConfig{
			Type:   config.StoreTypeFqueue,
			Folder: testFolder,
		},
	}

	tnames := []string{"test1", "test2", "test3"}
	for _, tn := range tnames {
		New(tn, c.Store)
	}

	topics := RestoreDefaultTopics(c.Store)
	if len(topics) != len(tnames) {
		t.Fatalf("Expected %d topics, got %d", len(tnames), len(topics))
	}

	for _, tn := range tnames {
		if _, ok := topics[tn]; !ok {
			t.Fatalf("Expected topic %s to exist", tn)
		}
	}
}
