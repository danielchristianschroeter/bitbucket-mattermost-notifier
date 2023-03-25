package main

import (
	"bitbucket-mattermost-notifier/eventpayload/datacenter"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
)

var version = "DEV"
var srv http.Server

func JsonUnmarshal(body []byte, payload interface{}, w http.ResponseWriter, eventKey string) error {
	var err error
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error in %v: %v", eventKey, err)
		return err
	}
	val := reflect.ValueOf(payload).Elem()
	if eventKey != val.FieldByName("EventKey").String() {
		http.Error(w, "Header value of key X-Event-Key ("+eventKey+") does not match with eventKey received in json payload ("+val.FieldByName("EventKey").String()+").", http.StatusInternalServerError)
		return fmt.Errorf("header value of key X-Event-Key (%s) does not match with eventKey received in json payload (%s)", eventKey, val.FieldByName("EventKey").String())
	}
	return err
}

func HandleEvent(w http.ResponseWriter, r *http.Request, eventKey string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var attachments []Attachment
	switch eventKey {
	case "repo:refs_changed":
		var payload datacenter.RepoRefsChanged
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		// Create a list of all changes
		var changesList string
		for _, change := range payload.Changes {
			changesList += "Branch: " + change.Ref.DisplayID + " - FromHash: " + change.FromHash + " - ToHash: " + change.ToHash + " - Type: " + change.Type + "\n"
		}
		attachments = []Attachment{
			{
				Fallback:  "One or more commits pushed by " + payload.Actor.DisplayName + " to repository " + payload.Repository.Name + " (" + payload.Repository.Project.Name + ")",
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Title:     "One or more commits pushed by " + payload.Actor.DisplayName,
				TitleLink: payload.Repository.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Repository",
						Value: payload.Repository.Name + " (" + payload.Repository.Project.Name + ")",
					},
					{
						Short: false,
						Title: "Changes",
						Value: changesList,
					},
				},
			},
		}
	case "repo:modified":
		var payload datacenter.RepoModified
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback:  "Name of repository updated by " + payload.Actor.DisplayName + " Old value: " + payload.Old.Name + " " + "(" + payload.Old.Project.Name + ")" + "=> New value: " + payload.New.Name + " " + "(" + payload.New.Project.Name + ")",
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Title:     "Name of repository updated by " + payload.Actor.DisplayName,
				TitleLink: payload.New.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Old repository name",
						Value: payload.Old.Name + " " + "(" + payload.Old.Project.Name + ")",
					},
					{
						Short: false,
						Title: "New repository name",
						Value: payload.New.Name + " " + "(" + payload.New.Project.Name + ")",
					},
				},
			},
		}
	case "repo:forked":
		var payload datacenter.RepoForked
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback:  "Repository " + payload.Repository.Origin.Name + " (" + payload.Repository.Origin.Project.Name + ") forked by " + payload.Actor.DisplayName + " to " + payload.Repository.Name,
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Title:     "Repository forked by " + payload.Actor.DisplayName,
				TitleLink: payload.Repository.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Origin repository",
						Value: payload.Repository.Origin.Name + " (" + payload.Repository.Origin.Project.Name + ")",
					},
					{
						Short: false,
						Title: "New forked repository",
						Value: payload.Repository.Name + " " + "(" + payload.Repository.Project.Name + ")",
					},
				},
			},
		}
	case "repo:comment:added":
		var payload datacenter.RepoCommentAdded
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback:  "Comment for commit created by " + payload.Actor.DisplayName + " for commit with hash " + payload.Commit + " in repository " + payload.Repository.Name + " (" + payload.Repository.Project.Name + ")",
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Title:     "Comment for commit created by " + payload.Actor.DisplayName,
				TitleLink: payload.Repository.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Repository",
						Value: payload.Repository.Name + " (" + payload.Repository.Project.Name + ")",
					},
					{
						Short: false,
						Title: "Commit",
						Value: payload.Commit,
					},
					{
						Short: false,
						Title: "Comment",
						Value: payload.Comment.Text,
					},
				},
			},
		}
	case "repo:comment:edited":
		var payload datacenter.RepoCommentEdited
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback: "Comment for commit edited by " + payload.Actor.DisplayName + " for commit with hash " + payload.Commit + " in repository " + payload.Repository.Name + " (" + payload.Repository.Project.Name + ")",
				Color:    "#" + GetConfigValue("INFO_COLOR"),
				Title:    "Comment for commit edited by " + payload.Actor.DisplayName,
				Fields: []Field{
					{
						Short: false,
						Title: "Repository",
						Value: payload.Repository.Name + " (" + payload.Repository.Project.Name + ")",
					},
					{
						Short: false,
						Title: "Commit",
						Value: payload.Commit,
					},
					{
						Short: false,
						Title: "Previous Comment",
						Value: payload.PreviousComment,
					},
					{
						Short: false,
						Title: "New Comment",
						Value: payload.Comment.Text,
					},
				},
			},
		}
	case "repo:comment:deleted":
		var payload datacenter.RepoCommentDeleted
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback:  "Comment for commit deleted by " + payload.Actor.DisplayName + " for commit with hash " + payload.Commit + " in repository " + payload.Repository.Name + " (" + payload.Repository.Project.Name + ")",
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Text:      "Comment for commit deleted by " + payload.Actor.DisplayName,
				TitleLink: payload.Repository.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Repository",
						Value: payload.Repository.Name + " (" + payload.Repository.Project.Name + ")",
					},
					{
						Short: false,
						Title: "Commit",
						Value: payload.Commit,
					},
					{
						Short: false,
						Title: "Deleted Comment",
						Value: payload.Comment.Text,
					},
				},
			},
		}
	case "mirror:repo_synchronized":
		var payload datacenter.RepoMirrorSynchronized
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		// Create a list of all changes
		var changesList string
		for _, change := range payload.Changes {
			changesList += "Branch: " + change.Ref.DisplayID + " - FromHash: " + change.FromHash + " - ToHash: " + change.ToHash + " - Type: " + change.Type + "\n"
		}
		attachments = []Attachment{
			{
				Fallback:  "Mirror server " + payload.MirrorServer.Name + " successfully synchronized for " + payload.Repository.Name + " (" + payload.Repository.Project.Name + ")",
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Title:     "Mirror server " + payload.MirrorServer.Name + " successfully synchronized",
				TitleLink: payload.Repository.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Sync type",
						Value: payload.SyncType,
					},
					{
						Short: false,
						Title: "Repository",
						Value: payload.Repository.Name + " (" + payload.Repository.Project.Name + ")",
					},
					{
						Short: false,
						Title: "Changes",
						Value: changesList,
					},
				},
			},
		}
	case "pr:opened":
		var payload datacenter.PROpened
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		// Create a comma seperated list for all reviewers
		var reviewersList string
		if len(payload.PullRequest.Reviewers) == 0 {
			reviewersList = ""
		} else {
			for _, reviewer := range payload.PullRequest.Reviewers {
				reviewersList += "@" + reviewer.User.Name + ", "
			}
			// Remove the last comma and space from the reviewersList
			reviewersList = reviewersList[:len(reviewersList)-2]
		}

		attachments = []Attachment{
			{
				Fallback:  "Pull request opened by " + payload.Actor.DisplayName,
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Title:     "Pull request opened by " + payload.Actor.DisplayName,
				TitleLink: payload.PullRequest.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "Description",
						Value: payload.PullRequest.Description,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
					{
						Short: false,
						Title: "Reviewers",
						Value: reviewersList,
					},
				},
			},
		}
	case "pr:from_ref_updated":
		var payload datacenter.PRFromRefUpdated
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback:  "Source branch of pull request updated by " + payload.Actor.DisplayName,
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Title:     "Source branch of pull request updated by " + payload.Actor.DisplayName,
				TitleLink: payload.PullRequest.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "Description",
						Value: payload.PullRequest.Description,
					},
					{
						Short: false,
						Title: "State",
						Value: payload.PullRequest.State,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
					{
						Short: false,
						Title: "Previous Commit",
						Value: payload.PreviousFromHash,
					},
				},
			},
		}
	case "pr:modified":
		var payload datacenter.PRModified
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback:  "Pull request modified by " + payload.Actor.DisplayName,
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Title:     "Pull request modified by " + payload.Actor.DisplayName,
				TitleLink: payload.PullRequest.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "Description",
						Value: payload.PullRequest.Description,
					},
					{
						Short: false,
						Title: "State",
						Value: payload.PullRequest.State,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
					{
						Short: false,
						Title: "Previous Title",
						Value: payload.PreviousTitle,
					},
					{
						Short: false,
						Title: "Previous Description",
						Value: payload.PreviousDescription,
					},
					{
						Short: false,
						Title: "Previous Target",
						Value: payload.PreviousTarget.DisplayID + " - " + payload.PreviousTarget.LatestCommit,
					},
				},
			},
		}
	case "pr:reviewer:updated":
		var payload datacenter.PRReviewersUpdated
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}

		// Create a comma seperated list for all reviewers
		var reviewersList string
		if len(payload.PullRequest.Reviewers) == 0 {
			reviewersList = ""
		} else {
			for _, reviewer := range payload.PullRequest.Reviewers {
				reviewersList += reviewer.User.DisplayName + ", "
			}
			// Remove the last comma and space from the reviewersList
			reviewersList = reviewersList[:len(reviewersList)-2]
		}

		// Create a comma seperated list for all added reviewers
		var reviewersAdded string
		if len(payload.AddedReviewers) == 0 {
			reviewersAdded = ""
		} else {
			for _, revieweradded := range payload.AddedReviewers {
				reviewersAdded += revieweradded.DisplayName + ", "
			}
			// Remove the last comma and space from the reviewersAdded
			reviewersAdded = reviewersAdded[:len(reviewersAdded)-2]
		}

		// Create a comma seperated list for all removed reviewers
		var reviewersRemoved string
		if len(payload.RemovedReviewers) == 0 {
			reviewersRemoved = ""
		} else {
			for _, reviewerremoved := range payload.RemovedReviewers {
				reviewersRemoved += reviewerremoved.DisplayName + ", "
			}
			// Remove the last comma and space from the reviewersAdded
			reviewersRemoved = reviewersRemoved[:len(reviewersRemoved)-2]
		}

		attachments = []Attachment{
			{
				Fallback: "Reviewers on pull request updated by " + payload.Actor.DisplayName,
				Color:    "#" + GetConfigValue("INFO_COLOR"),
				Title:    "Reviewers on pull request updated by " + payload.Actor.DisplayName,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "Description",
						Value: payload.PullRequest.Description,
					},
					{
						Short: false,
						Title: "State",
						Value: payload.PullRequest.State,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
					{
						Short: false,
						Title: "Reviewers",
						Value: reviewersList,
					},
					{
						Short: false,
						Title: "Added Reviewers",
						Value: reviewersAdded,
					},
					{
						Short: false,
						Title: "Removed Reviewers",
						Value: reviewersRemoved,
					},
				},
			},
		}
	case "pr:reviewer:approved":
		var payload datacenter.PRReviewerApproved
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		// Fetch all approved reviewers
		var approvedReviewersList string
		for _, reviewer := range payload.PullRequest.Reviewers {
			if reviewer.Status == "APPROVED" {
				approvedReviewersList += reviewer.User.DisplayName + ", "
			}
		}
		// Remove the last comma and space from the approvedReviewersList
		approvedReviewersList = approvedReviewersList[:len(approvedReviewersList)-2]
		// Count names in approvedReviewersList
		approvedReviewers := strings.Split(approvedReviewersList, ",")
		approvedReviewersCount := len(approvedReviewers)
		approvedReviewersCountString := strconv.Itoa(approvedReviewersCount)

		attachments = []Attachment{
			{
				Fallback:  "Pull request approved by " + payload.Actor.DisplayName,
				Color:     "#" + GetConfigValue("SUCCESS_COLOR"),
				Title:     "Pull request approved by " + payload.Actor.DisplayName,
				TitleLink: payload.PullRequest.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "Description",
						Value: payload.PullRequest.Description,
					},
					{
						Short: false,
						Title: "State",
						Value: payload.PullRequest.State + " / " + payload.PreviousStatus,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
					{
						Short: false,
						Title: "Approval Count",
						Value: approvedReviewersCountString + " - " + approvedReviewersList,
					},
				},
			},
		}
	case "pr:reviewer:unapproved":
		var payload datacenter.PRReviewerUnapproved
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback: "Pull request unapproved by " + payload.Actor.DisplayName,
				Color:    "#" + GetConfigValue("INFO_COLOR"),
				Title:    "Pull request unapproved by " + payload.Actor.DisplayName,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "Description",
						Value: payload.PullRequest.Description,
					},
					{
						Short: false,
						Title: "State",
						Value: payload.PullRequest.State + " / " + payload.PreviousStatus,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
				},
			},
		}
	case "pr:reviewer:needs_work":
		var payload datacenter.PRReviewerNeedsWork
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback: "Pull request marked as needs work by " + payload.Actor.DisplayName,
				Color:    "#" + GetConfigValue("WARNING_COLOR"),
				Title:    "Pull request marked as needs work by " + payload.Actor.DisplayName,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "Description",
						Value: payload.PullRequest.Description,
					},
					{
						Short: false,
						Title: "State",
						Value: payload.PullRequest.State + " / " + payload.PreviousStatus,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
				},
			},
		}
	case "pr:merged":
		var payload datacenter.PRMerged
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback:  "Pull request merged by " + payload.Actor.DisplayName,
				Color:     "#" + GetConfigValue("SUCCESS_COLOR"),
				Title:     "Pull request merged by " + payload.Actor.DisplayName,
				TitleLink: payload.PullRequest.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "Description",
						Value: payload.PullRequest.Description,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
					{
						Short: false,
						Title: "Merge Commit",
						Value: payload.PullRequest.Properties.MergeCommit.DisplayID + " - " + payload.PullRequest.Properties.MergeCommit.ID,
					},
				},
			},
		}
	case "pr:declined":
		var payload datacenter.PRDeclined
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback:  "Pull request declined by " + payload.Actor.DisplayName,
				Color:     "#" + GetConfigValue("DECLINED_COLOR"),
				Title:     "Pull request declined by " + payload.Actor.DisplayName,
				TitleLink: payload.PullRequest.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "Description",
						Value: payload.PullRequest.Description,
					},
					{
						Short: false,
						Title: "State",
						Value: payload.PullRequest.State,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
				},
			},
		}
	case "pr:deleted":
		var payload datacenter.PRDeleted
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback:  "Pull request deleted by " + payload.Actor.DisplayName,
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Title:     "Pull request deleted by " + payload.Actor.DisplayName,
				TitleLink: payload.PullRequest.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "Description",
						Value: payload.PullRequest.Description,
					},
					{
						Short: false,
						Title: "State",
						Value: payload.PullRequest.State,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
				},
			},
		}
	case "pr:comment:added":
		var payload datacenter.PRCommentAdded
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback:  "Comment for pull request added by " + payload.Actor.DisplayName,
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Title:     "Comment for pull request added by " + payload.Actor.DisplayName,
				TitleLink: payload.PullRequest.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "State",
						Value: payload.PullRequest.State,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
					{
						Short: false,
						Title: "Comment",
						Value: payload.Comment.Text,
					},
				},
			},
		}
	case "pr:comment:edited":
		var payload datacenter.PRCommentEdited
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback:  "Comment for pull request edited by " + payload.Actor.DisplayName,
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Title:     "Comment for pull request edited by " + payload.Actor.DisplayName,
				TitleLink: payload.PullRequest.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "State",
						Value: payload.PullRequest.State,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
					{
						Short: false,
						Title: "Previous Comment",
						Value: payload.PreviousComment,
					},
					{
						Short: false,
						Title: "New Comment",
						Value: payload.Comment.Text,
					},
				},
			},
		}
	case "pr:comment:deleted":
		var payload datacenter.PRCommentDeleted
		if err := JsonUnmarshal(body, &payload, w, eventKey); err != nil {
			log.Printf("Error in JsonUnmarshal for eventKey %v: %v", eventKey, err)
			return
		}
		attachments = []Attachment{
			{
				Fallback:  "Comment for pull request deleted by " + payload.Actor.DisplayName,
				Color:     "#" + GetConfigValue("INFO_COLOR"),
				Title:     "Comment for pull request deleted by " + payload.Actor.DisplayName,
				TitleLink: payload.PullRequest.Links.Self[0].Href,
				Fields: []Field{
					{
						Short: false,
						Title: "Title",
						Value: payload.PullRequest.Title,
					},
					{
						Short: false,
						Title: "State",
						Value: payload.PullRequest.State,
					},
					{
						Short: true,
						Title: "From",
						Value: payload.PullRequest.FromRef.Repository.Name + " (" + payload.PullRequest.FromRef.Repository.Project.Name + ") - " + payload.PullRequest.FromRef.DisplayID + " - " + payload.PullRequest.FromRef.LatestCommit,
					},
					{
						Short: true,
						Title: "To",
						Value: payload.PullRequest.ToRef.Repository.Name + " (" + payload.PullRequest.ToRef.Repository.Project.Name + ") - " + payload.PullRequest.ToRef.DisplayID + " - " + payload.PullRequest.ToRef.LatestCommit,
					},
					{
						Short: false,
						Title: "Deleted Comment",
						Value: payload.Comment.Text,
					},
				},
			},
		}
	case "diagnostics:ping":
		attachments = []Attachment{
			{
				Fallback: "Test connection successfully executed from Bitbucket",
				Color:    "#" + GetConfigValue("INFO_COLOR"),
				Title:    "Test connection successfully executed from Bitbucket",
			},
		}
	default:
		log.Printf("Unknown X-Event-Key received: %v", eventKey)
	}
	// If attachments struct is not empty, send message to mattermost webhook url
	if len(attachments) > 0 {
		if err := SendMessage(GetConfigValue("MATTERMOST_WEBHOOKURL"), GetConfigValue("MATTERMOST_CHANNEL"), GetConfigValue("MATTERMOST_USERNAME"), GetConfigValue("MESSAGE_ICON_EMOJI"), "", attachments); err != nil {
			log.Printf("Error sending message to %v: %v", GetConfigValue("MATTERMOST_WEBHOOKURL"), err)
		} else {
			log.Printf("Message for event %v sent successfully to %v", eventKey, GetConfigValue("MATTERMOST_WEBHOOKURL"))
		}
	}
}

