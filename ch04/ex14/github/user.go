// Copyright 2017 budougumi0617 All Rights Reserved.

package github

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

// User user information
// https://developer.github.com/v3/users/#get-all-users
// https://api.github.com/users
type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// Users manages GitHub user data.
type Users struct {
	Link  string
	Items []*User
}

// GetUsers gets issues from web.
func (u *Users) GetUsers(c *Client) error {
	req, _ := c.NewRequest(context.Background(), "GET", "https://api.github.com/users", nil)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Printf("Get users from %v\n", req)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status is not ok: %v", resp.StatusCode)
	}

	if err := DecodeBody(resp, &(u.Items)); err != nil {
		return err
	}
	return nil
}

// RenderUsers writes HTML table.
func (u *Users) RenderUsers(w io.Writer) {
	userList := template.Must(template.New("userList").Parse(`
	<h1>GitHub Users</h1>
	<table>
	  <tr><th>Name</th><th>URL</th></tr>
	  {{range .Items}}
	  <tr>
	    <td>{{.Login}}</td><td>{{.HTMLURL}}</td>
	  </tr>
	  {{end}}
	</table>
	`))

	if err := userList.Execute(w, &u); err != nil {
		fmt.Fprintf(w, "%v\n", err)
	}
}
