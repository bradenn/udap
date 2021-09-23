package template

type Media struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Thumbnail string `json:"thumbnail"`
	Duration  int    `json:"duration"`
	Id        string `json:"id"`
}

type State struct {
	Playing  bool
	Playback Media
}
