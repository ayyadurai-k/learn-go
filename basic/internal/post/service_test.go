package post

import "testing"

// fakeRepo and fakeUsers are stand-ins for the real GORM repository and the
// user service. Because the service depends on interfaces, we can test all of
// its rules here with no database and no HTTP server.

type fakeRepo struct {
	created *Post
}

func (f *fakeRepo) Create(p *Post) error            { p.ID = 1; f.created = p; return nil }
func (f *fakeRepo) GetByID(uint) (*Post, error)     { return nil, ErrNotFound }
func (f *fakeRepo) List() ([]Post, error)           { return nil, nil }
func (f *fakeRepo) ListByUser(uint) ([]Post, error) { return nil, nil }

type fakeUsers struct {
	exists bool
}

func (f fakeUsers) Exists(uint) (bool, error) { return f.exists, nil }

func TestCreate_RejectsUnknownAuthor(t *testing.T) {
	svc := NewService(&fakeRepo{}, fakeUsers{exists: false})

	_, err := svc.Create(CreateInput{UserID: 99, Title: "hello"})
	if err != ErrUnknownUser {
		t.Fatalf("want ErrUnknownUser, got %v", err)
	}
}

func TestCreate_RequiresTitle(t *testing.T) {
	svc := NewService(&fakeRepo{}, fakeUsers{exists: true})

	if _, err := svc.Create(CreateInput{UserID: 1, Title: "   "}); err == nil {
		t.Fatal("expected an error for a blank title")
	}
}

func TestCreate_Persists(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo, fakeUsers{exists: true})

	p, err := svc.Create(CreateInput{UserID: 1, Title: "hello", Body: "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID == 0 || repo.created == nil {
		t.Fatal("expected the post to be saved")
	}
}
