package persistence

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPersistCacheModes verifies the daemon's cache files are written with
// owner-only perms (0600 file, 0700 dir) so cached data — query history and
// agent-written notes — is not world-readable on a shared machine. It calls
// SaveNotes end-to-end so the guard at the call site is exercised, not just
// the helpers.
func TestPersistCacheModes(t *testing.T) {
	// dir is created by persistDir inside SaveNotes (MkdirAll does not change
	// the mode of an existing directory), so use a fresh subpath under the
	// already-created temp dir to exercise the 0700 creation path.
	cache := t.TempDir()
	dir := filepath.Join(cache, "repo-a")
	require.NoError(t, SaveNotes(dir, &NoteStore{Version: "t", RepoPath: dir}))

	di, err := os.Stat(dir)
	require.NoError(t, err)
	if perm := di.Mode().Perm(); perm != 0o700 {
		t.Errorf("cache dir mode = %o, want 0700", perm)
	}

	fi, err := os.Stat(filepath.Join(dir, notesFile))
	require.NoError(t, err)
	if perm := fi.Mode().Perm(); perm != 0o600 {
		t.Errorf("cache file mode = %o, want 0600", perm)
	}
}
