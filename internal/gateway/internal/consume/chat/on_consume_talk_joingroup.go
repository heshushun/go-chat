package chat

import (
	"context"
	"encoding/json"
	"strconv"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"
)

type ConsumeTalkJoinGroup struct {
	Gid  int   `json:"group_id"`
	Type int   `json:"type"`
	Uids []int `json:"uids"`
}

// 加入群房间
func (h *Handler) onConsumeTalkJoinGroup(ctx context.Context, body []byte) {

	var data ConsumeTalkJoinGroup
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error("[ChatSubscribe] onConsumeTalkJoinGroup Unmarshal err: ", err.Error())
		return
	}

	sid := h.config.ServerId()
	for _, uid := range data.Uids {
		cids := h.clientStorage.GetUidFromClientIds(ctx, sid, socket.Session.Chat.Name(), strconv.Itoa(uid))

		for _, cid := range cids {
			opts := &cache.RoomOption{
				Channel:  socket.Session.Chat.Name(),
				RoomType: entity.RoomImGroup,
				Number:   strconv.Itoa(data.Gid),
				Sid:      h.config.ServerId(),
				Cid:      cid,
			}

			if data.Type == 2 {
				_ = h.roomStorage.Del(ctx, opts)
			} else {
				_ = h.roomStorage.Add(ctx, opts)
			}
		}
	}
}
