package handlers

import (
	"DIVAYTHGRAM_BACKEND/internal/database"
	"DIVAYTHGRAM_BACKEND/internal/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func InitHandlers(r *mux.Router) {
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/authorization", authorization).Methods("GET")
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/stories", getStories).Methods("GET")
	r.HandleFunc("/stories", createStories).Methods("POST")
	r.HandleFunc("/avatars", getAva).Methods("GET")
	r.HandleFunc("/likesDislike", setLikeOrDislike).Methods("POST")
	r.HandleFunc("/likesDislikes", getLikeDislikes).Methods("GET")
	r.HandleFunc("/emojis", setEmoji).Methods("POST")
}

func checkUser(r *http.Request) bool {
	login, password, ok := r.BasicAuth()
	if !ok {
		return false
	}
	var user models.User
	err := database.GetDB().First(&user, "login = ?", login).Error
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func createUser(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	var user = models.User{Login: login, Password: string(hashedPassword)}
	err := database.GetDB().Create(&user).Error
	if err != nil {
		log.Println("User already exists")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	file, handler, _ := r.FormFile("ava")
	defer file.Close()
	f, _ := os.OpenFile("./assets/portrets/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	io.Copy(f, file)
	err = database.GetDB().Create(&models.UserAva{Login: login, Ava: f.Name()}).Error
	if err != nil {
		log.Println("Error with save Picture")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	log.Println("User " + login + "save")
	w.WriteHeader(http.StatusOK)
}

func authorization(w http.ResponseWriter, r *http.Request) {
	if checkUser(r) == true {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var userAva []models.UserAva
	var avaOfUsers []models.AvaOfUser
	_ = database.GetDB().Find(&userAva)
	for i := 0; i < len(userAva); i++ {
		fileBytes, _ := ioutil.ReadFile(userAva[i].Ava)
		avaOfUsers = append(avaOfUsers, models.AvaOfUser{Login: userAva[i].Login, Picture: fileBytes})
	}
	json.NewEncoder(w).Encode(avaOfUsers)
}

func getStories(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["login"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}
	login := keys[0]
	var stories []models.Stories
	var storiesResponse []models.StoriesResponse
	_ = database.GetDB().Find(&stories, "login = ?", login)
	for i := 0; i < len(stories); i++ {
		fileBytes, _ := ioutil.ReadFile(stories[i].Story)
		if !stories[i].Text.Valid {
			storiesResponse = append(storiesResponse, models.StoriesResponse{Id: stories[i].Id, Story: fileBytes, Like: stories[i].Like, Dislike: stories[i].Dislike,
				EmojiHaha: stories[i].EmojiHaha, EmojiZzz: stories[i].EmojiZzz, EmojiWho: stories[i].EmojiWho, EmojiRvota: stories[i].EmojiRvota})
		} else {
			storiesResponse = append(storiesResponse, models.StoriesResponse{Id: stories[i].Id, Story: fileBytes, Like: stories[i].Like, Dislike: stories[i].Dislike,
				EmojiHaha: stories[i].EmojiHaha, EmojiZzz: stories[i].EmojiZzz, EmojiWho: stories[i].EmojiWho,
				EmojiRvota: stories[i].EmojiRvota, Text: stories[i].Text.String, Y: stories[i].Y.String, X: stories[i].X.String,
				Color: stories[i].Color.String, Style: stories[i].Style.String})
		}
	}
	json.NewEncoder(w).Encode(storiesResponse)
}

func createStories(w http.ResponseWriter, r *http.Request) {
	if checkUser(r) == true {
		file, handler, _ := r.FormFile("picture")
		defer file.Close()
		login := r.FormValue("login")
		text := r.FormValue("text")
		color := r.FormValue("color")
		style := r.FormValue("style")
		y := r.FormValue("y")
		x := r.FormValue("x")
		f, _ := os.OpenFile("./assets/stories/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		defer f.Close()
		io.Copy(f, file)
		err := database.GetDB().Create(&models.Stories{Id: models.NewStories().Id, Login: login, Story: f.Name(), Text: models.NewNullString(text), Y: models.NewNullString(y),
			X: models.NewNullString(x), Color: models.NewNullString(color), Style: models.NewNullString(style)}).Error
		if err != nil {
			log.Println("Error with save Picture")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func getAva(w http.ResponseWriter, r *http.Request) {
	if checkUser(r) == true {
		keys, ok := r.URL.Query()["login"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}
		login := keys[0]
		var userAva models.UserAva
		err := database.GetDB().First(&userAva, "login = ?", login).Error
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		fileBytes, _ := ioutil.ReadFile(userAva.Ava)
		w.Write(fileBytes)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func setLikeOrDislike(w http.ResponseWriter, r *http.Request) {
	if checkUser(r) == true {
		var storiesLikeDislike models.StoriesLikeDislike
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&storiesLikeDislike)
		e := database.GetDB().Create(&storiesLikeDislike).Error
		if e != nil {
			log.Println("Error with save Like or Dislike")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		var stories models.Stories
		database.GetDB().First(&stories, "Id = ?", storiesLikeDislike.Id)
		if storiesLikeDislike.Like == true {
			val := stories.Like + 1
			database.GetDB().Model(&stories).Where("Id = ?", storiesLikeDislike.Id).Update("Like", val)
		}
		if storiesLikeDislike.Dislike == true {
			val := stories.Dislike + 1
			database.GetDB().Model(&stories).Where("Id = ?", storiesLikeDislike.Id).Update("Dislike", val)
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func getLikeDislikes(w http.ResponseWriter, r *http.Request) {
	if checkUser(r) == true {
		var storiesLikeDislike []models.StoriesLikeDislike
		login, _, _ := r.BasicAuth()
		err := database.GetDB().Find(&storiesLikeDislike, "login = ?", login).Error
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(storiesLikeDislike)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func setEmoji(w http.ResponseWriter, r *http.Request) {
	if checkUser(r) == true {
		var storesEmoji models.StoriesEmoji
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&storesEmoji)
		e := database.GetDB().Create(&storesEmoji).Error
		if e != nil {
			log.Println("Error with save Emoji")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		var stories models.Stories
		database.GetDB().First(&stories, "Id = ?", storesEmoji.Id)
		if storesEmoji.EmojiHaha == true {
			val := stories.EmojiHaha + 1
			database.GetDB().Model(&stories).Where("Id = ?", storesEmoji.Id).Update("EmojiHaha", val)
		}
		if storesEmoji.EmojiZzz == true {
			val := stories.EmojiZzz + 1
			database.GetDB().Model(&stories).Where("Id = ?", storesEmoji.Id).Update("EmojiZzz", val)
		}
		if storesEmoji.EmojiWho == true {
			val := stories.EmojiWho + 1
			database.GetDB().Model(&stories).Where("Id = ?", storesEmoji.Id).Update("EmojiWho", val)
		}
		if storesEmoji.EmojiRvota == true {
			val := stories.EmojiRvota + 1
			database.GetDB().Model(&stories).Where("Id = ?", storesEmoji.Id).Update("EmojiRvota", val)
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
