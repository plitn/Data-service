package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Data-service/internal/model"
	"github.com/Data-service/internal/service/data_service"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type service struct {
	dataService data_service.Service
}

func New(dataService data_service.Service) *service {
	return &service{dataService: dataService}
}

func (s *service) CreateFile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "form/json")

	err := r.ParseMultipartForm(32 << 30)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bucketName := r.MultipartForm.Value["user_id"][0]
	if bucketName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bucketName = fmt.Sprintf("usr%s", bucketName)
	files, ok := r.MultipartForm.File["file"]
	if !ok || len(files) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	itemData, err := s.getItemData(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileInfo := files[0]
	fileReader, err := fileInfo.Open()
	dto := model.CreateFileDTO{
		Name:   s.generateFileName(bucketName, fileInfo.Filename),
		Size:   fileInfo.Size,
		Reader: fileReader,
	}
	itemData.File.Name = dto.Name
	err = s.dataService.CreateFile(ctx, itemData, dto.Name, bucketName, dto.Size, dto.Reader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *service) CutItem(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 30)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	files, ok := r.MultipartForm.File["file"]
	if !ok || len(files) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileInfo := files[0]
	fileReader, err := fileInfo.Open()
	reader, _, err := s.dataService.ClearFile(s.generateFileName("firstCut", fileInfo.Filename), fileReader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytesData, err := ioutil.ReadAll(reader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(bytesData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *service) getItemData(r *http.Request) (model.Item, error) {
	var sizeNumber int
	var userIdNumber int64
	var sizeL, desc *string
	var col int
	userid := r.MultipartForm.Value["user_id"][0]
	if userid == "" {
		return model.Item{}, fmt.Errorf("no user_id provided")
	} else {
		userIdNumber, _ = strconv.ParseInt(userid, 10, 64)
	}
	itemName := r.MultipartForm.Value["name"][0]
	if itemName == "" {
		return model.Item{}, fmt.Errorf("no name provided")
	}
	var category int
	cat := r.MultipartForm.Value["category"][0]
	if cat == "" {
		return model.Item{}, fmt.Errorf("no category provided")
	}
	category, _ = strconv.Atoi(cat)
	size := r.MultipartForm.Value["numeric_size"][0]
	if size != "" {
		sizeNumber, _ = strconv.Atoi(size)
	}
	sizeLetter := r.MultipartForm.Value["size"][0]
	if sizeLetter != "" {
		sizeL = &sizeLetter
	}
	description := r.MultipartForm.Value["description"][0]
	if description != "" {
		desc = &description
	}
	color := r.MultipartForm.Value["color"][0]
	if color == "" {
		return model.Item{}, fmt.Errorf("no color provided")
	} else {
		col, _ = strconv.Atoi(color)
	}
	item := model.Item{
		UserId:      userIdNumber,
		Name:        itemName,
		Category:    category,
		SizeNumber:  &sizeNumber,
		SizeText:    sizeL,
		Description: desc,
		Color:       col,
	}
	return item, nil
}

func (s *service) GetItem(w http.ResponseWriter, r *http.Request) {
	var fileId int64
	fileIdString := r.URL.Query().Get("file_id")
	if fileIdString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		fileId, _ = strconv.ParseInt(fileIdString, 10, 64)
	}

	f, err := s.dataService.GetItem(r.Context(), fileId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *service) GetItemByUserId(w http.ResponseWriter, r *http.Request) {
	var userId int64
	param := r.URL.Query().Get("user_id")
	if param == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		userId, _ = strconv.ParseInt(param, 10, 64)
	}

	f, err := s.dataService.GetItemsByUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *service) CreateCapsule(w http.ResponseWriter, r *http.Request) {
	var req model.Capsule
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.dataService.CreateCapsule(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *service) GetCapsule(w http.ResponseWriter, r *http.Request) {
	var id int64
	param := r.URL.Query().Get("id")
	if param == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		id, _ = strconv.ParseInt(param, 10, 64)
	}
	fmt.Println(id)
	f, err := s.dataService.GetCapsule(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println(f)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *service) AddItemToCapsule(w http.ResponseWriter, r *http.Request) {
	var req model.ItemCapsule
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.dataService.AddItemToCapsule(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *service) DeleteItemFromCapsule(w http.ResponseWriter, r *http.Request) {
	var req model.ItemCapsule
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.dataService.DeleteItemFromCapsule(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *service) GetItemsFromCapsule(w http.ResponseWriter, r *http.Request) {
	var capsuleId int64
	param := r.URL.Query().Get("capsule_id")
	if param == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		capsuleId, _ = strconv.ParseInt(param, 10, 64)
	}

	f, err := s.dataService.GetItemsInCapsule(r.Context(), capsuleId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *service) GetCapsulesByUser(w http.ResponseWriter, r *http.Request) {
	var userId int64
	param := r.URL.Query().Get("user_id")
	if param == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		userId, _ = strconv.ParseInt(param, 10, 64)
	}

	f, err := s.dataService.GetCapsulesByUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *service) UpdateCapsule(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req model.Capsule
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.dataService.UpdateCapsule(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *service) DeleteCapsule(w http.ResponseWriter, r *http.Request) {
	var capsuleId int64
	param := r.URL.Query().Get("id")
	if param == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		capsuleId, _ = strconv.ParseInt(param, 10, 64)
	}

	err := s.dataService.DeleteCapsule(r.Context(), capsuleId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *service) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req model.Item
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.dataService.UpdateItem(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *service) DeleteItem(w http.ResponseWriter, r *http.Request) {
	var fileId int64
	fileIdString := r.URL.Query().Get("id")
	if fileIdString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		fileId, _ = strconv.ParseInt(fileIdString, 10, 64)
	}

	err := s.dataService.DeleteItem(r.Context(), fileId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *service) generateFileName(userId string, fileName string) string {
	timeMs := time.Now().UTC().UnixMilli()
	ext := ".png"
	return userId + "_" + strconv.FormatInt(timeMs, 10) + ext
}

func (s *service) CreateLook(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req model.Look
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.dataService.CreateLook(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *service) DeleteLook(w http.ResponseWriter, r *http.Request) {
	var lookId int64
	lookIdString := r.URL.Query().Get("id")
	if lookIdString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		lookId, _ = strconv.ParseInt(lookIdString, 10, 64)
	}

	err := s.dataService.DeleteLook(r.Context(), lookId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *service) AddToLook(w http.ResponseWriter, r *http.Request) {
	var lookId, itemId int64
	lookIdString := r.URL.Query().Get("look_id")
	if lookIdString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		lookId, _ = strconv.ParseInt(lookIdString, 10, 64)
	}

	itemIdString := r.URL.Query().Get("item_id")
	if itemIdString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		itemId, _ = strconv.ParseInt(itemIdString, 10, 64)
	}

	err := s.dataService.AddToLook(r.Context(), lookId, itemId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *service) GetLookData(w http.ResponseWriter, r *http.Request) {
	var lookId int64
	lookIdString := r.URL.Query().Get("look_id")
	if lookIdString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		lookId, _ = strconv.ParseInt(lookIdString, 10, 64)
	}

	f, err := s.dataService.GetLookData(r.Context(), lookId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
