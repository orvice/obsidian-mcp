package obsidianrest

import (
	"fmt"
	"os"
	"testing"
)

var (
	testClient *Client
)

func TestMain(m *testing.M) {
	addr := os.Getenv("OBSIDIAN_BASE_URL")
	apiKey := os.Getenv("OBSIDIAN_API_KEY")
	testClient = NewClient(addr, apiKey, WithInsecureSkipVerify(true))
	os.Exit(m.Run())
}

func TestGetVaultFile(t *testing.T) {
	note, err := testClient.GetVaultFile("test.md")
	if err != nil {
		t.Fatalf("failed to get note: %v", err)
	}
	fmt.Println(note)
}
