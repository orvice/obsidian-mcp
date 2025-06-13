package obsidianrest

import (
	"os"
	"testing"
	"time"
)

var (
	testClient *Client
	testPath   string
)

func TestMain(m *testing.M) {

	testPath = os.Getenv("OBSIDIAN_TEST_PATH")
	addr := os.Getenv("OBSIDIAN_BASE_URL")
	apiKey := os.Getenv("OBSIDIAN_API_KEY")
	testClient = NewClient(addr, apiKey, WithInsecureSkipVerify(true))
	os.Exit(m.Run())
}

func TestGetVaultFile(t *testing.T) {
	note, err := testClient.GetVaultFile(testPath)
	if err != nil {
		t.Fatalf("failed to get note path %s: %v", testPath, err)
	}
	t.Log(len(note.Content))
}

func TestUpdateVaultFile(t *testing.T) {
	err := testClient.UpdateVaultFile("test.md", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		t.Fatalf("failed to update note path %s: %v", "test.md", err)
	}
}
