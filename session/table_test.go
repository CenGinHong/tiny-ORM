package session

import "testing"

func TestSession_CreateTable(t *testing.T) {
	s := NewSession().Model(&User{})
	if err := s.DropTable(); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateTable(); err != nil {
		t.Fatal(err)
	}
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}
