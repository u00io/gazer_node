package system

import (
	"net/http"
	"sync"
	"time"

	"github.com/u00io/gomisc/logger"
)

type U00 struct {
	httpClient *http.Client
}

func NewU00() *U00 {
	var c U00
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		IdleConnTimeout:     3 * time.Second,
	}

	c.httpClient = &http.Client{
		Transport: transport,
		Timeout:   3 * time.Second,
	}
	return &c
}

func (c *U00) Run() {
}

type ItemToSet struct {
	Path  string
	Name  string
	Value string
	Uom   string
}

func (c *U00) Write(apiKey string, items []ItemToSet) error {
	parameters := make(map[string][]string)
	for _, item := range items {
		if item.Path == "/" {
			parameters[item.Path+"_name"] = []string{item.Name}
			parameters[item.Path+""] = []string{item.Value}
			parameters[item.Path+"_uom"] = []string{item.Uom}
		} else {
			parameters[item.Path+"/_name"] = []string{item.Name}
			parameters[item.Path+""] = []string{item.Value}
			parameters[item.Path+"/_uom"] = []string{item.Uom}
		}
	}

	/*for k, p := range parameters {
		logger.Println("Write parameter:", k, "=", p)
	}*/

	resp, err := c.httpClient.PostForm(GetServerURL()+"/set/"+apiKey, parameters)
	if err != nil {
		logger.Println("Error sending data:", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}

var servers = []string{
	"https://gazer.cloud",
	//"https://test.u00.io",
}

var mtx sync.Mutex
var currentServerIndex int

func GetServerURL() string {
	mtx.Lock()
	defer mtx.Unlock()
	url := servers[currentServerIndex]
	currentServerIndex++
	if currentServerIndex >= len(servers) {
		currentServerIndex = 0
	}
	return url
}
