package getDoc

import (
	"app/ext"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func GetPath(dir string, mode string, q string, meta *ext.Meta) (string, error) {
	dir = strings.TrimSuffix(dir, "/")

	docPath := fmt.Sprintf(
		"%v/%v/%v/%v/%v.xml",
		dir,
		(*meta).Main[q].Date.Year,
		(*meta).Main[q].Date.Month,
		(*meta).Main[q].Date.Day,
		q,
	)
	return docPath, nil
}

func GetRaw(dir string, mode string, q string, meta *ext.Meta) (string, error) {
	if mode == "id" {
		id, _ := strconv.Atoi(q)
		if id >= len((*meta).DocNos) {
			return "", fmt.Errorf("ID not found: %v", id)
		}
		q = (*meta).DocNos[id]
	} else if _, ok := (*meta).Main[q]; !ok {
		return "", errors.New("docno not found: " + q)
	}
	docPath, err := GetPath(dir, "docno", q, meta)
	if err != nil {
		return "", err
	}
	// get raw doc
	b, err := ioutil.ReadFile(docPath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func GetBody(dir string, mode string, q string, meta *ext.Meta) (string, error) {
	if mode == "id" {
		id, _ := strconv.Atoi(q)
		if id >= len((*meta).DocNos) {
			return "", fmt.Errorf("ID not found: %v", id)
		}
		q = (*meta).DocNos[id]
	} else if _, ok := (*meta).Main[q]; !ok {
		return "", errors.New("docno not found: " + q)
	}
	// get raw doc
	raw, err := GetRaw(dir, "docno", q, meta)
	if err != nil {
		return "", err
	}
	var docXML ext.DocTree
	xml.Unmarshal([]byte(raw), &docXML)
	re := regexp.MustCompile(`\n`)
	textContent := strings.Join(docXML.TEXT, "") +
		strings.Join(docXML.TABLE, "") +
		strings.Join(docXML.GRAPHIC, "")
	textContent = re.ReplaceAllString(textContent, "")

	return textContent, nil
}
