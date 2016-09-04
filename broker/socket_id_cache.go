package broker

import "errors"

type SocketIdCache struct {
	cache map[string]string
}

func NewSocketIdCache() *SocketIdCache{
	return &SocketIdCache{
		cache: make(map[string]string),
	}
}

func (s *SocketIdCache) Store(clientId, senderId string) {
	alreadyPresent := false

	for cid, _ := range s.cache {
		if cid == clientId {
			alreadyPresent = true
			break;
		}
	}

	if !(alreadyPresent){
		s.cache[clientId] = senderId
	}

	return
}

func (s *SocketIdCache) FetchAndDelete(clientId string) (error, string){
	var sid string = ""
	var err error
	found := false

	for cid, senderId := range s.cache {
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


