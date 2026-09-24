package access

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	Save(user *User) error
	NextID() string
}