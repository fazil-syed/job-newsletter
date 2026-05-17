package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"jobs-newsletter/internal/config"
	"jobs-newsletter/internal/db"
)

func SendNewsletter(cfg config.ResendConfig, content string) error {
	var subs []db.Subscriber
	db.DB.Find(&subs)

	fmt.Println("total subscribers:", len(subs))

	for _, s := range subs {
		fmt.Println("sending to:", s.Email)

		body := map[string]interface{}{
			"from":    cfg.FromEmail,
			"to":      []string{s.Email},
			"subject": "🔥 High-Paying Remote Jobs",
			"html":    content,
		}

		jsonData, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("failed:", s.Email, err)
			continue
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)

		fmt.Println("status:", resp.Status)
		fmt.Println("response:", string(bodyBytes))
		fmt.Println("finished:", s.Email)
	}

	return nil
}
func BuildEmailContent(cfg config.AppConfig, post db.Post, email string) string {
	baseURL := cfg.BASEUrl
	content := WrapLinks(baseURL, post.Content, email, post.ID)
	trackingPixel := fmt.Sprintf(
		`<img src="%s/open?email=%s&post_id=%d" width="1" height="1"/>`,
		baseURL, email, post.ID,
	)

	content = strings.ReplaceAll(content, "UNSUBSCRIBE_PLACEHOLDER", fmt.Sprintf("%s/unsubscribe?email=%s", baseURL, email))
	return content + "<br><br>" + trackingPixel
}
func SendPost(cfg *config.Config, post db.Post) {
	var subs []db.Subscriber
	db.DB.Find(&subs)

	for _, s := range subs {
		html := BuildEmailContent(cfg.App, post, s.Email)

		body := map[string]interface{}{
			"from":    cfg.Resend.FromEmail,
			"to":      []string{s.Email},
			"subject": post.Title,
			"html":    html,
		}

		jsonData, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+cfg.Resend.APIKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		client.Do(req)
	}
}

func WrapLinks(baseURL, content, email string, postID uint) string {
	re := regexp.MustCompile(`href="([^"]+)"`)

	return re.ReplaceAllStringFunc(content, func(match string) string {
		// extract original URL
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		original := submatches[1]

		encoded := url.QueryEscape(original)

		newURL := fmt.Sprintf(
			`href="%s/click?email=%s&post_id=%d&url=%s"`,
			baseURL,
			email,
			postID,
			encoded,
		)

		return newURL
	})
}
