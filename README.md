TimedCache
==========

Go package for local memory cache with automatic expiration and item refreshing. I'm using this on AppEngine Go-environment for managing active sessions. Sessions need to be renewed each time user makes a request and pruned away once user has been inactive for too long. This I have achieved using TimedCache package.

Features:
---------
  * Local memory cache for storing objects.
    * No need for memcache.
  * Automatically prunes away expired objects.
    * Pruning interval is changable.
  * Objects can be set to renew expiration time each time they are being accessed.
  * Cache allows a mix of refreshing and non-refreshing objects.

Usage
-----
  ```go
    import "github.com/tomihiltunen/timedcache"
    
    var (
      sessionStore  timedcache.TimedCache
    )
    
    func init() {
      // Create default
      sessionStore = timedcache.New()
      
      // Create with custom interval
      sessionStore = timedcache.NewWithCustomInterval(time.Hour)
      
      // Put a refreshing item
      sessionStore.PutRefreshing("itemkey1", sessionObject, 24*time.Hour)
      
      // Put a non-refreshing item
      sessionStore.PutRefreshing("itemkey2", sessionObject2, 24*time.Hour)
      
      // Remove item
      sessionStore.Delete("itemkey1")
      
      // Get item
      // Calling Get will cause refreshing items to renew expiration time
      sessionStore.Get("itemkey2")
      
      // Check that key exists
      sessionStore.KeyExist("itemkey2")
    }
  ```