func main() {
	log.Printf("Bitbucket-Mattermost-Notifier (v%s) started.\n", version)
	log.Printf("Loaded configuration:\n")
	log.Printf("LISTEN_PORT: %s\n", GetConfigValue("LISTEN_PORT"))
	log.Printf("MATTERMOST_WEBHOOKURL: %s\n", GetConfigValue("MATTERMOST_WEBHOOKURL"))
	log.Printf("MATTERMOST_USERNAME: %s\n", GetConfigValue("MATTERMOST_USERNAME"))
	log.Printf("MATTERMOST_CHANNEL: %s\n", GetConfigValue("MATTERMOST_CHANNEL"))
	log.Printf("MESSAGE_ICON_EMOJI: %s\n", GetConfigValue("MESSAGE_ICON_EMOJI"))
	log.Printf("INFO_COLOR: %s\n", GetConfigValue("INFO_COLOR"))
	log.Printf("WARNING_COLOR: %s\n", GetConfigValue("WARNING_COLOR"))
	log.Printf("SUCCESS_COLOR: %s\n", GetConfigValue("SUCCESS_COLOR"))
	SetupGracefulShutdown()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			eventKey := r.Header.Get("X-Event-Key")
			log.Printf("Incoming payload received for header key X-Event-Key with value %v", eventKey)
			// Handle different event keys from Bitbucket
			HandleEvent(w, r, eventKey)
		}
	})
	log.Print("Listening now on port ", GetConfigValue("LISTEN_PORT")+"...")
	srv.Addr = ":" + GetConfigValue("LISTEN_PORT")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Error received on listening on port: %v", err)
	}
}

func SetupGracefulShutdown() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		if err := srv.Close(); err != nil {
			log.Printf("Error on HTTP server Close: %v", err)
		}
		os.Exit(0)
	}()
}
