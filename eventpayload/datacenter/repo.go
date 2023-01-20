package datacenter

type RepoRefsChanged struct {
	EventKey string `json:"eventKey"`
	Date     string `json:"date"`
	Actor    struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		Active       bool   `json:"active"`
		DisplayName  string `json:"displayName"`
		ID           int    `json:"id"`
		Slug         string `json:"slug"`
		Type         string `json:"type"`
		Links        struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"actor"`
	Repository struct {
		Slug          string `json:"slug"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		HierarchyID   string `json:"hierarchyId"`
		ScmID         string `json:"scmId"`
		State         string `json:"state"`
		StatusMessage string `json:"statusMessage"`
		Forkable      bool   `json:"forkable"`
		Project       struct {
			Key    string `json:"key"`
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Public bool   `json:"public"`
			Type   string `json:"type"`
			Links  struct {
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"project"`
		Public   bool `json:"public"`
		Archived bool `json:"archived"`
		Links    struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"repository"`
	Changes []struct {
		Ref struct {
			ID        string `json:"id"`
			DisplayID string `json:"displayId"`
			Type      string `json:"type"`
		} `json:"ref"`
		RefID    string `json:"refId"`
		FromHash string `json:"fromHash"`
		ToHash   string `json:"toHash"`
		Type     string `json:"type"`
	} `json:"changes"`
}

type RepoModified struct {
	EventKey string `json:"eventKey"`
	Date     string `json:"date"`
	Actor    struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		Active       bool   `json:"active"`
		DisplayName  string `json:"displayName"`
		ID           int    `json:"id"`
		Slug         string `json:"slug"`
		Type         string `json:"type"`
		Links        struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"actor"`
	Old struct {
		Slug          string `json:"slug"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		HierarchyID   string `json:"hierarchyId"`
		ScmID         string `json:"scmId"`
		State         string `json:"state"`
		StatusMessage string `json:"statusMessage"`
		Forkable      bool   `json:"forkable"`
		Project       struct {
			Key    string `json:"key"`
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Public bool   `json:"public"`
			Type   string `json:"type"`
			Links  struct {
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"project"`
		Public   bool `json:"public"`
		Archived bool `json:"archived"`
		Links    struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"old"`
	New struct {
		Slug          string `json:"slug"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		HierarchyID   string `json:"hierarchyId"`
		ScmID         string `json:"scmId"`
		State         string `json:"state"`
		StatusMessage string `json:"statusMessage"`
		Forkable      bool   `json:"forkable"`
		Project       struct {
			Key    string `json:"key"`
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Public bool   `json:"public"`
			Type   string `json:"type"`
			Links  struct {
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"project"`
		Public   bool `json:"public"`
		Archived bool `json:"archived"`
		Links    struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"new"`
}

type RepoForked struct {
	EventKey string `json:"eventKey"`
	Date     string `json:"date"`
	Actor    struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		Active       bool   `json:"active"`
		DisplayName  string `json:"displayName"`
		ID           int    `json:"id"`
		Slug         string `json:"slug"`
		Type         string `json:"type"`
		Links        struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"actor"`
	Repository struct {
		Slug          string `json:"slug"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		HierarchyID   string `json:"hierarchyId"`
		ScmID         string `json:"scmId"`
		State         string `json:"state"`
		StatusMessage string `json:"statusMessage"`
		Forkable      bool   `json:"forkable"`
		Origin        struct {
			Slug          string `json:"slug"`
			ID            int    `json:"id"`
			Name          string `json:"name"`
			HierarchyID   string `json:"hierarchyId"`
			ScmID         string `json:"scmId"`
			State         string `json:"state"`
			StatusMessage string `json:"statusMessage"`
			Forkable      bool   `json:"forkable"`
			Project       struct {
				Key    string `json:"key"`
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Public bool   `json:"public"`
				Type   string `json:"type"`
				Links  struct {
					Self []struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"links"`
			} `json:"project"`
			Public   bool `json:"public"`
			Archived bool `json:"archived"`
			Links    struct {
				Clone []struct {
					Href string `json:"href"`
					Name string `json:"name"`
				} `json:"clone"`
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"origin"`
		Project struct {
			Key   string `json:"key"`
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Type  string `json:"type"`
			Owner struct {
				Name         string `json:"name"`
				EmailAddress string `json:"emailAddress"`
				Active       bool   `json:"active"`
				DisplayName  string `json:"displayName"`
				ID           int    `json:"id"`
				Slug         string `json:"slug"`
				Type         string `json:"type"`
				Links        struct {
					Self []struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"links"`
			} `json:"owner"`
			Links struct {
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"project"`
		Public   bool `json:"public"`
		Archived bool `json:"archived"`
		Links    struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"repository"`
}

type RepoCommentAdded struct {
	EventKey string `json:"eventKey"`
	Date     string `json:"date"`
	Actor    struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		Active       bool   `json:"active"`
		DisplayName  string `json:"displayName"`
		ID           int    `json:"id"`
		Slug         string `json:"slug"`
		Type         string `json:"type"`
		Links        struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"actor"`
	Comment struct {
		Properties struct {
			RepositoryID int `json:"repositoryId"`
		} `json:"properties"`
		ID      int    `json:"id"`
		Version int    `json:"version"`
		Text    string `json:"text"`
		Author  struct {
			Name         string `json:"name"`
			EmailAddress string `json:"emailAddress"`
			Active       bool   `json:"active"`
			DisplayName  string `json:"displayName"`
			ID           int    `json:"id"`
			Slug         string `json:"slug"`
			Type         string `json:"type"`
			Links        struct {
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"author"`
		CreatedDate int64         `json:"createdDate"`
		UpdatedDate int64         `json:"updatedDate"`
		Comments    []interface{} `json:"comments"`
		Severity    string        `json:"severity"`
		State       string        `json:"state"`
	} `json:"comment"`
	Repository struct {
		Slug          string `json:"slug"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		HierarchyID   string `json:"hierarchyId"`
		ScmID         string `json:"scmId"`
		State         string `json:"state"`
		StatusMessage string `json:"statusMessage"`
		Forkable      bool   `json:"forkable"`
		Project       struct {
			Key    string `json:"key"`
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Public bool   `json:"public"`
			Type   string `json:"type"`
			Links  struct {
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"project"`
		Public   bool `json:"public"`
		Archived bool `json:"archived"`
		Links    struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"repository"`
	Commit string `json:"commit"`
}

type RepoCommentEdited struct {
	EventKey string `json:"eventKey"`
	Date     string `json:"date"`
	Actor    struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		Active       bool   `json:"active"`
		DisplayName  string `json:"displayName"`
		ID           int    `json:"id"`
		Slug         string `json:"slug"`
		Type         string `json:"type"`
		Links        struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"actor"`
	Comment struct {
		Properties struct {
			RepositoryID int `json:"repositoryId"`
		} `json:"properties"`
		ID      int    `json:"id"`
		Version int    `json:"version"`
		Text    string `json:"text"`
		Author  struct {
			Name         string `json:"name"`
			EmailAddress string `json:"emailAddress"`
			Active       bool   `json:"active"`
			DisplayName  string `json:"displayName"`
			ID           int    `json:"id"`
			Slug         string `json:"slug"`
			Type         string `json:"type"`
			Links        struct {
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"author"`
		CreatedDate int64         `json:"createdDate"`
		UpdatedDate int64         `json:"updatedDate"`
		Comments    []interface{} `json:"comments"`
		Severity    string        `json:"severity"`
		State       string        `json:"state"`
	} `json:"comment"`
	Repository struct {
		Slug          string `json:"slug"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		HierarchyID   string `json:"hierarchyId"`
		ScmID         string `json:"scmId"`
		State         string `json:"state"`
		StatusMessage string `json:"statusMessage"`
		Forkable      bool   `json:"forkable"`
		Project       struct {
			Key    string `json:"key"`
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Public bool   `json:"public"`
			Type   string `json:"type"`
			Links  struct {
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"project"`
		Public   bool `json:"public"`
		Archived bool `json:"archived"`
		Links    struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"repository"`
	Commit          string `json:"commit"`
	PreviousComment string `json:"previousComment"`
}

type RepoCommentDeleted struct {
	EventKey string `json:"eventKey"`
	Date     string `json:"date"`
	Actor    struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		Active       bool   `json:"active"`
		DisplayName  string `json:"displayName"`
		ID           int    `json:"id"`
		Slug         string `json:"slug"`
		Type         string `json:"type"`
		Links        struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"actor"`
	Comment struct {
		ID      int    `json:"id"`
		Version int    `json:"version"`
		Text    string `json:"text"`
		Author  struct {
			Name         string `json:"name"`
			EmailAddress string `json:"emailAddress"`
			Active       bool   `json:"active"`
			DisplayName  string `json:"displayName"`
			ID           int    `json:"id"`
			Slug         string `json:"slug"`
			Type         string `json:"type"`
			Links        struct {
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"author"`
		CreatedDate int64         `json:"createdDate"`
		UpdatedDate int64         `json:"updatedDate"`
		Comments    []interface{} `json:"comments"`
		Severity    string        `json:"severity"`
		State       string        `json:"state"`
	} `json:"comment"`
	Repository struct {
		Slug          string `json:"slug"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		HierarchyID   string `json:"hierarchyId"`
		ScmID         string `json:"scmId"`
		State         string `json:"state"`
		StatusMessage string `json:"statusMessage"`
		Forkable      bool   `json:"forkable"`
		Project       struct {
			Key    string `json:"key"`
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Public bool   `json:"public"`
			Type   string `json:"type"`
			Links  struct {
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"project"`
		Public   bool `json:"public"`
		Archived bool `json:"archived"`
		Links    struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"repository"`
	Commit string `json:"commit"`
}

type RepoMirrorSynchronized struct {
	EventKey     string `json:"eventKey"`
	Date         string `json:"date"`
	MirrorServer struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"mirrorServer"`
	SyncType         string `json:"syncType"`
	RefLimitExceeded bool   `json:"refLimitExceeded"`
	Repository       struct {
		Slug          string `json:"slug"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		ScmID         string `json:"scmId"`
		State         string `json:"state"`
		StatusMessage string `json:"statusMessage"`
		Forkable      bool   `json:"forkable"`
		Project       struct {
			Key    string `json:"key"`
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Public bool   `json:"public"`
			Type   string `json:"type"`
		} `json:"project"`
		Public bool `json:"public"`
		Links  struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"repository"`
	Changes []struct {
		Ref struct {
			ID        string `json:"id"`
			DisplayID string `json:"displayId"`
			Type      string `json:"type"`
		} `json:"ref"`
		RefID    string `json:"refId"`
		FromHash string `json:"fromHash"`
		ToHash   string `json:"toHash"`
		Type     string `json:"type"`
	} `json:"changes"`
}
