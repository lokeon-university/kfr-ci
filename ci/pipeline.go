package main

var (
	supportedImages = map[string]string{
		"cpp":    "gcr.io/kfr-ci/kfr-cpp",
		"go":     "gcr.io/kfr-ci/kfr-go",
		"java":   "gcr.io/kfr-ci/kfr-java",
		"python": "gcr.io/kfr-ci/kfr-python",
	}
)

type pipeline struct {
	Branch      string                               `json:"branch,omitempty"`
	Language    string                               `json:"language,omitempty"`
	LogFileName string                               `json:"log_file_name,omitempty"`
	Owner       string                               `json:"owner,omitempty"`
	Repository  string                               `json:"repository,omitempty"`
	TelegramID  string                               `json:"telegram_id,omitempty"`
	URL         string                               `json:"url,omitempty"`
	Status      func(string, string, string, string) `json:"-,omitempty"`
}

func (p *pipeline) supportedLanguage() (ok bool) {
	_, ok = supportedImages[p.Language]
	return
}

func (p *pipeline) getImage() (image string) {
	image, _ = supportedImages[p.Language]
	return
}

func (p *pipeline) envVars() []string {
	return []string{
		keyValueEnv("REPO_BRANCH", p.Branch),
		keyValueEnv("REPO_NAME", p.Repository),
		keyValueEnv("REPO_URL", p.URL),
	}
}

func keyValueEnv(key, value string) string {
	return key + "=" + value
}

func (p *pipeline) status(status string) {
	p.Status(p.TelegramID, status, p.Repository, p.Owner)
}
