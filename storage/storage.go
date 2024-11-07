// storage/storage.go
package storage

import (
	"fealtyx-student-api/models"
	"sync"
)

var (
	Students  = make(map[int]models.Student)
	Mutex     = &sync.Mutex{}
	IDCounter = 1
)
