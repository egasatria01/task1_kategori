package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id`
	Name        string `json:"name"`
	Description string `json:"deskripsi"`
}

var category = []Category{
	{ID: 1, Name: "Indomie goreng", Description: "mie goreng favorite semua orang"},
	{ID: 2, Name: "Susu Ultra", Description: "susu UHT rasa coklat"},
}

func getCategoryID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	//ganti jadi string
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid kategori ID", http.StatusBadRequest)
		return
	}

	for _, c := range category {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	//GET id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	//ganti jadi int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid kategori ID", http.StatusBadRequest)
		return
	}

	var updateCategory Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	for i := range category {
		if category[i].ID == id {
			updateCategory.ID = id
			category[i] = updateCategory
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	//GET id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	//ganti jadi int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid kategori ID", http.StatusBadRequest)
		return
	}
	//loop category cari ID, dapet index yang mau dihapus
	for i := range category {
		if category[i].ID == id {
			category = append(category[:i], category[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "kategori berhasil dihapus",
			})
			return
		}
	}
	http.Error(w, "kategori tidak ditemukan", http.StatusBadRequest)
}

func main() {

	//GET localhost:8080/api/kategori{id}
	//PUT localhost:8080/api/kategori{id}
	//DELETE localhost:8080/api/kategori{id}
	http.HandleFunc("/api/kategori/", func(w http.ResponseWriter, r *http.Request) {
		//PUT localhost:8080 /api/kategori/
		if r.Method == "GET" {
			getCategoryID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})
	// GET localhost:8080/api/produk
	// POST localhost:8080/api/produk
	http.HandleFunc("/api/kategori", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
		} else if r.Method == "POST" {
			var newCategory Category
			err := json.NewDecoder(r.Body).Decode(&newCategory)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}

			// ini adalah untuk memasukkan data nya kedalam variable category
			newCategory.ID = len(category) + 1
			category = append(category, newCategory)

			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(newCategory)
		}
	})

	//lolcahost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Sedang Berjalan",
		})
	})
	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server gagal running")
	}
}
