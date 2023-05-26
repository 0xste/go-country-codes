package countrycodes

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tchap/go-patricia/patricia"
)

type CountryCode struct {
	Name       string
	ShortName  string
	Alpha2     string
	Alpha3     string
	Numeric    int
	Assignment Assignment
}

type Client struct {
	byAlpha2  map[string]CountryCode
	byName    map[string]CountryCode
	byAlpha3  map[string]CountryCode
	byNumeric map[int]CountryCode
	nameTrie  *patricia.Trie
}

const (
	codesFile = "./codes.csv"
)

func NewClient() (*Client, error) {
	codesFileBytes, err := os.ReadFile(codesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read csv error: %s", err.Error())
	}

	reader := csv.NewReader(bytes.NewReader(codesFileBytes))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv error: %s", err.Error())
	}

	//full_name,short_name_lower_case,remarks,independent_bool,territory_name,alpha2_code,alpha3_code,numeric_code,assignment_status
	var alpha2 = make(map[string]CountryCode)
	if len(rows) <= 1 {
		return nil, fmt.Errorf("no csv rows found: %s", err.Error())
	}
	for i, row := range rows[1:] {
		if len(row) != 9 {
			return nil, fmt.Errorf("csv formatting error at line %d error: %s", i, err.Error())
		}
		assignment, err := NewAssignment(row[8])
		if err != nil {
			return nil, err
		}
		numeric, err := strconv.Atoi(row[7])
		if err != nil {
			return nil, err
		}
		alpha2[row[5]] = CountryCode{
			Name:       row[0],
			ShortName:  row[1],
			Alpha2:     row[5],
			Alpha3:     row[6],
			Numeric:    numeric,
			Assignment: assignment,
		}
	}

	if err := validate(alpha2); err != nil {
		return nil, err
	}

	client := &Client{
		byAlpha2:  alpha2,
		byAlpha3:  make(map[string]CountryCode),
		byName:    make(map[string]CountryCode),
		byNumeric: make(map[int]CountryCode),
		nameTrie:  patricia.NewTrie(),
	}

	for _, cc := range alpha2 {
		if cc.Alpha3 != "" {
			client.byAlpha3[cc.Alpha3] = cc
		}
		client.byName[cc.Name] = cc
		client.byNumeric[cc.Numeric] = cc
		client.nameTrie.Insert(patricia.Prefix(strings.ToLower(cc.Name)), cc)
	}

	return client, nil
}

func validate(alpha2 map[string]CountryCode) error {
	var errorFields = make(map[string][]string)
	for code, detail := range alpha2 {
		if code == "" {
			errorFields[code] = append(errorFields[code], "nil")
			continue
		}
		if code != detail.Alpha2 {
			errorFields[code] = append(errorFields[code], "code and alpha2 mismatch")
			continue
		}
		if !detail.Assignment.Valid() {
			errorFields[code] = append(errorFields[code], fmt.Sprintf("invalid assignment '%d'", detail.Assignment))
			continue
		}
		if len(detail.Alpha2) != 2 {
			errorFields[code] = append(errorFields[code], fmt.Sprintf("invalid alpha2 length '%s'", detail.Alpha2))
			continue
		}
		if len(detail.Alpha3) != 3 {
			errorFields[code] = append(errorFields[code], fmt.Sprintf("invalid alpha3 length '%s'", detail.Alpha3))
			continue
		}
	}

	if len(errorFields) != 0 {
		return fmt.Errorf("invalid country codes configuration: %v", errorFields)
	}

	return nil
}

func (cc *Client) GetByAlpha2(a2 string) (CountryCode, bool) {
	code := cc.byAlpha2[a2]
	return code, code.Alpha2 != ""
}

func (cc *Client) GetByAlpha3(a3 string) (CountryCode, bool) {
	code := cc.byAlpha3[a3]
	return code, code.Alpha2 != ""
}

func (cc *Client) GetByName(name string) (CountryCode, bool) {
	code := cc.byName[name]
	return code, code.Alpha2 != ""
}

func (cc *Client) GetByNumeric(numeric int) (CountryCode, bool) {
	code := cc.byNumeric[numeric]
	return code, code.Alpha2 != ""
}

func (cc *Client) FindByName(prefix string) (matches []CountryCode) {
	matches = make([]CountryCode, 0)
	visit := func(prefix patricia.Prefix, item patricia.Item) error {
		matches = append(matches, item.(CountryCode))
		return nil
	}
	_ = cc.nameTrie.VisitSubtree(patricia.Prefix(strings.ToLower(prefix)), visit)
	return
}
