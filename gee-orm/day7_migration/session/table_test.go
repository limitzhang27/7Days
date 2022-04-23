package session

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}
