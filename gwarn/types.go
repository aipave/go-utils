package gwarn

type header struct {
    Title    element `json:"title"`
    Template string  `json:"template"`
}

type element struct {
    Tag     string `json:"tag"`
    Content string `json:"content"`
}

type cards struct {
    Header   header    `json:"header"`
    Elements []element `json:"elements"`
}

type msgCard struct {
    MsgType string `json:"msg_type"`
    Card    cards  `json:"card"`
}

type FontColor string

const (
    FontColorGreen FontColor = "green"
    FontColorRed   FontColor = "red"
    FontColorGrey  FontColor = "grey"
)
