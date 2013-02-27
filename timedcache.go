// TimedCache
// By Tomi Hiltunen 2013
// Contact: tomi@mitakuuluu.fi
//          http://tomi-hiltunen.com
//
// This package is a self-contained memory cache with expiration.
// This package is intended to be used for apps which can share data on single location.
// With this package you can create i.e. session repository with automatically expiring sessions.
// You can also add items that will keep on refreshing each time the item is accessed.
package timedcache


import (
    "time"
    "sync"
    "errors"
)


// Items contains the value and metadata for the stored objects
// Content - the stored object
// ExpiresAt - after this time the object is automatically removed
// KeepRefreshing - whether to refresh the expiration time when object is used
// RefreshBy - the extra duration given in a refresh
type cacheItem struct {
    Content         interface{}
    ExpiresAt       time.Time
    KeepRefreshing  bool
    RefreshBy       time.Duration
}


// TimedCache is the body for the storage cache
// items - The stored items
// mu    - Mutex
type TimedCache struct {
    items           map[string]cacheItem
    interval        time.Duration
    mu              sync.Mutex
}


// Creates a new instance of TimedCache
func New() (TimedCache) {
    newCache := TimedCache {
        items: make(map[string]cacheItem),
        interval: time.Minute,
    }
    go newCache.startInterval()
    return newCache
}


// Creates a new instance of TimedCache with a custom interval
// interval - the custom interval for expiration checking
func NewWithCustomInterval(interval time.Duration) (TimedCache) {
    newCache := TimedCache {
        items: make(map[string]cacheItem),
        interval: interval,
    }
    go newCache.startInterval()
    return newCache
}


// Add item without refreshing
func (tc *TimedCache) Put(key string, content interface{}, duration time.Duration) {
    tc.mu.Lock()
    tc.items[key] = cacheItem {
        Content: content,
        ExpiresAt: time.Now().UTC().Add(duration),
        KeepRefreshing: false,
    }
    tc.mu.Unlock()
}


// Add item with refreshing
func (tc *TimedCache) PutRefreshing(key string, content interface{}, duration time.Duration) {
    tc.mu.Lock()
    tc.items[key] = cacheItem {
        Content: content,
        ExpiresAt: time.Now().UTC().Add(duration),
        KeepRefreshing: true,
        RefreshBy: duration,
    }
    tc.mu.Unlock()
}


// Delete item
func (tc *TimedCache) Delete(key string) {
    tc.mu.Lock()
    delete(tc.items, key)
    tc.mu.Unlock()
}


// Get item
// If item is automatically refreshable, extends the expiration time
func (tc *TimedCache) Get(key string) (interface{}, error) {
    tc.mu.Lock()
    if theItem, ok := tc.items[key]; ok {
        if theItem.KeepRefreshing {
            theItem.ExpiresAt = time.Now().UTC().Add(theItem.RefreshBy)
            tc.items[key] = theItem
        }
        tc.mu.Unlock()
        return theItem.Content, nil
    }
    tc.mu.Unlock()
    return nil, errors.New("Item not found")
}


// Checks whether a key exists in the repository
func (tc *TimedCache) KeyExists(key string) (bool) {
    tc.mu.Lock()
    if _, ok := tc.items[key]; ok {
        tc.mu.Unlock()
        return true
    }
    tc.mu.Unlock()
    return false
}


// Start interval
// When triggered will activate removal of expired content
func (tc *TimedCache) startInterval() {
    for {
        time.Sleep(tc.interval)
        tc.RemoveExpired()
    }
}


// Remove expired
func (tc *TimedCache) RemoveExpired() {
    tc.mu.Lock()
    timeNow := time.Now().UTC()
    for key, curItem := range tc.items {
        if curItem.ExpiresAt.Before(timeNow) {
            delete(tc.items, key)
        }
    }
    tc.mu.Unlock()
}




