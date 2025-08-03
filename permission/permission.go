package permission

import (
	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/application"
)

// 権限管理操作のための公開インターフェース
// 内部マネージャのラップ、必要なメソッドのみ公開
type PermissionManager interface {
	// 権限チェック
	HasPermission(playerID uuid.UUID, permission string) bool

	// プレイヤー管理
	CreatePlayer(playerID uuid.UUID, playerName string) error
	UpdatePlayerName(playerID uuid.UUID, newName string) error

	// グループ管理
	CreateGroup(name string, permissions []string) error
	DeleteGroup(name string) error
	UpdateGroup(name string, permissions []string) error

	// プレイヤーとグループの関係
	AddPlayerToGroup(playerID uuid.UUID, playerName, groupName string) error
	RemovePlayerFromGroup(playerID uuid.UUID, groupName string) error

	// 個別プレイヤー権限
	AddPlayerPermission(playerID uuid.UUID, permission string) error
	RemovePlayerPermission(playerID uuid.UUID, permission string) error
	SetPlayerPermissions(playerID uuid.UUID, permissions []string) error
	GetPlayerPermissions(playerID uuid.UUID) []string

	// グループ権限管理
	AddPermissionToGroup(groupName, permission string) error
	RemovePermissionFromGroup(groupName, permission string) error

	// システム操作
	Save() error

	// イベントハンドリング
	OnPlayerJoin(playerID uuid.UUID, playerName string) error
	OnPlayerLeave(playerID uuid.UUID, playerName string) error
}

// internal.Managerの具体的実装
type permissionManager struct {
	internal *application.Manager
}

// オプションパターンでPermissionManagerを作成
// デフォルト設定をベースに、指定されたオプションを適用
func NewManager(opts ...Option) PermissionManager {
	config := buildConfig(opts...)
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

// 新しいプレイヤーを作成
// 既に存在する場合はエラー
func (p *permissionManager) CreatePlayer(playerID uuid.UUID, playerName string) error {
	return p.internal.CreatePlayer(playerID, playerName)
}

func (p *permissionManager) UpdatePlayerName(playerID uuid.UUID, newName string) error {
	return p.internal.UpdatePlayerName(playerID, newName)
}

// 指定した権限で新しいグループを作成
// 既存の場合はエラー
func (p *permissionManager) CreateGroup(name string, permissions []string) error {
	return p.internal.CreateGroup(name, permissions)
}

// 権限グループを削除
// 存在しない場合やシステムグループの場合はエラー
func (p *permissionManager) DeleteGroup(name string) error {
	return p.internal.DeleteGroup(name)
}

// グループの権限を更新（置き換え）
// グループが存在しない場合はエラー
func (p *permissionManager) UpdateGroup(name string, permissions []string) error {
	return p.internal.UpdateGroup(name, permissions)
}

// プレイヤーをグループに追加
// 存在しない場合はプレイヤーを作成
func (p *permissionManager) AddPlayerToGroup(playerID uuid.UUID, playerName, groupName string) error {
	return p.internal.AddPlayerToGroup(playerID, playerName, groupName)
}

// プレイヤーをグループから削除
// グループにいない場合はエラー
func (p *permissionManager) RemovePlayerFromGroup(playerID uuid.UUID, groupName string) error {
	return p.internal.RemovePlayerFromGroup(playerID, groupName)
}

// プレイヤーに直接権限を付与
// プレイヤーが存在しない場合はエラー
func (p *permissionManager) AddPlayerPermission(playerID uuid.UUID, permission string) error {
	return p.internal.AddPlayerPermission(playerID, permission)
}

// プレイヤーから権限を削除
// 権限を持っていない場合はエラー
func (p *permissionManager) RemovePlayerPermission(playerID uuid.UUID, permission string) error {
	return p.internal.RemovePlayerPermission(playerID, permission)
}

// グループに権限を追加
// グループが存在しない場合はエラー
func (p *permissionManager) AddPermissionToGroup(groupName, permission string) error {
	return p.internal.AddPermissionToGroup(groupName, permission)
}

// グループから権限を削除
// 権限を持っていない場合はエラー
func (p *permissionManager) RemovePermissionFromGroup(groupName, permission string) error {
	err := p.internal.RemovePermissionFromGroup(groupName, permission)
	return err
}

// プレイヤーの全権限を設定（置き換え）
// プレイヤーが存在しない場合はエラー
func (p *permissionManager) SetPlayerPermissions(playerID uuid.UUID, permissions []string) error {
	return p.internal.SetPlayerPermissions(playerID, permissions)
}

// プレイヤーの全権限を取得（個人権限+グループ権限）
// プレイヤーが存在しない場合は空のスライスを返す
func (p *permissionManager) GetPlayerPermissions(playerID uuid.UUID) []string {
	return p.internal.GetPlayerPermissions(playerID)
}

// 現在の権限データをストレージに保存
func (p *permissionManager) Save() error {
	return p.internal.Save()
}

// プレイヤー参加時の処理
func (p *permissionManager) OnPlayerJoin(playerID uuid.UUID, playerName string) error {
	return p.internal.OnPlayerJoin(playerID, playerName)
}

// プレイヤー退出時の処理
func (p *permissionManager) OnPlayerLeave(playerID uuid.UUID, playerName string) error {
	return p.internal.OnPlayerLeave(playerID, playerName)
}
