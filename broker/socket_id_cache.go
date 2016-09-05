package broker

import (
	"errors"
	"sync"
)

var cache map[string]string
var mutex sync.Mutex

//type SocketIdCache struct {
//	cache map[string]string
//}
//
//func NewSocketIdCache() *SocketIdCache{
//	return &SocketIdCache{
//		cache: make(map[string]string),
//	}
//}

func init(){
	cache = make(map[string]string)
}

func SocketIdCacheStore(clientId, senderId string) {
	defer mutex.Unlock()

	mutex.Lock()
	alreadyPresent := false

	for cid, _ := range cache {
		if cid == clientId {
			alreadyPresent = true
			break;
		}
	}

	if !(alreadyPresent){
		cache[clientId] = senderId
	}

	return
}

func SocketIdCacheFetch(clientId string) (error, string){
	defer mutex.Unlock()

	mutex.Lock()
	var sid string = ""
	var err error
	found := false

	for cid, senderId := range cache {
		if cid == clientId {
			sid = senderId
			found = true
			break
		}
	}

	if !(found) {
		err = errors.New("SenderId not found")
	}

	return err, sid
}


