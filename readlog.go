package ccruncher

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type CCLog struct {
	entries    map[string][]LogEntry
	requestIDs map[string][]string
}

type LogEntry struct {
	Timestamp  float64 `json:",omitempty"`
	Message    string  `json:",omitempty"`
	LogLevel   string  `json:"log_level,omitempty"`
	Source     string  `json:",omitempty"`
	Data       *LogEntryData
	File       string `json:",omitempty"`
	LineNumber int    `json:"lineno,omitempty"`
	Method     string `json:",omitempty"`
	Request    *Request
}

type LogEntryData struct {
	RequestGUID string `json:"request_guid,omitempty"`
	ProcessGUID string `json:"process_guid,omitempty"`
}

type Request struct {
	RequestID string
	AppGUID   string
	Method    string
	URIPath   string
}

func (e LogEntry) RequestID() string {
	return e.Request.RequestID
}

func (e LogEntry) AppGUID() string {
	return e.Request.AppGUID
}

func (e LogEntry) HttpMethod() string {
	return e.Request.Method
}

func (e LogEntry) URIPath() string {
	return e.Request.URIPath
}

const unspecifiedRequestID = "unspecified"
const unspecifiedAppGUID = "unspecified"
const vcapRequestIDRegex = "vcap-request-id: [\\w\\-\\:]+"
const appGUIDRegex = "\\/v2\\/apps/[\\w\\-]+"
const methodAndURIPathRegex = `Started [A-Z]+ \"[\w\/\-]+`

func (c *CCLog) EntriesForRequest(requestID string) []LogEntry {
	return c.entries[requestID]
}

func (c *CCLog) Entries() []LogEntry {
	var entries []LogEntry

	// e is of type []LogEntry
	for _, e := range c.entries {
		entries = append(entries, e...)
	}

	return entries
}

func (c *CCLog) Apps() []string {
	var guids []string

	for guid, _ := range c.requestIDs {
		guids = append(guids, guid)
	}

	return guids
}

func (c *CCLog) RequestsForApp(appGUID string) []string {
	if ids, exists := c.requestIDs[appGUID]; exists == true {
		return ids
	}

	return nil
}

func ParseLog(logFile io.Reader) (*CCLog, error) {
	var ccLog *CCLog
	var entry LogEntry
	var lineNumber uint

	ccLog = &CCLog{
		entries:    make(map[string][]LogEntry),
		requestIDs: make(map[string][]string),
	}

	scanner := bufio.NewScanner(logFile)

	for scanner.Scan() {
		lineNumber++

		entry.Data = &LogEntryData{}

		err := json.Unmarshal([]byte(scanner.Text()), &entry)

		if err != nil {
			return ccLog, fmt.Errorf("ParseLog() error at line %d: %s", lineNumber, err)
		}

		id, _ := requestIDFromMessage(entry.Message)

		if id == "" {
			if entry.Data.RequestGUID != "" {
				id = entry.Data.RequestGUID
			} else {
				id = unspecifiedRequestID
			}
		}

		if existingRequests := ccLog.EntriesForRequest(id); len(existingRequests) > 0 {
			entry.Request = existingRequests[0].Request
		} else {
			// We should only be here if we've never seen the request ID before
			entry.Request = &Request{}
			entry.Request.RequestID = id
			appGUID, _ := appGUIDFromMessage(entry.Message)
			if appGUID == "" {
				if entry.Data.ProcessGUID != "" {
					appGUID = entry.Data.ProcessGUID
				} else {
					appGUID = unspecifiedAppGUID
				}
			}

			method, uriPath, _ := methodAndURIPathFromMessage(entry.Message)
			entry.Request.AppGUID = appGUID
			entry.Request.Method = method
			entry.Request.URIPath = uriPath

			ccLog.requestIDs[appGUID] = append(ccLog.requestIDs[appGUID], id)
		}

		ccLog.entries[id] = append(ccLog.entries[id], entry)

	}

	return ccLog, nil
}

func requestIDFromMessage(message string) (string, error) {
	rx, err := regexp.Compile(vcapRequestIDRegex)
	if err != nil {
		return "", fmt.Errorf("Regexp compilation failure: %s", err)
	}

	var firstMatch string

	if firstMatch = rx.FindString(message); firstMatch == "" {
		return "", fmt.Errorf("Unable to extract vcap-request-id from %s", message)
	}

	return strings.Split(firstMatch, " ")[1], nil
}

func appGUIDFromMessage(message string) (string, error) {
	rx, err := regexp.Compile(appGUIDRegex)
	if err != nil {
		return "", fmt.Errorf("App Guid Regexp compilation failure: %s", err)
	}

	var firstMatch string
	if firstMatch = rx.FindString(message); firstMatch == "" {
		return "", nil
	}

	return strings.TrimPrefix(firstMatch, "/v2/apps/"), nil
}

func methodAndURIPathFromMessage(message string) (string, string, error) {
	rx, err := regexp.Compile(methodAndURIPathRegex)
	if err != nil {
		return "", "", fmt.Errorf("Method and URI Path Regexp compilation error: %s", err)
	}

	var firstMatch string
	if firstMatch = rx.FindString(message); firstMatch == "" {
		return "", "", nil
	}

	tokens := strings.Split(firstMatch, " ")

	return tokens[1], strings.TrimLeft(tokens[2], `"`), nil
}
