package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/logger"
)

func isListmonkConfigured() bool {
	return os.Getenv("LISTMONK_ENABLED") != "false" && os.Getenv("LISTMONK_URL") != ""
}

type listmonkProvider struct{}

func (p *listmonkProvider) sendTemplate(to string, tmpl Template, data map[string]any) error {
	templateID := tmpl.listmonkTemplateID()
	if templateID == 0 {
		return fmt.Errorf("listmonk: no template ID configured for this email (set %s)", tmpl.listmonkIDEnv)
	}

	url := os.Getenv("LISTMONK_URL")
	username := os.Getenv("LISTMONK_USERNAME")
	password := os.Getenv("LISTMONK_PASSWORD")

	payload := bson.M{
		"subscriber_email": to,
		"template_id":      templateID,
		"data":             data,
		"content_type":     "html",
	}
	if tmpl.FromEmail != "" {
		payload["from_email"] = tmpl.FromEmail
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/tx", url), bytes.NewBuffer(body))
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (p *listmonkProvider) ensureSubscriber(email string, sendMarketing bool) {
	if exists, _ := p.subscriberExists(email); !exists {
		p.addSubscriber(email, "", "", "", nil, sendMarketing)
	}
}

func (p *listmonkProvider) subscriberExists(email string) (bool, *int) {
	url := os.Getenv("LISTMONK_URL")
	username := os.Getenv("LISTMONK_USERNAME")
	password := os.Getenv("LISTMONK_PASSWORD")

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/subscribers?query=subscribers.email='%s'", url, email), nil)
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
		return false, nil
	}
	defer resp.Body.Close()

	var response struct {
		Data struct {
			Results []struct {
				Id int `json:"id"`
			} `json:"results"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logger.StdErr.Println(err)
		return false, nil
	}

	if len(response.Data.Results) > 0 {
		return true, &response.Data.Results[0].Id
	}
	return false, nil
}

func (p *listmonkProvider) addSubscriber(email, firstName, lastName, picture string, subscriberID *int, sendMarketing bool) {
	url := os.Getenv("LISTMONK_URL")
	username := os.Getenv("LISTMONK_USERNAME")
	password := os.Getenv("LISTMONK_PASSWORD")
	listIDStr := os.Getenv("LISTMONK_LIST_ID")

	args := bson.M{
		"email":  email,
		"name":   firstName + " " + lastName,
		"status": "enabled",
		"attribs": bson.M{
			"firstName": firstName,
			"lastName":  lastName,
			"picture":   picture,
		},
		"preconfirm_subscriptions": true,
	}
	if sendMarketing {
		if listID, err := strconv.Atoi(listIDStr); err == nil {
			args["lists"] = bson.A{listID}
		}
	}

	body, _ := json.Marshal(args)

	var req *http.Request
	if subscriberID != nil {
		req, _ = http.NewRequest("PUT", fmt.Sprintf("%s/api/subscribers/%d", url, *subscriberID), bytes.NewBuffer(body))
	} else {
		req, _ = http.NewRequest("POST", fmt.Sprintf("%s/api/subscribers", url), bytes.NewBuffer(body))
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
		return
	}
	defer resp.Body.Close()
}
