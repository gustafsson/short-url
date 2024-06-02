package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"salbo.ai/short-url/service"

	"github.com/gorilla/mux"
)

type ShortenRequest struct {
	LongURL string `json:"longUrl"`
}

type ShortenResponse struct {
	ShortURL string `json:"shortUrl"`
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	router.HandleFunc("/", redirectQuery).Methods("GET")
	router.HandleFunc("/shorten", shortenURL).Methods("POST")
	router.HandleFunc("/{id}", redirectPath).Methods("GET")
	router.HandleFunc("/{id}/", redirectPath).Methods("GET")
	router.HandleFunc("/{id}/qr", generateQRCode).Methods("GET")
	router.ServeHTTP(w, r)
}

func redirectQuery(w http.ResponseWriter, r *http.Request) {
	redirectId(w, r, r.URL.RawQuery)
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := service.ShortenURL(req.LongURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ShortenResponse{ShortURL: shortURL})
}

func redirectPath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	redirectId(w, r, id)
}

func redirectId(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		http.NotFound(w, r)
		return
	}

	link, file, err := service.GetRedirect(id)
	req := headersToMap(r.Header)
	if err != nil {
		req["_ResponseErr"] = err.Error()
		log.Panicf("%s", err.Error())
		http.Error(w, "Det där gick visst åt pipsvängen.", http.StatusInternalServerError)
	} else if file != nil {
		req["_ResponseFile"] = string(file)
		w.Header().Set("Content-Type", "text/vcard")
		w.Header().Set("Content-Disposition", "attachment; filename=\"contact.vcf\"")
		w.Write(file)
	} else if link != "" {
		req["_ResponseLink"] = link
		http.Redirect(w, r, link, http.StatusFound)
	} else {
		req["_ResponseNotFound"] = true
		http.NotFound(w, r)
	}
	req["_id"] = id
	service.SaveRequest(req)
}

func generateQRCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	opt := service.DefaultQRCodeOptions()
	if err := populateStructFromQuery(r.URL.Query(), &opt); err != nil {
		http.Error(w, "Error parsing query parameters: "+err.Error(), http.StatusBadRequest)
		return
	}

	qrCode, err := service.GenerateQRCode(service.ShortUrl(id), opt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(qrCode)
}
