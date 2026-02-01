package main

import (
	"categories-sesi-2/database"
	"categories-sesi-2/handlers"
	"categories-sesi-2/repositories"
	"categories-sesi-2/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

// var category = []Category{
// 	{ID: 1, Name: "Indomie goreng", Description: "mie goreng favorite semua orang"},
// 	{ID: 2, Name: "Susu Ultra", Description: "susu UHT rasa coklat"},
// }

// func getCategoryID(w http.ResponseWriter, r *http.Request) {
// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
// 	//ganti jadi string
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "invalid kategori ID", http.StatusBadRequest)
// 		return
// 	}

// 	for _, c := range category {
// 		if c.ID == id {
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(c)
// 			return
// 		}
// 	}
// }

// func updateCategory(w http.ResponseWriter, r *http.Request) {
// 	//GET id dari request
// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
// 	//ganti jadi int
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "invalid kategori ID", http.StatusBadRequest)
// 		return
// 	}

// 	var updateCategory Category
// 	err = json.NewDecoder(r.Body).Decode(&updateCategory)
// 	if err != nil {
// 		http.Error(w, "invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	for i := range category {
// 		if category[i].ID == id {
// 			updateCategory.ID = id
// 			category[i] = updateCategory
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(updateCategory)
// 			return
// 		}
// 	}
// }

// func deleteCategory(w http.ResponseWriter, r *http.Request) {
// 	//GET id
// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
// 	//ganti jadi int
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "invalid kategori ID", http.StatusBadRequest)
// 		return
// 	}
// 	//loop category cari ID, dapet index yang mau dihapus
// 	for i := range category {
// 		if category[i].ID == id {
// 			category = append(category[:i], category[i+1:]...)
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(map[string]string{
// 				"message": "kategori berhasil dihapus",
// 			})
// 			return
// 		}
// 	}
// 	http.Error(w, "kategori tidak ditemukan", http.StatusBadRequest)
// }

func main() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// ----------- SETUP DATABASE ------------
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	defer db.Close()

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	//setup routes
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoriesByID)

	//lolcahost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Sedang Berjalan",
		})
	})
	if config.Port == "" {
		config.Port = "8080"
	}

	addr := ":" + config.Port
	fmt.Println("Server running di", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println("DB_CONN:", config.DBConn)
	}
}
