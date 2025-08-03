package application

import "github.com/google/uuid"

// // 参加時のイベントハンドラ
// func (mgr *Manager) OnPlayerJoin(playerID uuid.UUID, defaultGroup string) error {
// 	// キャッシュ削除処理
// 	mgr.cache.InvalidatePlayer(playerID)
// 	return nil
// }

// 退出時のイベントハンドラ
func (mgr *Manager) OnPlayerLeave(playerID uuid.UUID) error {
	// キャッシュ削除処理
	mgr.cache.InvalidatePlayer(playerID)
	return nil
}
