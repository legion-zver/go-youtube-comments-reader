package main

import (
    "os"
    "fmt"
    "time"
    "strings"
    "io/ioutil"
    "math/rand"    
    "encoding/json"    
    "youtube-haiters/youtube"    
)

// CommentInfo - информация о окмментариях
type CommentInfo struct {    
    Text    string      `json:"text"`
    Date    time.Time   `json:"date"`
}

// AuthorInfo - информация об авторе
type AuthorInfo struct {
    DisplayName     string `json:"displayName"`
    ProfileImageURL string `json:"profileImageURL"`
    ChannelURL      string `json:"channelURL"`
}

func main()  {    
    videoID := ""
    apiKey  := ""
    args := os.Args[1:]
    for _, arg := range args {
        if pos := strings.Index(arg, "-v="); pos >= 0 {
            videoID = arg[pos+3:]
        } else if pos := strings.Index(arg, "-api-key="); pos >= 0 {
            apiKey = arg[pos+9:]
        }
    }        
    if len(videoID) > 0 && len(apiKey) > 0 {
        result, err := youtube.GetListComments(apiKey, videoID, "")
        if err != nil {
            fmt.Println("Error:", err)
            return
        }    
        var authors = make(map[string]*AuthorInfo)
        var comments = make(map[string][]*CommentInfo)
        currentPage := "<start page>"
        for len(result.NextPageToken) > 1 {
            fmt.Println("Get page: ", currentPage)
            currentPage = result.NextPageToken
            fmt.Println("Page Info: ", result.PageInfo)
            for _, item := range result.Items {
                ID := item.Snippet.TopLevelComment.Content.AuthorChannelID.Value
                if authors[ID] == nil {
                    authors[ID] = &AuthorInfo{DisplayName: item.Snippet.TopLevelComment.Content.AuthorDisplayName,
                        ProfileImageURL: item.Snippet.TopLevelComment.Content.AuthorProfileImageURL,
                        ChannelURL: item.Snippet.TopLevelComment.Content.AuthorChannelURL}
                }
                if comments[ID] == nil {
                    comments[ID] = make([]*CommentInfo, 0)
                }
                comments[ID] = append(comments[ID], &CommentInfo{
                    Text: item.Snippet.TopLevelComment.Content.TextDisplay,
                    Date: item.Snippet.TopLevelComment.Content.UpdatedAt,
                })
            }
            time.Sleep(250*time.Millisecond + (time.Millisecond * time.Duration(rand.Int()%850)))
            result, err = youtube.GetListComments(apiKey, videoID, result.NextPageToken)
            if err != nil {
                fmt.Println("Error:", err)
                return
            }            
        }
        // Сохраняем
        var r = make(map[string]interface{})
        r["authors"] = &authors
        r["comments"] = &comments
        b, err := json.Marshal(&r)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        err = ioutil.WriteFile("./commentsByVideo-"+videoID+".json", b, 0755)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        fmt.Println("Read complete, please see file commentsByVideo-"+videoID+".json")
        return
    }
    fmt.Println("Help:")
    fmt.Println("--v=... - youtube video id")
    fmt.Println("--api-key=... - Google API Key")
}