package Models

import (
	"gitlab.com/ptflp/goboilerplate/internal/infrastructure/db/types"
	"time"
)

//go:generate easytags $GOFILE json,db,db_ops,db_type,db_default,db_index

type UserDTO struct {
	ID            int              `json:"id" db:"id" db_type:"BIGSERIAL primary key" db_default:"not null"`
	Name          types.NullString `json:"name" db:"name" db_type:"varchar(55)" db_default:"default null" db_ops:"create,update"`
	Phone         types.NullString `json:"phone" db:"phone" db_type:"varchar(34)" db_default:"default null" db_index:"index,unique" db_ops:"create,update"`
	Email         types.NullString `json:"email" db:"email" db_type:"varchar(89)" db_default:"default null" db_index:"index,unique" db_ops:"create,update"`
	Password      types.NullString `json:"password" db:"password" db_type:"varchar(144)" db_default:"default null" db_ops:"create,update"`
	Status        int              `json:"status" db:"status" db_type:"int" db_default:"default 0" db_ops:"create,update"`
	Role          int              `json:"role" db:"role" db_type:"int" db_default:"not null" db_ops:"create,update"`
	Verified      bool             `json:"verified" db:"verified" db_type:"boolean" db_default:"not null" db_ops:"create,update"`
	EmailVerified bool             `json:"email_verified" db:"email_verified" db_type:"boolean" db_default:"not null" db_ops:"create,update"`
	PhoneVerified bool             `json:"phone_verified" db:"phone_verified" db_type:"boolean" db_default:"not null" db_ops:"create,update"`
	CreatedAt     time.Time        `json:"created_at" db:"created_at" db_type:"timestamp" db_default:"default (now()) not null" db_index:"index"`
	UpdatedAt     time.Time        `json:"updated_at" db:"updated_at" db_type:"timestamp" db_default:"default (now()) not null" db_index:"index"`
	DeletedAt     types.NullTime   `json:"deleted_at" db:"deleted_at" db_type:"timestamp" db_default:"default null" db_index:"index"`
}

func (u *UserDTO) TableName() string {
	return "users"
}

func (u *UserDTO) OnCreate() []string {
	return []string{}
}

func (u *UserDTO) SetID(id int) *UserDTO {
	u.ID = id
	return u
}

func (u *UserDTO) GetID() int {
	return u.ID
}

func (u *UserDTO) SetName(name string) *UserDTO {
	u.Name = types.NewNullString(name)
	return u
}

func (u *UserDTO) GetName() string {
	return u.Name.String
}

func (u *UserDTO) SetPhone(phone string) *UserDTO {
	u.Phone = types.NewNullString(phone)
	return u
}

func (u *UserDTO) GetPhone() string {
	return u.Phone.String
}

func (u *UserDTO) SetEmail(email string) *UserDTO {
	u.Email = types.NewNullString(email)
	return u
}

func (u *UserDTO) GetEmail() string {
	return u.Email.String
}

func (u *UserDTO) SetPassword(password string) *UserDTO {
	u.Password = types.NewNullString(password)
	return u
}

func (u *UserDTO) GetPassword() string {
	return u.Password.String
}

func (u *UserDTO) SetStatus(status int) *UserDTO {
	u.Status = status
	return u
}

func (u *UserDTO) GetStatus() int {
	return u.Status
}

func (u *UserDTO) SetRole(role int) *UserDTO {
	u.Role = role
	return u
}

func (u *UserDTO) GetRole() int {
	return u.Role
}

func (u *UserDTO) SetVerified(verified bool) *UserDTO {
	u.Verified = verified
	return u
}

func (u *UserDTO) GetVerified() bool {
	return u.Verified
}

func (u *UserDTO) SetEmailVerified(emailVerified bool) *UserDTO {
	u.EmailVerified = emailVerified
	return u
}

func (u *UserDTO) GetEmailVerified() bool {
	return u.EmailVerified
}

func (u *UserDTO) SetPhoneVerified(phoneVerified bool) *UserDTO {
	u.PhoneVerified = phoneVerified
	return u
}

func (u *UserDTO) GetPhoneVerified() bool {
	return u.PhoneVerified
}

func (s *UserDTO) SetCreatedAt(createdAt time.Time) *UserDTO {
	s.CreatedAt = createdAt
	return s
}

func (s *UserDTO) GetCreatedAt() time.Time {
	return s.CreatedAt
}

func (s *UserDTO) SetUpdatedAt(updatedAt time.Time) *UserDTO {
	s.UpdatedAt = updatedAt
	return s
}

func (s *UserDTO) GetUpdatedAt() time.Time {
	return s.UpdatedAt
}

func (s *UserDTO) SetDeletedAt(deletedAt time.Time) *UserDTO {
	s.DeletedAt.Time.Time = deletedAt
	s.DeletedAt.Time.Valid = true
	return s
}

func (s *UserDTO) GetDeletedAt() time.Time {
	return s.DeletedAt.Time.Time
}
