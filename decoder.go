package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type OptionCodeList map[string]OptionCode

type OptionCode struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

var OptionCodeListInstance OptionCodeList = nil

func InitOptionCodeMap() (OptionCodeList, error) {
	log.Println("Reading option code map...")
	s, err := ioutil.ReadFile("optioncodes.json")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(s, &OptionCodeListInstance); err != nil {
		return nil, err
	}
	return OptionCodeListInstance, nil
}

func ReadFile(url string) ([]byte, error) {
	if strings.Index(url, "file://") == 0 {
		return ioutil.ReadFile(strings.TrimPrefix(url, "file://"))
	}
	client := http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid status code " + strconv.Itoa(resp.StatusCode))
	}
	contentType := resp.Header.Get("Content-Type")
	if strings.Index(contentType, "text/html") != 0 {
		return nil, errors.New("invalid content type " + contentType)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if len(body) > 1000000 {
		return nil, errors.New("response body too large")
	}
	return body, err
}

func ExtractCodes(file []byte) ([]string, error) {
	re := regexp.MustCompile("\"ManufacturingOptionCodeList\":\"([A-Za-z0-9\\,]+)\"")
	match := re.FindSubmatch(file)
	if len(match) > 0 {
		codes := strings.Split(string(match[1]), ",")
		return codes, nil
	}
	return nil, errors.New("ManufacturingOptionCodeList not found")
}
