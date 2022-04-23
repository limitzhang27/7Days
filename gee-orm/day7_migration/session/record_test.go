package session

import (
	"log"
	"testing"
)

var (
	user1 = &User{
		Name: "zhangjixian",
		Age:  26,
	}
	user2 = &User{
		Name: "zhengshanzhong",
		Age:  27,
	}
	user3 = &User{
		Name: "zhengtangcheng",
		Age:  27,
	}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().Model(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(user3)

	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to find record")
	}
}

func TestSession_Limit(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	err := s.Limit(1).Find(&users)
	if err != nil || len(users) != 1 {
		t.Fatal("failed to query with limit condition")
	}
}

func TestSession_Update(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "zhangjixian").Update("Age", 30)
	u := &User{}
	_ = s.OrderBy("Age desc").First(u)
	if affected != 1 || u.Age != 30 {
		log.Println(affected, u)
		t.Fatal("failed to update")
	}
}

func TestSession_DeleteAndCount(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "zhangjixian").Delete()
	count, _ := s.Count()
	if affected != 1 || count != 1 {
		log.Println(affected, count)
		t.Fatal("failed to delete or count")
	}
}
