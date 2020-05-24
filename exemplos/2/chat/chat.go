package chat

import (
	"log"
	"net/http"
)

func Start(port string) {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Println("Iniciando Servidor de Chat na Porta:", port)

	rooms := []string{"reactjs", "golang", "rust", "flutter"}
	for _, room := range rooms {
		r := NewRoom(room)
		http.HandleFunc("/chat/"+room, r.Handler)
		go r.Run()
	}

	log.Fatal(http.ListenAndServe(port, nil))
}
