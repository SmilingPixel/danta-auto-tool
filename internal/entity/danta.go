package entity

// Config represents the entire configuration parsed from the TOML file.
// See https://github.com/SmilingPixel/DanXi-Backend/blob/main/public/tmp_wait_for_json_editor.toml
type DantaAppContentConfig struct {
	UserAgent       string                `json:"user_agent" toml:"user_agent"`
	StopWords       []string              `json:"stop_words" toml:"stop_words"`
	ChangeLog      string                `json:"change_log" toml:"change_log"`
	HighlightTagIDs []int                 `json:"highlight_tag_ids" toml:"highlight_tag_ids"`
	Banners         []Banner              `json:"banners" toml:"banners"`
	SemesterStart  map[int]string        `json:"semester_start_date" toml:"semester_start_date"`
	Celebrations    []Celebration         `json:"celebrations" toml:"celebrations"`
	LatestVersion   map[string]string     `json:"latest_version" toml:"latest_version"`
}

// Banner represents a single banner item.
type Banner struct {
	Title  string `json:"title" toml:"title"`
	Action string `json:"action" toml:"action"`
	Button string `json:"button" toml:"button"`
}

// BannerApplication represents a single banner application entry.
type BannerApplication struct {
	Banner
	ApplicantEmail string `json:"applicant_email" toml:"applicant_email"`
}

// Celebration represents a single celebration entry.
type Celebration struct {
	Date  string   `json:"date" toml:"date"`
	Words []string `json:"words" toml:"words"`
}