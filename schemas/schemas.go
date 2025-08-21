package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole define os tipos de usuários
// @enum admin,fotografo,usuario
type UserRole string

const (
	RoleAdmin     UserRole = "admin"     // Administrador do sistema
	RoleFotografo UserRole = "fotografo" // Fotógrafo profissional
	RoleUsuario   UserRole = "usuario"   // Usuário comum
)

// User representa um usuário do sistema (baseado no teste.go)
type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Nome         string    `json:"nome" gorm:"not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	Senha        string    `json:"-" gorm:"not null"`
	Tipo         UserRole  `json:"tipo" gorm:"type:varchar(20);not null;default:'usuario'"`
	Creditos     float64   `json:"creditos" gorm:"default:0"`
	DataCadastro time.Time `json:"data_cadastro" gorm:"autoCreateTime"`
	UltimoAcesso time.Time `json:"ultimo_acesso" gorm:"autoUpdateTime"`

	// Campos específicos para fotógrafos
	Biografia *string `json:"biografia,omitempty"`
	URLPerfil *string `json:"url_perfil,omitempty"`
	URLCapa   *string `json:"url_capa,omitempty"`
	Telefone  *string `json:"telefone,omitempty"`
	Website   *string `json:"website,omitempty"`
}

// UserResponse para retornos da API
type UserResponse struct {
	ID           uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Nome         string    `json:"nome" example:"João Silva"`
	Email        string    `json:"email" example:"joao@exemplo.com"`
	Tipo         UserRole  `json:"tipo" example:"usuario"`
	Creditos     float64   `json:"creditos" example:"100.50"`
	DataCadastro time.Time `json:"data_cadastro" example:"2025-08-19T14:00:00Z"`
	UltimoAcesso time.Time `json:"ultimo_acesso" example:"2025-08-19T15:00:00Z"`
	Biografia    *string   `json:"biografia,omitempty" example:"Fotógrafo profissional"`
	URLPerfil    *string   `json:"url_perfil,omitempty" example:"https://exemplo.com/perfil.jpg"`
	URLCapa      *string   `json:"url_capa,omitempty" example:"https://exemplo.com/capa.jpg"`
	Telefone     *string   `json:"telefone,omitempty" example:"(11) 99999-9999"`
	Website      *string   `json:"website,omitempty" example:"https://joaosilva.com"`
}

// Adicione ao arquivo schemas/models.go
type TwoFactorSettings struct {
	gorm.Model
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	User           User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Enabled        bool      `gorm:"default:false"`
	LastVerifiedAt time.Time
}

type TwoFactorToken struct {
	gorm.Model
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Token     string    `gorm:"size:8;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
}

// Adicione estas respostas também
type TwoFactorSettingsResponse struct {
	ID             uint      `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Enabled        bool      `json:"enabled"`
	LastVerifiedAt time.Time `json:"last_verified_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
