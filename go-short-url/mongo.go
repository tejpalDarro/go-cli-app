package main 

import (
    "context"
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "log"
    "regexp"
    "time"
    "math/rand"
    "strings"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    // "go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	client        *mongo.Client
	urlCollection *mongo.Collection
)

type Url struct {
	ID          uint      `bson:"_id,omitempty`
	OriginalUrl string    `bson:"original_url,omitempty`
	ShortUrl    string    `bson:"short_url,omitempty`
	CreatedAt   time.Time `bson:"created_at,omitempty`
	ExpiresAt   time.Time `bson:"expires_at,omitempty`
	HashCode    string    `bson:"hash_code,omitempty`
}

func mongoConnect(str string) (string, error) {
    connectDB()
    if urlCollection != nil {
            fmt.Println("urlCollection initialized succesfully")
    } else {
            fmt.Println("Failed to initialize urlCollection")
    }


    if len(str) == 0 {
            return str , fmt.Errorf("string length 0") 
    }

    simplifiedUrl, err := simplifyString(str)
    if err != nil {
            fmt.Println(err)
            return str,fmt.Errorf("error creating simplified Url") 
    }

    fmt.Println(simplifiedUrl)

    flag := checkHashCodeInDb(&simplifiedUrl)    
    if !flag {
        addUrlToMongo(&simplifiedUrl)
    } else {
        fmt.Println("url already present")
        fmt.Println("getting the url...")
    }


    // addDemoUrl()

    return simplifiedUrl, nil 
}

func generateShortURL() string {
    const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    rand.NewSource(time.Now().UnixNano())
    // fmt.Println("rnadom number: ",randd)
    var shortURL strings.Builder
    for i := 0; i < 6; i++ {
            shortURL.WriteByte(chars[rand.Intn(len(chars))])
    }
    return shortURL.String()
}


func checkHashCodeInDb(url *string) bool {
    checkConnection()

    hashValue := hashString(*url) 

    filter := bson.M{"hashcode": hashValue}
    var result Url
    err := urlCollection.FindOne(context.Background(), filter).Decode(&result)
    if err == mongo.ErrNoDocuments {
        fmt.Println("HashCode not found")
        return false
    } else if err != nil {
        fmt.Println("Error querying MongoDB:", err)
        return true 
    } else {
        fmt.Println("HashCode already exists")
        return true
    }
}

func addUrlToMongo(url *string) {
    checkConnection() 

    //calculate hash of string
    hashValue := hashString(*url)

    //generate random string
    randString := generateShortURL()
    fmt.Println("random String:", randString)

    //concate random string to localhost
    var strBuilder strings.Builder
    strBuilder.WriteString("localhost:8080/")
    strBuilder.WriteString(randString)
    result := strBuilder.String()

    singlurl := Url{
            OriginalUrl: *url,
            ShortUrl:    result,
            CreatedAt:   time.Now(),
            ExpiresAt:   time.Now().AddDate(0, 1, 0),
            HashCode:    hashValue,
    }


    collect, err := urlCollection.InsertOne(context.TODO(), &singlurl)
    if err != nil {
            log.Fatalln("Error Inserting Document", err)
    }

    log.Println("url inserted", collect)

}

func simplifyString(input string) (string, error) {
    // Regular expression pattern to match everything before //
    pattern := `(https?:\/\/)`

    // Compile the regular expression
    re := regexp.MustCompile(pattern)

    // Find the first match
    matches := re.FindStringSubmatch(input)

    // fmt.Println(matches)
    // fmt.Println(remainingString)


    if len(matches) > 1 {
        // If a match is found, select the substring after //
        // newString := matches[1]
        // fmt.Println("Selected substring after //:", newString)

        remainingString := re.ReplaceAllString(input, "")
        return remainingString, nil
    } else {
            fmt.Println("No match found.")
            return input, fmt.Errorf("cannot find //")
    }
}



func addDemoUrl() {
    checkConnection()

    url := "localhost:8080/testit.com"
    hashValue := hashString(url)
    // fmt.Println(hashValue)
    // for i:=0; i<3; i++ {
    //     fmt.Println(hashString(url))
    // }

    singlurl := Url{
            ID:          1,
            OriginalUrl: url,
            ShortUrl:    "localhost:8080/be31dd",
            CreatedAt:   time.Now(),
            ExpiresAt:   time.Now().AddDate(1, 1, 1),
            HashCode:    hashValue,
    }

    collect, err := urlCollection.InsertOne(context.TODO(), &singlurl)
    if err != nil {
            log.Fatalln("Error Inserting Document", err)
    }

    log.Println("demo url inserted", collect)
}

func hashString(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func connectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to the MongoDB server
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Set a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ping the MongoDB server to check the connection
	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}
	fmt.Println("Ping successful!")

	urlCollection = client.Database("sorturl").Collection("urls")
}

func checkConnection() {
	if client == nil {
		fmt.Println("Not connected to MongoDB!")
		return
	}

	// Ping the MongoDB server to check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		fmt.Println("Connection lost!")
	} else {
		fmt.Println("Connection is active!")
	}
}
