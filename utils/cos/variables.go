/*
Contains all variables for cloud object storage handling
*/
package cos

import (
	"sync"
)

/*
Lock for single execution
Used for avoid concurrent writes/reads to/from the temporary memory holding
the parts read from Cloud Object Storage
*/
var writeToPipeLock sync.Mutex
