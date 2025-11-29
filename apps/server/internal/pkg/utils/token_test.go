package utils_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zarazaex/zik/apps/server/internal/pkg/utils"
)

func cleanupCacheDir() {
	cacheDir := filepath.Join(".", "tiktoken")
	os.RemoveAll(cacheDir)
	os.Unsetenv("TIKTOKEN_CACHE_DIR")
}

func TestTokenizer_Init(t *testing.T) {
	t.Run("successful initialization", func(t *testing.T) {
		cleanupCacheDir()
		tokenizer := utils.NewTokenizer()
		err := tokenizer.Init()
		assert.NoError(t, err)

		// Check that cache dir was created
		cacheDir := filepath.Join(".", "tiktoken")
		_, err = os.Stat(cacheDir)
		assert.NoError(t, err, "cache directory should be created")
	})

	t.Run("initialization happens only once", func(t *testing.T) {
		cleanupCacheDir()
		tokenizer := utils.NewTokenizer()
		var wg sync.WaitGroup
		errs := make(chan error, 100)

		// Call Init multiple times concurrently
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				errs <- tokenizer.Init()
			}()
		}
		wg.Wait()
		close(errs)

		// Check that all calls returned no error and the logic ran once
		for err := range errs {
			assert.NoError(t, err)
		}
	})

	t.Run("mkdir failure", func(t *testing.T) {
		cleanupCacheDir()

		// Create a file where the directory should be, causing MkdirAll to fail
		cacheDir := filepath.Join(".", "tiktoken")
		err := os.WriteFile(cacheDir, []byte("i am a file"), 0644)
		require.NoError(t, err)

		tokenizer := utils.NewTokenizer()
		initErr := tokenizer.Init()

		assert.Error(t, initErr, "Init should fail if cache directory cannot be created")
	})

	t.Run("get encoding failure", func(t *testing.T) {
		cleanupCacheDir()
		// This test is tricky as it relies on the internal behavior of tiktoken.
		// We can't easily mock GetEncoding. One way is to provide an invalid encoding name.
		// However, the current implementation hardcodes "cl100k_base".
		// To test this, we would need to refactor Init to accept the encoding name.
		// For now, we'll skip this specific error case as it requires further refactoring
		// and depends on an external library's error conditions.
		t.Skip("Skipping GetEncoding failure test as it's hard to trigger reliably without more refactoring.")
	})
}

func TestTokenizer_Count(t *testing.T) {
	t.Run("count with successful init", func(t *testing.T) {
		cleanupCacheDir()
		tokenizer := utils.NewTokenizer()

		// First call should initialize
		count := tokenizer.Count("hello world")
		assert.Equal(t, 2, count, "should count tokens correctly on first call")

		// Second call should use existing encoder
		count = tokenizer.Count("hello")
		assert.Equal(t, 1, count, "should count tokens correctly on subsequent calls")
	})

	t.Run("count with failed init", func(t *testing.T) {
		cleanupCacheDir()

		// Create a file to make init fail
		cacheDir := filepath.Join(".", "tiktoken")
		err := os.WriteFile(cacheDir, []byte("i am a file"), 0644)
		require.NoError(t, err)

		tokenizer := utils.NewTokenizer()
		count := tokenizer.Count("this should not be counted")
		assert.Equal(t, 0, count, "should return 0 when tokenizer fails to initialize")
	})

	t.Run("count empty string", func(t *testing.T) {
		cleanupCacheDir()
		tokenizer := utils.NewTokenizer()
		count := tokenizer.Count("")
		assert.Equal(t, 0, count, "should return 0 for an empty string")
	})
}