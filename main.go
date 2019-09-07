package main

// Checked-190818-2145
import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

// Checked-190818-2147
const difficulty = 1

// Checked-190818-2146
type Block struct {
	Index       int
	Timestamp   string
	PayloadData string
	Hash        string
	PrevHash    string

	// PoW Struture
	Difficulty int
	Nonce      string
}

// Checked-190818-2147
var Blockchain []Block
var genesisBlock Block

// Checked-190818-2147
// type Message struct {
// 	PayloadData string
// }

// Checked-190818-2147
var mutex = &sync.Mutex{}

// Checked-190818-2158
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.PayloadData + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Modified-190818-2016
func generateBlock(oldBlock Block, PayloadData string, recieverId string) Block {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PayloadData = PayloadData
	newBlock.PrevHash = oldBlock.Hash

	// PoW Struture
	newBlock.Difficulty = difficulty
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			msg := "Mining: [" + calculateHash(newBlock) + "]"
			bot.PushMessage(recieverId, linebot.NewTextMessage(msg)).Do()
			fmt.Println(msg)
			// time.Sleep(time.Second)
			continue
		} else {
			msg := "Mining succeeded: [" + calculateHash(newBlock) + "]"
			bot.PushMessage(recieverId, linebot.NewTextMessage(msg)).Do()
			fmt.Println(msg)
			newBlock.Hash = calculateHash(newBlock)
			break
		}
	}
	return newBlock
}

// Checked-190818-2156
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// func replaceChain(newBlocks []Block) {
// 	if len(newBlocks) > len(Blockchain) {
// 		Blockchain = newBlocks
// 	}
// }

// Checked-190818-2148
func run() error {
	mux := makeMuxRouter()
	httpAddr := "8080" //"os.Getenv("ADDR")"
	log.Println("Listening on ", "8080" /*os.Getenv("ADDR")*/)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// Checked-190818-2149
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	// muxRouter.HandleFunc("/webhook", handleGetBlockchain).Methods("Get")
	muxRouter.HandleFunc("/webhook", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					usrMsg := message.Text

					tmpMsg := make([]string, 2)
					tmpMsg = strings.Fields(usrMsg)

					recieverId := event.Source.UserID
					if event.Source.GroupID != "" {
						recieverId = event.Source.GroupID
					} else if event.Source.RoomID != "" {
						recieverId = event.Source.RoomID
					}
					// event.Source.UserID,
					// event.Source.GroupID,
					// event.Source.RoomID,
					// recieverId :=

					if strings.EqualFold(tmpMsg[0], ":SEND") {
						handleWriteBlock(w, req, event.ReplyToken, recieverId, tmpMsg[1])
					} else if strings.EqualFold(tmpMsg[0], ":GET") {
						tmpIndex, err := strconv.ParseUint(tmpMsg[1], 10, 64)
						if err != nil {
							log.Fatal(err)
						}
						handleGetBlockchain(w, req, event.ReplyToken, tmpIndex)
					}
					// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					// 	log.Print(err)
					// }
				}
			}
		}
	}).Methods("POST")
	return muxRouter
}

// Checked-190818-2149
// w http.ResponseWriter, r *http.Request,
func handleGetBlockchain(w http.ResponseWriter, r *http.Request, replyToken string, index uint64) {
	msg, err := json.MarshalIndent(Blockchain, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpMsg := "All Block: " + string(msg)
	splitMsg := strings.Split(tmpMsg, "}")

	bot.ReplyMessage(replyToken, linebot.NewTextMessage(splitMsg[index])).Do()
	fmt.Println("Line User request All Block content")
	// io.WriteString(w, string(bytes))
}

// Checked-190818-2154
// w http.ResponseWriter, r *http.Request,
func handleWriteBlock(w http.ResponseWriter, r *http.Request, replyToken string, recieverId string, payloadData string) {
	w.Header().Set("Content-Type", "application/json")
	// var m Message

	// decoder := json.NewDecoder(r.Body)
	// if err := decoder.Decode(&m); err != nil {
	// 	respondWithJSON(w, r, http.StatusInternalServerError, r.Body)
	// 	return
	// }
	// defer r.Body.Close()

	// Ensure atomicity when creating new block
	mutex.Lock()
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], payloadData, recieverId)
	mutex.Unlock()

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		Blockchain = append(Blockchain, newBlock)

		msg, err := json.MarshalIndent(newBlock, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpMsg := "New Block: " + string(msg)
		bot.ReplyMessage(replyToken, linebot.NewTextMessage(tmpMsg)).Do()
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

// Checked-190818-2152
func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

// Created-190818-2015
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

var bot *linebot.Client

func main() {
	// Line Message API Authentication
	var bot_t, err = linebot.New(
		"<CHANNEL_SECRET>",
		"<CHANNEL_TOKEN>",
	)
	bot = bot_t
	if err != nil {
		log.Fatal(err)
	}
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		// genesisBlock := Block{}
		genesisBlock := Block{0, t.String(), "INIT. Payload Data", calculateHash(genesisBlock), "", difficulty, ""}
		spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())
}
