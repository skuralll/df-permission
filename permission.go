package dfpermission

import (
	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/application"
)

// 権限管理操作のための公開インターフェース
// 内部マネージャのラップ、必要なメソッドのみ公開
type PermissionManager interface {
	// 権限チェック
	HasPermission(playerID uuid.UUID, permission string) bool

	// グループ管理
	CreateGroup(name string, permissions []string) error
	DeleteGroup(name string) error

	// プレイヤーとグループの関係
	AddPlayerToGroup(playerID uuid.UUID, playerName, groupName string) error
	RemovePlayerFromGroup(playerID uuid.UUID, groupName string) error

	// 個別プレイヤー権限
	AddPlayerPermission(playerID uuid.UUID, permission string) error
	RemovePlayerPermission(playerID uuid.UUID, permission string) error

	// グループ権限管理
	AddPermissionToGroup(groupName, permission string) error
	RemovePermissionFromGroup(groupName, permission string) error

	// システム操作
	Save() error
}

// internal.Managerの具体的実装
type permissionManager struct {
	internal *application.Manager
}

// 指定された設定でPermissionManagerを作成
// 内部マネージャをラップし、安定した公開APIを提供
func NewManager(config ManagerConfig) PermissionManager {
	internalMgr := application.NewManager(config)
	return &permissionManager{
		internal: internalMgr,
	}
}

// プレイヤーが特定の権限を持つか確認
// 直接またはグループ経由で権限を持つ場合true
func (p *permissionManager) HasPermission(playerID uuid.UUID, permission string) bool {
	return p.internal.HasPermission(playerID, permission)
}

// 指定した権限で新しいグループを作成
// 既存の場合はエラー
func (p *permissionManager) CreateGroup(name string, permissions []string) error {
	err := p.internal.CreateGroup(name, permissions)
	return convertError(err)
}

// 権限グループを削除
// 存在しない場合やシステムグループの場合はエラー
func (p *permissionManager) DeleteGroup(name string) error {
	err := p.internal.DeleteGroup(name)
	return convertError(err)
}

// プレイヤーをグループに追加
// 存在しない場合はプレイヤーを作成
func (p *permissionManager) AddPlayerToGroup(playerID uuid.UUID, playerName, groupName string) error {
	err := p.internal.AddPlayerToGroup(playerID, playerName, groupName)
	return convertError(err)
}

// プレイヤーをグループから削除
// グループにいない場合はエラー
func (p *permissionManager) RemovePlayerFromGroup(playerID uuid.UUID, groupName string) error {
	err := p.internal.RemovePlayerFromGroup(playerID, groupName)
	return convertError(err)
}

// プレイヤーに直接権限を付与
// プレイヤーが存在しない場合はエラー
func (p *permissionManager) AddPlayerPermission(playerID uuid.UUID, permission string) error {
	err := p.internal.AddPlayerPermission(playerID, permission)
	return convertError(err)
}

// プレイヤーから権限を削除
// 権限を持っていない場合はエラー
func (p *permissionManager) RemovePlayerPermission(playerID uuid.UUID, permission string) error {
	err := p.internal.RemovePlayerPermission(playerID, permission)
	return convertError(err)
}

// グループに権限を追加
// グループが存在しない場合はエラー
func (p *permissionManager) AddPermissionToGroup(groupName, permission string) error {
	err := p.internal.AddPermissionToGroup(groupName, permission)
	return convertError(err)
}

// グループから権限を削除
// 権限を持っていない場合はエラー
func (p *permissionManager) RemovePermissionFromGroup(groupName, permission string) error {
	err := p.internal.RemovePermissionFromGroup(groupName, permission)
	return convertError(err)
}

// 現在の権限データをストレージに保存
func (p *permissionManager) Save() error {
	err := p.internal.Save()
	return convertError(err)
}
