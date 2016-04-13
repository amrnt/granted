package granted

import "testing"

var authorize *Authorize

func init() {
	config := &Config{}
	authorize = New(config)

	authorize.AddRole(Role{
		Name: "everybody",
		CheckFunc: func(i Instance) bool {
			return true
		},
	})

	authorize.AddRole(Role{
		Name: "nobody",
		CheckFunc: func(i Instance) bool {
			return false
		},
	})
}

type User struct {
	ID    uint    `grant:"nobody"`
	Name  string  `grant:"everybody"`
	Email string  `grant:"nobody"`
	Phone string  `grant:"nobody"`
	Posts []*Post `grant:"everybody"`
}

type Post struct {
	Title  string `grant:""`
	Count  int    `grant:"nobody"`
	Author *User  `grant:"everybody"`
}

func Test1(t *testing.T) {
	u := User{
		ID:    1,
		Name:  "Amr",
		Email: "amr@email",
		Phone: "123123123",
		Posts: []*Post{
			{
				Title: "Hello world",
				Count: 1,
				Author: &User{
					ID:    199,
					Email: "Hello",
					Posts: []*Post{
						{Title: "Hello world InsiDe 1", Count: 122},
						{Title: "Hello world InsiDe 2", Count: 144, Author: &User{
							ID:    788,
							Email: "Hello from 788",
							Posts: []*Post{
								{Title: "Hello world InsiDe 1", Count: 122},
								{Title: "Hello world InsiDe 2", Count: 144},
							},
						}},
					},
				},
			},
			{
				Title: "Hello world II",
				Count: 2,
				Author: &User{
					ID:    299,
					Email: "Hello",
					Posts: []*Post{
						{Title: "Hello world II InsiDe 1", Count: 222},
						{Title: "Hello world II InsiDe 2", Count: 244},
					},
				},
			},
		},
	}

	authorize.filter(&u)

	if u.ID == 1 {
		t.Error("Error")
	}

	if u.Email == "amr@email" {
		t.Error("Error")
	}

	if u.Phone == "123123123" {
		t.Error("Error")
	}

	if u.Posts[0].Author.ID == 199 {
		t.Error("Error")
	}

	if u.Posts[0].Author.Email == "Hello" {
		t.Error("Error")
	}

	if u.Posts[0].Author.Posts[1].Author.ID == 788 {
		t.Error("Error")
	}

	if u.Posts[0].Author.Posts[1].Author.Email == "Hello from 788" {
		t.Error("Error")
	}

	if u.Posts[1].Author.ID == 299 {
		t.Error("Error")
	}

	if u.Posts[1].Author.Email == "Hello" {
		t.Error("Error")
	}
}

type Customer struct {
	ID               int         `grant:"everybody"`
	Age              int         `grant:"nobody"`
	Emails           []string    `grant:"everybody"`
	Something        interface{} `grant:"nobody"`
	PrimaryContact   Contact     `grant:"everybody"`
	SecondaryContact *Contact    `grant:"nobody"`
	Name             string      `grant:"nobody"`
}

type Contact struct {
	ID         int         `grant:"everybody"`
	CustomerID int         `grant:"nobody"`
	Emails     []string    `grant:"nobody"`
	Something  interface{} `grant:"everybody"`
	IsPrimary  bool        `grant:"nobody"`
	Email      string      `grant:"everybody"`
}

func Test2(t *testing.T) {
	x := Customer{
		ID:     123,
		Emails: []string{"amr@", "lala@"},
		PrimaryContact: Contact{
			ID:        321,
			Emails:    []string{"amr1@", "lala1@"},
			IsPrimary: false,
			Email:     "amr@email",
			Something: "{0, 1, 2}",
		},
		SecondaryContact: &Contact{
			ID:        456,
			Emails:    []string{"amr2@", "lala2@"},
			Something: "{0, 1, 2}",
		},
	}

	authorize.filter(&x)

	if x.Emails == nil {
		t.Error("Error")
	}

	if x.PrimaryContact.ID != 321 {
		t.Error("Error")
	}

	if x.PrimaryContact.Emails != nil {
		t.Error("Error")
	}

	if x.PrimaryContact.IsPrimary != false {
		t.Error("Error")
	}

	if x.PrimaryContact.Email != "amr@email" {
		t.Error("Error")
	}

	if x.SecondaryContact != nil {
		t.Error("Error")
	}
}

func TestWithAnonStruct(t *testing.T) {
	x := struct {
		ID    int    `grant:"everybody"`
		Title string `grant:"nobody"`
	}{
		1,
		"Title",
	}

	authorize.filter(&x)

	if x.ID != 1 {
		t.Error("Error")
	}

	if x.Title == "Title" {
		t.Error("Error")
	}
}

func TestWithSliceOfAnonStruct(t *testing.T) {
	var x = []struct {
		Name string `grant:"everybody"`
		Code rune   `grant:"nobody"`
		Num  int    `grant:"nobody"`
	}{
		{"a A x", 'A', 2},
		{"some_text=some_value", '=', 9},
		{"☺a", 'a', 3},
		{"a☻☺b", '☺', 4},
	}

	authorize.filter(&x)

	if x[0].Name != "a A x" {
		t.Error("Error")
	}

	if x[0].Code == 65 {
		t.Error("Error")
	}

	if x[0].Num == 2 {
		t.Error("Error")
	}
}

func TestWithMapOfStruct(t *testing.T) {
	posts := map[string]*Post{
		"amr": &Post{Title: "Hello world", Count: 999},
	}

	authorize.filter(&posts)

	if posts["amr"].Title != "Hello world" {
		t.Error("Error")
	}

	if posts["amr"].Count == 999 {
		t.Error("Error")
	}
}
