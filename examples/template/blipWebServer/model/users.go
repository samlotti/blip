package model

type User struct {
	Name    string
	Title   string
	EMail   string
	Active  bool
	Profile string
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
		Name:    "Bob",
		Title:   "Developer",
		EMail:   "bob@test.com",
		Active:  true,
		Profile: "https://images.generated.photos/kOx1xdZUI-VgtuysQ2QgdeBZG9w9mipqasQqK0sIJlU/rs:fit:256:256/czM6Ly9pY29uczgu/Z3Bob3Rvcy1wcm9k/LnBob3Rvcy90cmFu/c3BhcmVudF92My92/M18wNjgzOTI2LnBu/Zw.png",
	})
	users = append(users, &User{
		Name:    "Steve",
		Title:   "Manager",
		EMail:   "steve@test.com",
		Active:  true,
		Profile: "https://images.generated.photos/4SXX-TjLhF7-ywXSlJPeJQNI7ZAuQ0H_6Cm6ZX6ntmw/rs:fit:256:256/czM6Ly9pY29uczgu/Z3Bob3Rvcy1wcm9k/LnBob3Rvcy90cmFu/c3BhcmVudF92My92/M18wOTk1NjQyLnBu/Zw.png",
	})
	users = append(users, &User{
		Name:    "Sandy",
		Title:   "Director",
		EMail:   "sandy@test.com",
		Active:  true,
		Profile: "https://images.generated.photos/jjtp-fMHeqIKSa1CqcWQcCIoyvp49VcgHBbnqPLgNiA/rs:fit:256:256/czM6Ly9pY29uczgu/Z3Bob3Rvcy1wcm9k/LnBob3Rvcy92M18w/ODU5MTQ0XzAzOTY2/NTVfMDMxNzMzNy5q/cGc.jpg",
	})
	users = append(users, &User{
		Name:    "joe",
		Title:   "programmer",
		EMail:   "prog1l@test.com",
		Active:  false,
		Profile: "https://images.generated.photos/piaKB_nHX4z9xpVEGzLf206e6V37IuLdnWYzjvfTNOs/rs:fit:256:256/czM6Ly9pY29uczgu/Z3Bob3Rvcy1wcm9k/LnBob3Rvcy90cmFu/c3BhcmVudF92Mi92/Ml8wMTkzOTAyLnBu/Zw.png",
	})
	return users
}
