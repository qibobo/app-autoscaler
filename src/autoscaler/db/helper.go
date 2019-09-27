package db

import (
	"net/url"
	"strings"
)

func DetectDBDriver(dbUrl string) (string, error) {
	url, err := url.Parse(dbUrl)
	if err != nil {
		return "", nil
	}
	if strings.ToLower(url.Scheme) == "postgres" {
		return PostgresDriverName, nil
	} else {
		return MysqlDriverName, nil
	}
	// else {
	// 	return "", fmt.Errorf("failed to detect database type from url: %s", dbUrl)
	// }
}
