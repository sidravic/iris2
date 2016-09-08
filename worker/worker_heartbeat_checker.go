package worker

import "time"

func (worker *Worker) UpdateHeartBeatExpiry(){
	worker.HeartBeatExpiry = time.Now().Add(HEARTBEAT_INTERVAL)
}

func (worker *Worker) IsDisconnectAndReconnect() bool{
	now := time.Now()

	if worker.HeartBeatExpiry.After(now) {
		return false
	} else {
		return true
	}
}
