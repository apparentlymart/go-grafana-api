package gapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
)

// Annotation represents a Grafana API Annotation
type Annotation struct {
	ID          int64    `json:"id,omitempty"`
	AlertID     int64    `json:"alertId,omitempty"`
	DashboardID int64    `json:"dashboardId"`
	PanelID     int64    `json:"panelId"`
	UserID      int64    `json:"userId,omitempty"`
	UserName    string   `json:"userName,omitempty"`
	NewState    string   `json:"newState,omitempty"`
	PrevState   string   `json:"prevState,omitempty"`
	Time        int64    `json:"time"`
	TimeEnd     int64    `json:"timeEnd,omitempty"`
	Text        string   `json:"text"`
	Metric      string   `json:"metric,omitempty"`
	RegionID    int64    `json:"regionId,omitempty"`
	Type        string   `json:"type,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	IsRegion    bool     `json:"isRegion,omitempty"`
}

// GraphiteAnnotation represents a Grafana API annotation in Graphite format
type GraphiteAnnotation struct {
	What string   `json:"what"`
	When int64    `json:"when"`
	Data string   `json:"data"`
	Tags []string `json:"tags,omitempty"`
}

// Annotations fetches the annotations queried with the params it's passed
func (c *Client) Annotations(params url.Values) ([]Annotation, error) {
	req, err := c.newRequest("GET", "/api/annotation", params, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := []Annotation{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// Annotation fetches the annotation queried with the ID and params it's passed.
// It returns an error if no annotation with a matching ID is found.
func (c *Client) Annotation(id int64, params url.Values) (Annotation, error) {
	as, err := c.Annotations(params)
	if err != nil {
		return Annotation{}, err
	}

	for _, a := range as {
		if a.ID == id {
			return a, nil
		}
	}

	return Annotation{}, fmt.Errorf("annotation %v not found", id)
}

// NewAnnotation creates a new annotation with the Annotation it is passed
func (c *Client) NewAnnotation(a *Annotation) (int64, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return 0, err
	}
	req, err := c.newRequest("POST", "/api/annotations", nil, bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, errors.New(resp.Status)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	result := struct {
		ID int64 `json:"id"`
	}{}
	err = json.Unmarshal(data, &result)
	return result.ID, err
}

// NewGraphiteAnnotation creates a new annotation with the GraphiteAnnotation it is passed
func (c *Client) NewGraphiteAnnotation(gfa *GraphiteAnnotation) (int64, error) {
	data, err := json.Marshal(gfa)
	if err != nil {
		return 0, err
	}
	req, err := c.newRequest("POST", "/api/annotations/graphite", nil, bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, errors.New(resp.Status)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	result := struct {
		ID int64 `json:"id"`
	}{}
	err = json.Unmarshal(data, &result)
	return result.ID, err
}

// UpdateAnnotation updates all properties an existing annotation with the Annotation it is passed.
func (c *Client) UpdateAnnotation(id int64, a *Annotation) (string, error) {
	path := fmt.Sprintf("/api/annotations/%d", id)
	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	req, err := c.newRequest("PUT", path, nil, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	result := struct {
		Message string `json:"message"`
	}{}
	err = json.Unmarshal(data, &result)
	return result.Message, err
}

// PatchAnnotation updates one or more properties of an existing annotation that matches the specified ID.
func (c *Client) PatchAnnotation(id int64, a *Annotation) (string, error) {
	path := fmt.Sprintf("/api/annotations/%d", id)
	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	req, err := c.newRequest("PATCH", path, nil, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	result := struct {
		Message string `json:"message"`
	}{}
	err = json.Unmarshal(data, &result)
	return result.Message, err
}

// DeleteAnnotation deletes the annotation of the ID it is passed
func (c *Client) DeleteAnnotation(id int64) (string, error) {
	path := fmt.Sprintf("/api/annotations/%d", id)
	req, err := c.newRequest("DELETE", path, nil, bytes.NewBuffer(nil))
	if err != nil {
		return "", err
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	result := struct {
		Message string `json:"message"`
	}{}
	err = json.Unmarshal(data, &result)
	return result.Message, err
}

// DeleteAnnotationByRegionID deletes the annotation corresponding to the region ID it is passed
func (c *Client) DeleteAnnotationByRegionID(id int64) (string, error) {
	path := fmt.Sprintf("/api/annotations/region/%d", id)
	req, err := c.newRequest("DELETE", path, nil, bytes.NewBuffer(nil))
	if err != nil {
		return "", err
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	result := struct {
		Message string `json:"message"`
	}{}
	err = json.Unmarshal(data, &result)
	return result.Message, err
}
