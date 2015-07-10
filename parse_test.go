package parse_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tmc/parse"
)

type GameScore struct {
	parse.ParseObject
	CheatMode  bool    `json:"cheatMode,omitempty"`
	PlayerName string  `json:"playerName,omitempty"`
	Score      float64 `json:"score,omitempty"`
}

func mkclient(t *testing.T) *parse.Client {
	appID := os.Getenv("APPLICATION_ID")
	apiKey := os.Getenv("REST_API_KEY")
	client, err := parse.NewClient(appID, apiKey)
	if testing.Verbose() {
		client.TraceOn(log.New(os.Stderr, "[parse test] ", log.LstdFlags))
	}
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestClassesEndpoints(t *testing.T) {
	client := mkclient(t)

	// Create
	objID, err := client.Create(&GameScore{
		CheatMode:  true,
		PlayerName: "Sean Plott",
		Score:      1337,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Update
	upd := &GameScore{
		Score: 31337,
	}
	upd.ID = objID
	if _, err := client.Update(upd); err != nil {
		t.Fatal(err)
	}

	// Get
	obj := new(GameScore)
	if err := client.Get(objID, obj); err != nil {
		t.Fatal(err)
	}
	if obj.Score != 31337 {
		t.Errorf("Update failed, Score is not 31337\nObj:%+v", obj)
	}

	// Query
	objs := []GameScore{}
	if err := client.Query(nil, &objs); err != nil {
		t.Fatal(err)
	}

	// Delete
	if err := client.Delete(obj); err != nil {
		t.Fatal(err)
	}

	// Query again and ensure count is reduced
	prevLen := len(objs)
	objs = []GameScore{}
	if err := client.Query(nil, &objs); err != nil {
		t.Fatal(err)
	}
	if len(objs) >= prevLen {
		t.Fatal("Number of objects didn't decrease")
	}

}

func TestFileOperations(t *testing.T) {
	client := mkclient(t)
	client = client.WithMasterKey(os.Getenv("MASTER_KEY"))

	f, err := client.UploadFile("answer.txt", strings.NewReader("42"), "text/plain")
	if err != nil {
		t.Fatal(err)
	}
	err = client.DeleteFile(f.Name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserOperations(t *testing.T) {
	client := mkclient(t)

	type MyParseUser struct {
		parse.ParseUser
	}
	user := MyParseUser{}
	user.Username = "joe" + fmt.Sprint(time.Now().Unix())
	user.Password = "kinginyell0"

	// Create user
	loadedUser, err := client.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	// Login user
	loggedInUser := MyParseUser{}
	err = client.LoginUser(user.Username, user.Password, &loggedInUser)
	if err != nil {
		t.Fatal(err)
	}

	// Delete user
	client = client.WithSessionToken(loadedUser.SessionToken)
	err = client.DeleteUser(loggedInUser)
	if err != nil {
		t.Fatal(err)
	}

}
