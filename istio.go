// Copyright 2017 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type IstioClient struct {
	httpClient          *http.Client
	configAPIService    string
	mixerAPIService     string
	serviceGraphService string
	password            string
	username            string
}

func NewIstioClient(username, password, configAPIService, mixerAPIService, serviceGraphService string) *IstioClient {
	httpClient := &http.Client{}
	return &IstioClient{
		httpClient:          httpClient,
		configAPIService:    configAPIService,
		mixerAPIService:     mixerAPIService,
		serviceGraphService: serviceGraphService,
		password:            password,
		username:            username,
	}
}

func (c *IstioClient) AllowAccess(to, from string) error {
	urlStr := fmt.Sprintf("http://%s/api/v1/scopes/global/subjects/%s.default.svc.cluster.local/rules", c.mixerAPIService, to)
	request, err := http.NewRequest("DELETE", urlStr, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		log.Printf("AllowAccess error: non-200 status code: %d", resp.StatusCode)
		return fmt.Errorf("AllowAccess error: non-200 status code: %d", resp.StatusCode)
	}

	return nil

}

func (c *IstioClient) DenyAccess(to, from string) error {
	urlStr := fmt.Sprintf("http://%s/api/v1/scopes/global/subjects/%s.default.svc.cluster.local/rules", c.mixerAPIService, to)

	mixerRule := MixerRule{
		Rules: []Rule{
			Rule{
				Selector: fmt.Sprintf("source.labels[\"app\"]==\"%s\"", from),
				Aspects: []Aspect{
					Aspect{Kind: "denials"},
				},
			},
		},
	}

	body, err := json.Marshal(&mixerRule)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", urlStr, bytes.NewReader(body))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		log.Printf("DenyAccess error: non-200 status code: %d", resp.StatusCode)
		return fmt.Errorf("DenyAccess error: non-200 status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *IstioClient) GetRouteRule(name string) (*RouteRule, error) {
	urlStr := fmt.Sprintf("http://%s/v1alpha1/config/route-rule/default/%s-default", c.configAPIService, name)

	request, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(c.username, c.password)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Printf("GetRouteRule error: non-200 status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("GetRouteRule error: non-200 status code: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var routeRule RouteRule
	err = json.Unmarshal(data, &routeRule)
	if err != nil {
		return nil, err
	}

	return &routeRule, nil
}

func (c *IstioClient) GetTopology() (*Topology, error) {
	urlStr := fmt.Sprintf("http://%s/graph", c.serviceGraphService)
	request, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Printf("GetTopology error: non-200 status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("GetTopology error: non-200 status code: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t Topology
	err = json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
