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
	RepositoryID int64        `json:"repository_id,omitempty"`
	URL          string       `json:"url,omitempty"`
	Repository   string       `json:"repository,omitempty"`
	Branch       string       `json:"branch,omitempty"`
	LogFileName  string       `json:"log_file_name,omitempty"`
	Language     string       `json:"language,omitempty"`
	Status       func(string) `json:"-"`
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
		keyValueEnv("REPO_URL", p.URL),
		keyValueEnv("REPO_NAME", p.Repository),
		keyValueEnv("REPO_BRANCH", p.Branch),
	}
}

func keyValueEnv(key, value string) string {
	return key + "=" + value
}
