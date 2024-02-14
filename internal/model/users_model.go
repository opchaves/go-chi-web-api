package model

func (u *User) CanLogin() bool {
	return u.Verified
}
