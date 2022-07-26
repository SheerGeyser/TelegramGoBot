package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

//точка входа
func main() {
	// https://api.telegram.org/bot<token>/METHOD_NAME
	botApi := "https://api.telegram.org/bot" + getToken() + "/"
	offset := 0
	for {
		updates, err := getUpdates(botApi, offset)
		if err != nil {
			log.Println("Somth went wrong", err.Error())
		}
		for _, update := range updates {
			err = respond(botApi, update)
			offset = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}

func getUpdates(botApi string, offset int) ([]Update, error) {
	resp, err := http.Get(botApi + "getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func respond(botApi string, update Update) error {
	var botMessage BotMessage
	var botSendPhoto BotSendPhoto
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = update.Message.Text

	var chatIdString string = strconv.Itoa(update.Message.Chat.ChatId)
	photoBytes, err := ioutil.ReadFile(chatIdString + ".png")
	if err != nil {
		panic(err)
	}
	botSendPhoto.ChatId = update.Message.Chat.ChatId
	botSendPhoto.Photo = string(photoBytes)

	bufMsg, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	postAnsMsg, err := http.Post(botApi+"sendMessage", "application/json", bytes.NewBuffer(bufMsg))
	if err != nil {
		log.Println("post status:" + postAnsMsg.Status)
		return err
	}

	bufImg, err := json.Marshal(botSendPhoto)
	if err != nil {
		return err
	}
	postAnsImg, err := http.Post(botApi+"sendPhoto", "application/json", bytes.NewBuffer(bufImg))
	fmt.Println(bytes.NewBuffer(bufImg))
	if err != nil {
		log.Println("post status:" + postAnsImg.Status)
		return err
	}
	return nil
}

func createBarcode(chatId string) error {
	// Create the barcode
	qrCode, _ := qr.Encode(chatId, qr.M, qr.Auto)

	// Scale the barcode to 200x200 pixels
	qrCode, _ = barcode.Scale(qrCode, 200, 200)

	// create the output file
	name := chatId + ".png"

	file, err := os.Create(name)
	defer file.Close()

	// encode the barcode as png
	png.Encode(file, qrCode)

	return err

}
