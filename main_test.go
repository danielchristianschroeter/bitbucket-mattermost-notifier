package main

import (
	"bitbucket-mattermost-notifier/eventpayload/datacenter"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestJsonUnmarshal(t *testing.T) {
	//Example json payload content from file
	f, err := os.Open("./eventpayload/datacenter/example_payload/repo_refs_changed.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	body := make([]byte, 500)
	_, err = f.Read(body)
	if err != nil {
		t.Fatal(err)
	}
	// initialize payload struct
	var payload datacenter.RepoPush
	// create a test http response writer
	w := &responseWriter{}
	// call JsonUnmarshal
	err = JsonUnmarshal(body, &payload, w, "repo:refs_changed")
	if err != nil {
		t.Errorf("JsonUnmarshal returned an error: %v", err)
	}
	// check if payload has been correctly unmarshalled
	if reflect.TypeOf(payload) != reflect.TypeOf(datacenter.RepoPush{}) {
		t.Errorf("payload type is %v, expected %v", reflect.TypeOf(payload), reflect.TypeOf(datacenter.RepoPush{}))
	}
	// check if the EventKey in the payload is correct
	val := reflect.ValueOf(payload).Elem()
	if "repo:refs_changed" != val.FieldByName("EventKey").String() {
		t.Errorf("EventKey in payload is %s, expected %s", val.FieldByName("EventKey").String(), "repo:refs_changed")
	}
	// check if the http response has been correctly written
	if w.status != http.StatusInternalServerError {
		t.Errorf("http status is %v, expected %v", w.status, http.StatusInternalServerError)
	}
}

type responseWriter struct {
	status int
}

func (w *responseWriter) Header() http.Header {
	return nil
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return 0, nil
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

func TestMainRefsChanged(t *testing.T) {
	var payload datacenter.RepoPush
	//Example json payload content from file
	f, err := os.Open("./eventpayload/datacenter/example_payload/repo_refs_changed.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	body := make([]byte, 500)
	_, err = f.Read(body)
	if err != nil {
		t.Fatal(err)
	}
	// call main function
	main()
	// check if the attachments are created correctly
	if len(attachments) != 1 {
		t.Errorf("attachments length is %d, expected %d", len(attachments), 1)
	}
	// check if the attachments have the correct Fallback
	if attachments[0].Fallback != "One or more commits from "+payload.Actor.DisplayName+" to repository "+payload.Repository.Name+" ("+payload.Repository.Project.Name+")" {
		t.Errorf("attachments Fallback is %s, expected %s", attachments[0].Fallback, "One or more commits from "+payload.Actor.DisplayName+" to repository "+payload.Repository.Name+" ("+payload.Repository.Project.Name+")")
	}
	// check if the attachments have the correct Title
	if attachments[0].Title != "One or more commits from "+payload.Actor.DisplayName {
		t.Errorf("attachments Title is %s, expected %s", attachments[0].Title, "One or more commits from "+payload.Actor.DisplayName)
	}
	// check if the attachments have the correct Color
	if attachments[0].Color != "#"+GetConfigValue("INFO_COLOR") {
		t.Errorf("attachments Color is %s, expected %s", attachments[0].Color, "#"+GetConfigValue("INFO_COLOR"))
	}
	// check if the attachments Fields are correct
	//...
}
