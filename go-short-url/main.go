package main

import (
	"fmt"
	"net/http"
	"time"
   "errors"
	"github.com/gin-gonic/gin"
)

type URL struct {
   ID          string     `json:"id"`
   OriginalURL string   `json:"originalId"`
   ShortURL    string   `json:"shortUrl"`
   CreatedAt   time.Time`json:"createdAt"`
   ExpiresAt   time.Time`json:"expiresAt"`
}

var urls = []URL {
   {  ID: "1",
      OriginalURL: "https://www.youtube.com/watch?v=3tCV6GfuL9U",
      ShortURL: "https://www.youtube.com/dd32f",
      CreatedAt: time.Now(),
      ExpiresAt: time.Now().Add(24 * time.Hour),
   },
   {  ID: "2",
      OriginalURL: "https://www.reddit.com/watch?v=3tCV6GfuL9U",
      ShortURL: "https://www.reddit.com/dd32f",
      CreatedAt: time.Now(),
      ExpiresAt: time.Now().Add(24 * time.Hour),
   },
}

func getSartUrls(context *gin.Context) {
   fmt.Println("get request called")
   context.IndentedJSON(http.StatusOK, urls)
}

func addSartUrls(context *gin.Context) {
   var newUrl URL

   if err:= context.BindJSON(&newUrl); err != nil {
      return
   }

   urls = append(urls, newUrl)
   context.IndentedJSON(http.StatusOK, newUrl)
}

func getSartById(id string) (*URL, error) {
   
   for i, url := range urls {
      if url.ID == id {
         return &urls[i], nil
      }
   }
   return nil, errors.New("url not found")

}

func getSart(context *gin.Context) {
   id := context.Param("id")
   URL, err := getSartById(id)

   if err != nil {
      context.IndentedJSON(http.StatusNotFound,
         gin.H{"message": "url not found"})
      return
   }

   context.IndentedJSON(http.StatusOK, URL)
}

func main() {
   router := gin.Default()
   router.GET("/sartlist", getSartUrls)
   router.GET("/sartbyid/:id", getSart)
   router.POST("/sartlist", addSartUrls)

   router.Run("localhost:9090")
}
