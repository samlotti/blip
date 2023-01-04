package model

type User struct {
	Name   string
	Title  string
	EMail  string
	Active bool
}

// GetUser
// Returns a specific user by 'offset' in array
func GetUser(id int) *User {
	return GetUsers()[id]

}

// GetUsers --
// Return some sample users.
func GetUsers() []*User {
	users := make([]*User, 0)
	users = append(users, &User{
		Name:   "Bob",
		Title:  "Developer",
		EMail:  "bob@test.com",
		Active: true,
	})
	users = append(users, &User{
		Name:   "Steve",
		Title:  "Manager",
		EMail:  "steve@test.com",
		Active: true,
	})
	users = append(users, &User{
		Name:   "Mike",
		Title:  "Director",
		EMail:  "mike@test.com",
		Active: true,
	})
	users = append(users, &User{
		Name:   "joe",
		Title:  "programmer",
		EMail:  "prog1l@test.com",
		Active: false,
	})
	return users
}
