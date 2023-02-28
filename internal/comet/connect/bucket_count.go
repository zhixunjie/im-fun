package connect

func (b *Bucket) IPCount() (res map[string]struct{}) {
	var ip string
	b.rwLock.RLock()
	res = make(map[string]struct{}, len(b.ipCount))
	for ip = range b.ipCount {
		res[ip] = struct{}{}
	}
	b.rwLock.RUnlock()
	return
}

// GetRoomsOnline 获取所有在线人数大于0的房间
func (b *Bucket) GetRoomsOnline() map[string]struct{} {
	var roomId string
	var room *Room

	res := make(map[string]struct{})
	b.rwLock.RLock()
	for roomId, room = range b.rooms {
		if room.Online > 0 {
			res[roomId] = struct{}{}
		}
	}
	b.rwLock.RUnlock()
	return res
}

// UpdateRoomOnline 更新每个房间的所有在线人数
func (b *Bucket) UpdateRoomOnline(roomCountMap map[string]int32) {
	var roomId string
	var room *Room

	b.rwLock.RLock()
	for roomId, room = range b.rooms {
		room.AllOnline = roomCountMap[roomId]
	}
	b.rwLock.RUnlock()
}
