package internal

import "sync"

var GlobalMutex = sync.Mutex{}
