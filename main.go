package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/RailwayTickets/backend-go/controller"
	"github.com/RailwayTickets/backend-go/entity"
	h "github.com/RailwayTickets/backend-go/handler"
)

func main() {
	http.Handle("/register", h.Chain(http.HandlerFunc(registerHandler),
		h.SetContentTypeJSON,
		h.RequiredPost))
	http.Handle("/login", h.Chain(http.HandlerFunc(loginHandler),
		h.SetContentTypeJSON,
		h.RequiredPost))
	http.Handle("/search", h.Chain(http.HandlerFunc(searchHandler),
		h.CheckAndUpdateToken,
		h.SetContentTypeJSON,
		h.RequiredPost))
	http.Handle("/buy", h.Chain(http.HandlerFunc(buyHandler),
		h.CheckAndUpdateToken,
		h.SetContentTypeJSON))
	http.Handle("/return", h.Chain(http.HandlerFunc(returnHandler),
		h.CheckAndUpdateToken,
		h.SetContentTypeJSON))
	http.Handle("/return/valid", h.Chain(http.HandlerFunc(validReturnHandler),
		h.CheckAndUpdateToken,
		h.SetContentTypeJSON))
	http.Handle("/directions", h.Chain(http.HandlerFunc(allDirectionsHandler),
		h.CheckAndUpdateToken,
		h.SetContentTypeJSON))
	http.Handle("/departures", h.Chain(http.HandlerFunc(allDeparturesHandler),
		h.CheckAndUpdateToken,
		h.SetContentTypeJSON))
	http.Handle("/profile", h.Chain(http.HandlerFunc(profileHandler),
		h.CheckAndUpdateToken,
		h.SetContentTypeJSON))
	http.Handle("/profile/tickets", h.Chain(http.HandlerFunc(myTicketsHandler),
		h.CheckAndUpdateToken,
		h.SetContentTypeJSON))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	registrationInfo := new(entity.RegistrationInfo)
	err := json.NewDecoder(r.Body).Decode(registrationInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := registrationInfo.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	creds, err := controller.Register(registrationInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(creds)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	loginInfo := new(entity.LoginInfo)
	err := json.NewDecoder(r.Body).Decode(loginInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := loginInfo.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	creds, err := controller.Login(loginInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(creds)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := new(entity.TicketSearchParams)
	err := json.NewDecoder(r.Body).Decode(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tickets, err := controller.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tickets)
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
	ticketID := r.URL.Query().Get("id")
	ctx := r.Context()
	err := controller.Buy(ctx.Value(h.LoginKey).(string), ticketID)
	if err == controller.AlreadyTaken {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func returnHandler(w http.ResponseWriter, r *http.Request) {
	ticketID := r.URL.Query().Get("id")
	ctx := r.Context()
	err := controller.Return(ctx.Value(h.LoginKey).(string), ticketID)
	if err == controller.AlreadyIssued || err == controller.NotYourTicket {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func validReturnHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tickets, err := controller.ValidReturn(ctx.Value(h.LoginKey).(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tickets)
}

func allDirectionsHandler(w http.ResponseWriter, r *http.Request) {
	directions, err := controller.GetDirections()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(directions)
}

func allDeparturesHandler(w http.ResponseWriter, r *http.Request) {
	departures, err := controller.GetDepartures()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(departures)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		updateProfileHandler(w, r)
	case http.MethodGet:
		getProfileHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := new(entity.User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err = controller.UpdateProfile(ctx.Value(h.LoginKey).(string), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profile, err := controller.GetProfile(ctx.Value(h.LoginKey).(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(profile)
}

func myTicketsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tickets, err := controller.GetMyTickets(ctx.Value(h.LoginKey).(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tickets)
}
