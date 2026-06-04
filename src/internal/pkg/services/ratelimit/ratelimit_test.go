package ratelimit

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew verifies the constructor creates a usable limiter.
func TestNew(t *testing.T) {
	rl := New(Config{MaxAttempts: 5, Window: time.Minute})
	require.NotNil(t, rl)
}

// TestAllow_FirstCall_ReturnsTrue verifies the very first call is always allowed.
func TestAllow_FirstCall_ReturnsTrue(t *testing.T) {
	rl := New(Config{MaxAttempts: 3, Window: time.Minute})

	assert.True(t, rl.Allow("alice"))
	assert.True(t, rl.Allow("alice"))
	assert.True(t, rl.Allow("alice"))
}

// TestAllow_ExceedsLimit_Blocks verifies Allow returns false after hitting the limit.
func TestAllow_ExceedsLimit_Blocks(t *testing.T) {
	rl := New(Config{MaxAttempts: 3, Window: time.Minute})

	assert.True(t, rl.Allow("alice"))
	assert.True(t, rl.Allow("alice"))
	assert.True(t, rl.Allow("alice"))
	assert.False(t, rl.Allow("alice"))
	assert.False(t, rl.Allow("alice"))
}

// TestAllow_WindowExpiry_Resets verifies that after the window passes, attempts reset.
func TestAllow_WindowExpiry_Resets(t *testing.T) {
	rl := New(Config{MaxAttempts: 2, Window: 50 * time.Millisecond})

	assert.True(t, rl.Allow("alice"))
	assert.True(t, rl.Allow("alice"))
	assert.False(t, rl.Allow("alice"))

	// wait for window to expire
	time.Sleep(60 * time.Millisecond)

	assert.True(t, rl.Allow("alice"), "window should have expired and reset counter")
}

// TestAllow_MultipleKeys_Independent verifies different keys have separate counters.
func TestAllow_MultipleKeys_Independent(t *testing.T) {
	rl := New(Config{MaxAttempts: 2, Window: time.Minute})

	assert.True(t, rl.Allow("alice"))
	assert.True(t, rl.Allow("bob"))
	assert.True(t, rl.Allow("alice"))
	assert.True(t, rl.Allow("bob"))

	// both are at limit now
	assert.False(t, rl.Allow("alice"))
	assert.False(t, rl.Allow("bob"))

	// carol — untouched
	assert.True(t, rl.Allow("carol"))
}

// TestAllow_SlidingWindow_OldEntriesExpire verifies that old timestamps within the
// window are pruned, so a burst followed by a pause frees up a slot.
func TestAllow_SlidingWindow_OldEntriesExpire(t *testing.T) {
	rl := New(Config{MaxAttempts: 2, Window: 80 * time.Millisecond})

	assert.True(t, rl.Allow("alice"))
	time.Sleep(90 * time.Millisecond) // first entry falls outside a 80ms window

	// after the sleep, only 1 recent attempt (this one) counts
	assert.True(t, rl.Allow("alice"))
	assert.True(t, rl.Allow("alice")) // still within limit: [recent, this]
	assert.False(t, rl.Allow("alice"))
}

// TestAllow_ZeroMaxAttempts_Disabled verifies that with MaxAttempts=0 the limiter
// always allows (disabled / no-op mode).
func TestAllow_ZeroMaxAttempts_AlwaysAllows(t *testing.T) {
	rl := New(Config{MaxAttempts: 0, Window: time.Minute})

	for i := 0; i < 100; i++ {
		assert.True(t, rl.Allow("alice"), "disabled limiter should always allow (attempt %d)", i)
	}
}

// TestRemaining_ReturnsFullLimit_ForUnknownKey verifies an unseen key returns MaxAttempts.
func TestRemaining_ReturnsFullLimit_ForUnknownKey(t *testing.T) {
	rl := New(Config{MaxAttempts: 5, Window: time.Minute})

	assert.Equal(t, 5, rl.Remaining("unknown"))
}

// TestRemaining_Decreases_AfterEachCall verifies Remaining decrements correctly.
func TestRemaining_Decreases_AfterEachCall(t *testing.T) {
	rl := New(Config{MaxAttempts: 3, Window: time.Minute})

	assert.Equal(t, 3, rl.Remaining("alice"))
	rl.Allow("alice")
	assert.Equal(t, 2, rl.Remaining("alice"))
	rl.Allow("alice")
	assert.Equal(t, 1, rl.Remaining("alice"))
	rl.Allow("alice")
	assert.Equal(t, 0, rl.Remaining("alice"))
	rl.Allow("alice") // blocked, but doesn't add a timestamp
	assert.Equal(t, 0, rl.Remaining("alice"))
}

