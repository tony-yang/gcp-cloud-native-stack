package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
	pb "github.com/tony-yang/gcp-cloud-native-stack/frontend/genproto"
)

var (
	templates = template.Must(template.New("").ParseGlob("templates/*.html"))
)

func (f *frontendServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Home handler.")
	products, err := f.getProducts(r.Context())
	if err != nil {
		renderHTTPError(w, r, fmt.Errorf("could not retrieve products: %w", err), http.StatusInternalServerError)
		return
	}

	type productView struct {
		Item  *pb.Product
		Price *pb.Money
	}
	var pv []productView
	for _, p := range products {
		pv = append(pv, productView{
			Item:  p,
			Price: p.GetPriceUsd(),
		})
	}

	if err := templates.ExecuteTemplate(w, "home", map[string]interface{}{
		"session_id":   sessionID(r),
		"products":     pv,
		"banner_color": os.Getenv("BANNER_COLOR"),
	}); err != nil {
		log.Error(err)
		renderHTTPError(w, r, fmt.Errorf("failed to render home: %w", err), http.StatusInternalServerError)
	}
}

func (f *frontendServer) productHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Product handler.")
	id := mux.Vars(r)["id"]
	if id == "" {
		renderHTTPError(w, r, fmt.Errorf("product id not specified"), http.StatusBadRequest)
		return
	}

	p, err := f.getProduct(r.Context(), id)
	if err != nil {
		renderHTTPError(w, r, fmt.Errorf("could not retrieve product: %w", err), http.StatusInternalServerError)
		return
	}

	recommendations, err := f.getRecommendations(r.Context(), sessionID(r), []string{id})
	if err != nil {
		renderHTTPError(w, r, fmt.Errorf("failed to get product recommendations: %w", err), http.StatusInternalServerError)
		return
	}

	product := struct {
		Item  *pb.Product
		Price *pb.Money
	}{p, p.GetPriceUsd()}

	if err := templates.ExecuteTemplate(w, "product", map[string]interface{}{
		"session_id":      sessionID(r),
		"product":         product,
		"recommendations": recommendations,
	}); err != nil {
		log.Error(err)
		renderHTTPError(w, r, fmt.Errorf("failed to render product: %w", err), http.StatusInternalServerError)
	}
}

func renderHTTPError(w http.ResponseWriter, r *http.Request, err error, code int) {
	errMsg := fmt.Sprintf("%+v", err)
	w.WriteHeader(code)
	templates.ExecuteTemplate(w, "error", map[string]interface{}{
		"session_id": sessionID(r),
		"error":      errMsg,
		"statu_code": code,
		"status":     http.StatusText(code),
	})
}

func sessionID(r *http.Request) string {
	v := r.Context().Value(ctxKeySessionID{})
	if v != nil {
		return v.(string)
	}
	return ""
}
