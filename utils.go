package ujeebu

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

const CreditsHeader = "ujb-credits"

func getUjeebuCreditsFromResponse(resp *resty.Response) int {
	for k, v := range resp.RawResponse.Header {
		if strings.EqualFold(k, CreditsHeader) {
			if len(v) > 0 {
				cv := v[0]
				if c, err := strconv.Atoi(cv); err == nil {
					return c
				}
			}
		}
	}
	return 0
}

// encodeBase64 encodes a string to Base64
func encodeBase64(value string) string {
	if value == "" {
		return value
	}
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// shouldEncodeWaitFor checks if WaitFor should be encoded
func shouldEncodeWaitFor(value string) bool {
	if len(value) == 0 {
		return false
	}

	// If it's a number, return false (do not encode)
	if _, err := strconv.Atoi(value); err == nil {
		return false
	}

	// If it's a short string (<= 100 chars), return false
	if len(value) <= 100 {
		return false
	}

	return true
}