// TestRemaining_Refills_AfterWindowExpiry verifies that after the window passes
// remaining returns to full.
func TestRemaining_Refills_AfterWindowExpiry(t *testing.T) {
	rl := New(Config{MaxAttempts: 2, Window: 50 * time.Millisecond})

	rl.Allow("alice")
	rl.Allow("alice")
	assert.Equal(t, 0, rl.Remaining("alice"))

	time.Sleep(60 * time.Millisecond)

	assert.Equal(t, 2, rl.Remaining("alice"))
}

// TestResetAfter_Zero_ForUnknownKey verifies ResetAfter returns 0 for unseen key.
func TestResetAfter_Zero_ForUnknownKey(t *testing.T) {
	rl := New(Config{MaxAttempts: 5, Window: time.Minute})
	assert.Equal(t, time.Duration(0), rl.ResetAfter("unknown"))
}

// TestResetAfter_ReturnsPositive_ForActiveKey verifies ResetAfter is > 0 for a key with attempts.
func TestResetAfter_ReturnsPositive_ForActiveKey(t *testing.T) {
	rl := New(Config{MaxAttempts: 3, Window: 500 * time.Millisecond})

	rl.Allow("alice")
	d := rl.ResetAfter("alice")
	assert.Greater(t, d, time.Duration(0))
	// it should be well under the window (we just made the call)
	assert.LessOrEqual(t, d, 500*time.Millisecond)
}

// TestResetAfter_ReturnsZero_AfterWindowExpiry verifies ResetAfter returns 0
// once the window has elapsed since the first attempt.
func TestResetAfter_ReturnsZero_AfterWindowExpiry(t *testing.T) {
	rl := New(Config{MaxAttempts: 2, Window: 50 * time.Millisecond})

	rl.Allow("alice")
	time.Sleep(60 * time.Millisecond)
	assert.Equal(t, time.Duration(0), rl.ResetAfter("alice"))
}

// TestConcurrent_Safety verifies the limiter is safe under concurrent access.
func TestConcurrent_Safety(t *testing.T) {
	rl := New(Config{MaxAttempts: 10, Window: time.Minute})
	var wg sync.WaitGroup
	n := 50

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rl.Allow("shared")
			rl.Remaining("shared")
			rl.ResetAfter("shared")
		}()
	}
	wg.Wait()

	// After 50 goroutines hitting Allow in parallel, at most 10 should have passed;
	// the point here is no race / no panic.
	remaining := rl.Remaining("shared")
	t.Logf("shared key remaining after %d concurrent calls: %d", n, remaining)
}

// TestAllow_MixedCalls_DifferentKeys verifies independent keys don't interfere.
func TestAllow_MixedCalls_DifferentKeys(t *testing.T) {
	rl := New(Config{MaxAttempts: 1, Window: time.Minute})

	assert.True(t, rl.Allow("a"))
	assert.False(t, rl.Allow("a"))
	assert.True(t, rl.Allow("b"))
	assert.False(t, rl.Allow("a"))
	assert.True(t, rl.Allow("c"))
}

// TestRemaining_NeverGoesNegative verifies Remaining never returns less than 0.
func TestRemaining_NeverGoesNegative(t *testing.T) {
	rl := New(Config{MaxAttempts: 2, Window: time.Minute})

	rl.Allow("alice")
	rl.Allow("alice")
	rl.Allow("alice") // blocked
	rl.Allow("alice") // blocked

	remaining := rl.Remaining("alice")
	assert.GreaterOrEqual(t, remaining, 0)
	assert.Equal(t, 0, remaining)
}

// TestResetAfter_WithMultipleAttempts verifies the oldest timestamp determines the window.
func TestResetAfter_WithMultipleAttempts(t *testing.T) {
	rl := New(Config{MaxAttempts: 3, Window: 200 * time.Millisecond})

	rl.Allow("alice")
	time.Sleep(50 * time.Millisecond)
	rl.Allow("alice")

	// the window should expire based on the oldest timestamp (~150ms left from first call)
	d := rl.ResetAfter("alice")
	assert.Greater(t, d, time.Duration(0))
	assert.Less(t, d, 200*time.Millisecond)
}

// TestAllow_SameTimestamp_Edge verifies that calling Allow rapidly doesn't overflow.
func TestAllow_SameTimestamp_Edge(t *testing.T) {
	rl := New(Config{MaxAttempts: 1000, Window: time.Minute})

	for i := 0; i < 1000; i++ {
		assert.True(t, rl.Allow("fast"))
	}
	assert.False(t, rl.Allow("fast"))
	assert.Equal(t, 0, rl.Remaining("fast"))
}

// TestDifferentWindows_PerKey verifies that keys with no attempts have
// no effect on keys that are rate-limited.
func TestDifferentWindows_PerKey(t *testing.T) {
	rl := New(Config{MaxAttempts: 2, Window: time.Minute})

	rl.Allow("alice")
	assert.Equal(t, 1, rl.Remaining("alice"))
	assert.Equal(t, 2, rl.Remaining("bob"))
	rl.Allow("alice")
	assert.Equal(t, 0, rl.Remaining("alice"))
	assert.Equal(t, 2, rl.Remaining("bob"))
}
