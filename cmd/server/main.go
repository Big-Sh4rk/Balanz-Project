package main

import (
	"context"
	"github.com/Big-Sh4rk/Balanz-Project/internal/dolar"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

//Servidor con Graceful Shutdown
const puerto = ":8080"

type myServer struct {
	http.Server
	shutdownReq chan bool
	reqCount    uint32
}

func NewServer() *myServer {
	//Crear servidor
	s := &myServer{
		Server: http.Server{
			Addr:         puerto,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		shutdownReq: make(chan bool),
	}

	router := mux.NewRouter()

	//set http server handler
	s.Handler = router

	return s
}

func (s *myServer) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	//Espere la solicitud de interrupción o apagado a través de /shutdown
	select {
	case sig := <-irqSig:
		log.Printf("Shutdown request (signal: %v)", sig)
	case sig := <-s.shutdownReq:
		log.Printf("Shutdown request (/shutdown %v)", sig)
	}

	log.Printf("Stoping http server ...")

	//Crear contexto de apagado con 10 segundos de tiempo de espera
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Apagar el servidor
	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}
}

func (s *myServer) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Gorilla MUX!\n"))
}

func (s *myServer) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shutdown server"))

	//No hacer nada si ya se emitió la solicitud de apagado
	//si s.reqCount == 0 entonces se establece en 1, devuelve verdadero de lo contrario falso
	if !atomic.CompareAndSwapUint32(&s.reqCount, 0, 1) {
		log.Printf("Shutdown through API call in progress...")
		return
	}

	go func() {
		s.shutdownReq <- true
	}()
}

func main() {
	//Inicia el servidor
	server := NewServer()

	done := make(chan bool)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Listen and serve: %v", err)
		}
		done <- true
	}()

	//Consumo de la API
	dolar.ConsumeAPI()

	//Consumo del WebSocket
	dolar.ConsumeSocket()

	//esperar apagado
	server.WaitShutdown()

	<-done
	log.Printf("DONE!")
}
