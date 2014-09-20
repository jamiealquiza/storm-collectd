// The MIT License (MIT)
//
// Copyright (c) 2014 Jamie Alquiza
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
)

var (
	stormServer = "localhost"
	stormPort   = "8080"
	stormApi    = "/api/v1"
)

func main() {
	topologies := getTop()
	for i := range topologies {
		getTopInfo(topologies[i])
	}
}

func getTop() []string {
	topologies := make([]string, 0)
	// Query API
	resp, err := http.Get("http://" + stormServer + ":" + stormPort + stormApi + "/topology/summary")
	if err != nil {
		fmt.Println("Connection Failed:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// Decode response json
	var data map[string]interface{}
	json.Unmarshal(body, &data)
	topes := data["topologies"].([]interface{})

	for i := range topes {
		topologies = append(topologies, topes[i].(map[string]interface{})["id"].(string))
	}
	return topologies
}

func getTopInfo(endpoint string) {
	// Query API
	resp, err := http.Get("http://" + stormServer + ":" + stormPort + stormApi + "/topology/" + endpoint)
	if err != nil {
		fmt.Println("Connection Failed:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// Decode response json
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	// Parse decoded and populate new map
	// Most values from the Storm API are float64, some are strings.
	// We create a new map 'info' that holds a nested map of strings for each boltId.
	// Each boltId map is populated with the respective bolt's metrics. We collect
	// all metrics with value types of float64 (converting them to strings) and selectively
	// capture some string based metric values that represent meaningful statistics (e.g.
	// bolt 'capacity') while discarding others (e.g. bolt 'lastError').
	r := regexp.MustCompile("capacity|executeLatency|processLatency")
	info := make(map[string]map[string]string)
	// Grab 'bolts' object from full topology info output
	bolts := data["bolts"].([]interface{})
	for i := range bolts {
		// includes top name: // id := endpoint + ":" + bolts[i].(map[string]interface{})["boltId"].(string)
		id := bolts[i].(map[string]interface{})["boltId"].(string)
		m, ok := info[id]
		if !ok {
			m = make(map[string]string)
			info[id] = m
		}
		for k, v := range bolts[i].(map[string]interface{}) {
			if reflect.ValueOf(v).Kind() == reflect.Float64 {
				m[k] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
			} else if r.MatchString(k) {
				m[k] = v.(string)
			}
		}
	}

	// Format output for Collectd
	var hostname = os.Getenv("COLLECTD_HOSTNAME")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}

	for bolt := range info {
		for k, v := range info[bolt] {
			fmt.Printf("PUTVAL %s/storm/gauge-%s-%s N:%s\n", hostname, bolt, k, v)
		}
	}
}
