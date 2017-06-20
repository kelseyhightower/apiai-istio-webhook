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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kelseyhightower/api/ai"
)

var (
	configAPIService string
	mixerAPIService  string
	password         string
	username         string
)

func main() {
	flag.StringVar(&configAPIService, "config-api-service", "istio-pilot:8081", "The Istio config API service.")
	flag.StringVar(&mixerAPIService, "mixer-api-service", "istio-mixer:9094", "The mixer API service.")
	flag.StringVar(&password, "password", "", "The Istio config service password")
	flag.StringVar(&username, "username", "", "The Istio config service username")
	flag.Parse()

	istioClient := NewIstioClient(username, password, configAPIService, mixerAPIService)

	http.Handle("/", webhookServer(istioClient))
	log.Printf("Starting the Istio Google Action service...")
	err := http.ListenAndServeTLS(":443", "/etc/istio-webhook/tls.crt", "/etc/istio-webhook/tls.key", nil)
	log.Fatal(err)
}

type webhookHandler struct {
	istioClient *IstioClient
}

func webhookServer(istioClient *IstioClient) http.Handler {
	return &webhookHandler{
		istioClient: istioClient,
	}
}

func (h *webhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var aiRequest ai.Request

	aiRequestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", 500)
		return
	}
	r.Body.Close()

	err = json.Unmarshal(aiRequestData, &aiRequest)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Failed to decode request body", 500)
		return
	}

	action := aiRequest.Result.Action
	parameters := aiRequest.Result.Parameters

	var aiResponse *ai.Response

	log.Println("New request for", action)

	switch action {
	case "allowAccess":
		aiResponse, err = allowAccess(parameters, h.istioClient)
	case "denyAccess":
		aiResponse, err = denyAccess(parameters, h.istioClient)
	case "getTopology":
		aiResponse, err = getTopology()
	case "setRoute":
		aiResponse, err = setRoute(parameters, h.istioClient)
	case "getRoute":
		aiResponse, err = getRoute(parameters, h.istioClient)
	}

	if err != nil {
		log.Printf("Failed to perform action %s: %v", action, err)
		http.Error(w, "Failed to perform action", 500)
		return
	}

	out, err := json.MarshalIndent(aiResponse, "", "  ")
	if err != nil {
		log.Printf("Failed to generate response: %v", err)
		http.Error(w, "unable to marshal response", 500)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.Write(out)
}

func allowAccess(params map[string]string, istioClient *IstioClient) (*ai.Response, error) {
	message := "setAccess function was called"
	return &ai.Response{
		DisplayText: message,
		Speech:      message,
		Source:      "Istio Action",
	}, nil
}

func denyAccess(params map[string]string, istioClient *IstioClient) (*ai.Response, error) {
	to := params["to"]
	from := params["from"]
	err := istioClient.DenyAccess(to, from)
	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf("Access to the %s service is prohibited from the %s service.", to, from)
	return &ai.Response{
		DisplayText: message,
		Speech:      message,
		Source:      "Istio Action",
	}, nil
}

func getTopology() (*ai.Response, error) {
	message := "getTopology function was called"
	return &ai.Response{
		DisplayText: message,
		Speech:      message,
		Source:      "Istio Action",
	}, nil
}

func setRoute(params map[string]string, istioClient *IstioClient) (*ai.Response, error) {
	message := "setRoute function was called"
	return &ai.Response{
		DisplayText: message,
		Speech:      message,
		Source:      "Istio Action",
	}, nil
}

func getRoute(params map[string]string, istioClient *IstioClient) (*ai.Response, error) {
	name := params["serviceName"]
	routeRule, err := istioClient.GetRouteRule(name)
	if err != nil {
		return nil, err
	}

	retryCount := routeRule.Spec.HttpReqRetries.SimpleRetry.Attempts

	message := fmt.Sprintf("The %s route has HTTP retries set to %d", name, retryCount)
	return &ai.Response{
		DisplayText: message,
		Speech:      message,
		Source:      "Istio Action",
	}, nil
}
