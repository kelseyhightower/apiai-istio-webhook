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

type MixerRule struct {
	Rules []Rule
}

type Rule struct {
	Aspects  []Aspect
	Selector string
}

type Aspect struct {
	Kind string
}

type RouteRule struct {
	Type string
	Name string
	Spec RouteSpec
}

type RouteSpec struct {
	Destination    string
	HttpReqRetries HttpReqRetries
	HttpReqTimeout HttpReqTimeout
	Precedence     int64
	Route          []Route
}

type HttpReqRetries struct {
	SimpleRetry SimpleRetry
}

type HttpReqTimeout struct {
	SimpleTimeout SimpleTimeout
}

type Route struct {
	Tags   map[string]string
	Weight int64
}

type SimpleRetry struct {
	Attempts      int64
	PerTryTimeout string
}

type SimpleTimeout struct {
	Timeout string
}

type Topology struct {
	Nodes map[string]map[string]string
	Edges []Edge
}

type Edge struct {
	Source string
	Target string
	Labels map[string]string
}
