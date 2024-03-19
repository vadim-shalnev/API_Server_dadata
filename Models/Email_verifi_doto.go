package Models

import "time"

//go:generate easytags $GOFILE json,db,db_ops,db_type,db_default,db_index
type EmailVerifyDTO struct {
	ID        int       `json:"id" db:"id" db_type:"BIGSERIAL primary key" db_default:"not null"`
	Email     string    `json:"email" db:"email" db_type:"varchar(89)" db_default:"not null" db_index:"index,unique" db_ops:"create,update"`
	UserID    int       `json:"user_id,omitempty" db:"user_id" db_ops:"create" db_type:"int" db_default:"not null" db_index:"index"`
	Hash      string    `json:"hash,omitempty" db:"hash" db_ops:"create" db_type:"char(36)" db_default:"not null" db_index:"index"`
	Verified  bool      `json:"verified" db:"verified" db_type:"boolean" db_default:"not null" db_ops:"create,update"`
	CreatedAt time.Time `json:"created_at" db:"created_at" db_type:"timestamp" db_default:"default (now()) not null" db_index:"index"`
}

func (e EmailVerifyDTO) TableName() string {
	return "email_verify"
}

func (e EmailVerifyDTO) OnCreate() []string {
	return nil
}

func (e *EmailVerifyDTO) SetID(id int) *EmailVerifyDTO {
	e.ID = id
	return e
}

func (e *EmailVerifyDTO) GetID() int {
	return e.ID
}

func (e *EmailVerifyDTO) SetEmail(email string) *EmailVerifyDTO {
	e.Email = email
	return e
}

func (e *EmailVerifyDTO) GetEmail() string {
	return e.Email
}

func (e *EmailVerifyDTO) SetUserID(userID int) *EmailVerifyDTO {
	e.UserID = userID
	return e
}

func (e *EmailVerifyDTO) GetUserID() int {
	return e.UserID
}

func (e *EmailVerifyDTO) SetHash(hash string) *EmailVerifyDTO {
	e.Hash = hash
	return e
}

func (e *EmailVerifyDTO) GetHash() string {
	return e.Hash
}

func (e *EmailVerifyDTO) SetVerified(verified bool) *EmailVerifyDTO {
	e.Verified = verified
	return e
}

func (e *EmailVerifyDTO) GetVerified() bool {
	return e.Verified
}

func (e *EmailVerifyDTO) SetCreatedAt(createdAt time.Time) *EmailVerifyDTO {
	e.CreatedAt = createdAt
	return e
}

func (e *EmailVerifyDTO) GetCreatedAt() time.Time {
	return e.CreatedAt
}
