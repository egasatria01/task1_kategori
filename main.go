package main

import (
	"categories-sesi-2/database"
	"categories-sesi-2/handlers"
	"categories-sesi-2/repositories"
	"categories-sesi-2/services"
	"encoding/json"
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

	// .env hanya untuk lokal
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// ================= VALIDASI =================
	if config.DBConn == "" {
		log.Fatal("❌ DB_CONN kosong! Pastikan sudah diset di Leapcell ENV")
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	// ================= DATABASE =================
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("❌ Gagal koneksi ke database:", err)
	}
	defer db.Close()

	// ================= DEPENDENCY =================
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// ================= ROUTES =================
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoriesByID)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
		})
	})

	addr := ":" + config.Port
	log.Println("🚀 Server running on port", config.Port)

	log.Fatal(http.ListenAndServe(addr, nil))
}
