package handler

import "net/http"

type Handler interface {
	CreateFile(w http.ResponseWriter, r *http.Request)
	GetItem(w http.ResponseWriter, r *http.Request)
	GetItemByUserId(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
	UpdateItem(w http.ResponseWriter, r *http.Request)

	CreateCapsule(w http.ResponseWriter, r *http.Request)
	GetCapsulesByUser(w http.ResponseWriter, r *http.Request)
	UpdateCapsule(w http.ResponseWriter, r *http.Request)
	DeleteCapsule(w http.ResponseWriter, r *http.Request)
	GetCapsule(w http.ResponseWriter, r *http.Request)
	AddItemToCapsule(w http.ResponseWriter, r *http.Request)
	DeleteItemFromCapsule(w http.ResponseWriter, r *http.Request)
	GetItemsFromCapsule(w http.ResponseWriter, r *http.Request)

	GetLookData(w http.ResponseWriter, r *http.Request)
	AddToLook(w http.ResponseWriter, r *http.Request)
	DeleteLook(w http.ResponseWriter, r *http.Request)
	CreateLook(w http.ResponseWriter, r *http.Request)
}
