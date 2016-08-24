package youtube

import (
    "time"    
    "net/http"
    "encoding/json"
)

// PageInfo - данные о странице
type PageInfo struct {
    TotalResults int `json:"totalResults"`
    ResultsPerPage int `json:"resultsPerPage"`
}

// ChannelID - 
type ChannelID struct {
    Value string `json:"value"`
}

// ContentCommentSnippet -
type ContentCommentSnippet struct {    
    AuthorDisplayName     string      `json:"authorDisplayName"`
    AuthorProfileImageURL string      `json:"authorProfileImageUrl"`
    AuthorChannelURL      string      `json:"authorChannelUrl"`
    AuthorChannelID       ChannelID   `json:"authorChannelId"`
    VideoID               string      `json:"videoId"`
    TextDisplay           string      `json:"textDisplay"`
    CanRate               bool        `json:"canRate"`
    ViewerRating          string      `json:"viewerRating"`
    LikeCount             int         `json:"likeCount"`
    PublishedAt           time.Time   `json:"publishedAt"`
    UpdatedAt             time.Time   `json:"updatedAt"`
}

// TopLevelComment - 
type TopLevelComment struct {
    Kind            string                `json:"kind"`
    Etag            string                `json:"etag"`
    ID              string                `json:"id"`    
    Content         ContentCommentSnippet `json:"snippet"`
}

// CommentSnippet -
type CommentSnippet struct {
    VideoID            string          `json:"videoId"`
    TopLevelComment    TopLevelComment `json:"topLevelComment"`
}

// CommentItem - комментарий
type CommentItem struct {
    Kind            string         `json:"kind"`
    Etag            string         `json:"etag"`
    ID              string         `json:"id"`
    CanReply        bool           `json:"canReply"`
    TotalReplyCount int            `json:"totalReplyCount"`
    IsPublic        bool           `json:"isPublic"`
    Snippet         CommentSnippet `json:"snippet"`
}

// ResponseListComments - ответ от сервера по комментариям
type ResponseListComments struct {
    Kind          string        `json:"kind"`
    Etag          string        `json:"etag"`
    NextPageToken string        `json:"nextPageToken"`
    PageInfo      PageInfo      `json:"pageInfo"`
    Items        []CommentItem  `json:"items"`
}

// GetListComments - получить список комментриев к видео
func GetListComments(apiKey string, videoID string, pageToken string) (resp *ResponseListComments, err error) {
    URL := "https://www.googleapis.com/youtube/v3/commentThreads?part=snippet%2Creplies&videoId="+videoID
    if len(pageToken) > 0 {
        URL += "&pageToken="+pageToken
    }
    URL += "&key="+apiKey
    URL += "&maxResults=100"
    
    r, e := http.Get(URL)
    if e != nil {
        err = e
        return
    }
    resp = new(ResponseListComments)
    decoder := json.NewDecoder(r.Body)
    err = decoder.Decode(resp)
    return
}