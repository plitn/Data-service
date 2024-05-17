package data_service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Data-service/internal/model"
	"github.com/Data-service/internal/repository"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"time"
)

type service struct {
	repository  repository.Repository
	minioClient *minio.Client
}

func NewDataService(repo repository.Repository, endpoint, accessId, accessKey string) *service {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessId, accessKey, ""),
		Secure: false,
	})
	if err != nil {
		fmt.Printf("cannot create minio client, %v", err)
		return nil
	}
	return &service{
		minioClient: minioClient,
		repository:  repo,
	}
}

func (s *service) CreateFile(ctx context.Context, item model.Item, fileName string, bucketName string, fileSize int64, reader io.Reader) error {
	exists, err := s.minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Printf("cannot check if bucket exists: %s\n", err.Error())
		return err
	}
	if !exists {
		err := s.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Printf("cannot create bucket: %s\n", err.Error())
			return err
		}
	}
	//newReader := reader
	newReader, newSize, err := s.ClearFile(fileName, reader)
	if err != nil {
		return err
	}
	_, err = s.minioClient.PutObject(ctx, bucketName, fileName, newReader, newSize,
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		})
	if err != nil {
		fmt.Printf("cannot upload file: %s\n", err.Error())
		return err
	}
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	presignedURL, err := s.minioClient.PresignedGetObject(context.Background(), bucketName, fileName,
		time.Second*24*60*60*7, reqParams)
	if err != nil {
		fmt.Println(err)
	}
	item.Url = presignedURL.String()
	err = s.repository.SaveItem(ctx, item)
	fmt.Println(item)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// тут получаем ссылку на файл и идем сохранять в бд данные по инфе о итеме
	return nil
}

func (s *service) ClearFile(fileName string, reader io.Reader) (io.Reader, int64, error) {
	outputPath := fileName
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, 0, err
	}
	tempFile, err := ioutil.TempFile("", "temp_input_*.jpg")
	if err != nil {
		fmt.Println("Ошибка при создании временного файла:", err)
		return nil, 0, err
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(data)
	if err != nil {
		fmt.Println("Ошибка при записи во временный файл:", err)
		return nil, 0, err
	}
	tempFile.Close()

	cmd := exec.Command("python3", "build/script.py", tempFile.Name(), outputPath)
	//cmd.Path = "/Users/platondmitrin/go/src/github.com/Data-service/.venv/bin/python3"

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Ошибка при выполнении команды:", err)
		return nil, 0, err
	}

	outputData, err := ioutil.ReadFile(outputPath)
	if err != nil {
		fmt.Println("Ошибка при чтении выходного файла:", err)
		return nil, 0, err
	}
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		fmt.Println("Ошибка при получении информации о файле:", err)
		return nil, 0, err
	}
	fileSize := fileInfo.Size()

	os.Remove(outputPath)
	return bytes.NewReader(outputData), fileSize, nil
}

func (s *service) GetItem(ctx context.Context, itemId int64) (model.Item, error) {
	item, err := s.repository.GetItem(ctx, itemId)
	if err != nil {
		return model.Item{}, fmt.Errorf("cannot get item %d: %s", itemId, err.Error())
	}
	return item, nil
}

func (s *service) GetItemsByUser(ctx context.Context, userId int64) ([]model.Item, error) {
	item, err := s.repository.GetItemsByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("cannot get user items, userId %d: %s", userId, err.Error())
	}
	return item, nil
}

func (s *service) UpdateItem(ctx context.Context, item model.Item) error {
	err := s.repository.UpdateItem(ctx, item)
	if err != nil {
		return fmt.Errorf("cannot update item %d: %s", item.Id, err.Error())
	}
	return nil
}

func (s *service) DeleteItem(ctx context.Context, itemId int64) error {
	itemInfo, err := s.repository.GetItem(ctx, itemId)
	if err != nil {
		return fmt.Errorf("cannot get item %d: %s", itemId, err.Error())
	}
	err = s.minioClient.RemoveObject(ctx, fmt.Sprintf("usr%d", itemInfo.UserId), itemInfo.FileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("cannot delete item from storage %d: %s", itemId, err.Error())
	}
	err = s.repository.DeleteItem(ctx, itemId)
	if err != nil {
		return fmt.Errorf("cannot delete item from storage %d: %s", itemId, err.Error())
	}
	return nil
}

func (s *service) CreateCapsule(ctx context.Context, capsule model.Capsule) error {
	err := s.repository.CreateCapsule(ctx, capsule)
	if err != nil {
		return fmt.Errorf("cannot create capsule: %s", err.Error())
	}
	return nil
}

func (s *service) UpdateCapsule(ctx context.Context, capsule model.Capsule) error {
	err := s.repository.UpdateCapsuleName(ctx, capsule)
	if err != nil {
		return fmt.Errorf("cannot update capsule: %s", err.Error())
	}
	return nil
}

func (s *service) DeleteCapsule(ctx context.Context, capsuleId int64) error {
	err := s.repository.DeleteCapsule(ctx, capsuleId)
	if err != nil {
		return fmt.Errorf("cannot delete capsule: %s", err.Error())
	}
	return nil
}

func (s *service) GetCapsule(ctx context.Context, capsuleId int64) (model.Capsule, error) {
	capsule, err := s.repository.GetCapsule(ctx, capsuleId)
	if err != nil {
		return model.Capsule{}, fmt.Errorf("cannot get capsule: %s", err.Error())
	}
	return capsule, nil
}

func (s *service) GetCapsulesByUser(ctx context.Context, userId int64) ([]model.Capsule, error) {
	capsules, err := s.repository.GetCapsulesByUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("cannot get capsules for user %d: %s", userId, err.Error())
	}
	return capsules, nil
}

func (s *service) AddItemToCapsule(ctx context.Context, itemCapsule model.ItemCapsule) error {
	err := s.repository.AddItemToCapsule(ctx, itemCapsule)
	if err != nil {
		return fmt.Errorf("cannot add item to capsule: %s", err.Error())
	}
	return nil
}

func (s *service) DeleteItemFromCapsule(ctx context.Context, itemCapsule model.ItemCapsule) error {
	err := s.repository.DeleteItemFromCapsule(ctx, itemCapsule)
	if err != nil {
		return fmt.Errorf("cannot delete item from capsule: %s", err.Error())
	}
	return nil
}

func (s *service) GetItemsInCapsule(ctx context.Context, capsuleId int64) ([]model.Item, error) {
	itemsIds, err := s.repository.GetCapsuleItems(ctx, capsuleId)
	if err != nil {
		return nil, fmt.Errorf("cannot get items ids in capsule: %s", err.Error())
	}
	items, err := s.repository.GetItemsByIds(ctx, itemsIds)
	if err != nil {
		return nil, fmt.Errorf("cannot get items in capsule: %s", err.Error())
	}
	return items, nil
}

func (s *service) GetFilesByUser(ctx context.Context, userId string) ([]*model.File, error) {
	objects, err := s.getBucketFiles(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get objects. err: %w", err)
	}
	if len(objects) == 0 {
		return nil, fmt.Errorf("no objects found for user %d", userId)
	}

	var files []*model.File
	for _, obj := range objects {
		stat, err := obj.Stat()
		if err != nil {
			fmt.Printf("failed to get objects. err: %v", err)
			continue
		}
		buffer := make([]byte, stat.Size)
		_, err = obj.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("failed to get objects. err: %v", err)
			continue
		}
		f := model.File{
			ID:    stat.Key,
			Name:  stat.UserMetadata["Name"],
			Size:  stat.Size,
			Bytes: buffer,
		}
		files = append(files, &f)
		obj.Close()
	}

	return files, nil
}

func (s *service) getBucketFiles(ctx context.Context, bucketName string) ([]*minio.Object, error) {
	var files []*minio.Object
	for bucketItem := range s.minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{WithMetadata: true}) {
		if bucketItem.Err != nil {
			fmt.Printf("failed to list object from minio bucket %s. err: %v", bucketName, bucketItem.Err)
			continue
		}
		object, err := s.minioClient.GetObject(ctx, bucketName, bucketItem.Key, minio.GetObjectOptions{})
		if err != nil {
			fmt.Printf("failed to get object key=%s from minio bucket %s. err: %v", bucketItem.Key, bucketName, bucketItem.Err)
			continue
		}
		files = append(files, object)
	}
	return files, nil
}

func (s *service) GetFilesArray(ctx context.Context, bucketName string, fileNames map[string]struct{}) ([]*model.File, error) {
	objects, err := s.getBucketFiles(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get objects. err: %w", err)
	}
	if len(objects) == 0 {
		return nil, fmt.Errorf("no objects found for user %d", bucketName)
	}

	var files []*model.File
	for _, obj := range objects {
		stat, err := obj.Stat()
		if err != nil {
			fmt.Printf("failed to get objects. err: %v", err)
			continue
		}
		if _, ok := fileNames[stat.Key]; !ok {
			continue
		}
		buffer := make([]byte, stat.Size)
		_, err = obj.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("failed to get objects. err: %v", err)
			continue
		}
		f := model.File{
			ID:    stat.Key,
			Name:  stat.Key,
			Size:  stat.Size,
			Bytes: buffer,
		}
		files = append(files, &f)
		obj.Close()
	}

	return files, nil
}

func (s *service) DeleteFile(ctx context.Context, bucketName, fileName string) error {
	err := s.minioClient.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file. err: %w", err)
	}
	return nil
}

//func (s *service) GetItemsArray(ctx context.Context, ids []int64) ([]model.Item, error) {
//	items, err := s.repository.GetItemsArray(ctx, ids)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get items array. err: %w", err)
//	}
//	return items, nil
//}

func (s *service) GetLookData(ctx context.Context, lookId int64) (model.LookData, error) {
	look, err := s.repository.GetLook(ctx, lookId)
	if err != nil {
		fmt.Println(err)

		return model.LookData{}, fmt.Errorf("failed to get look. err: %w", err)
	}
	fmt.Println(look)
	lk, err := s.repository.GetLookItems(ctx, lookId)
	if err != nil {
		fmt.Println(err)

		return model.LookData{}, fmt.Errorf("failed to get look items. err: %w", err)
	}
	fmt.Println(lk)
	itemIds := make([]int64, len(lk))
	for _, item := range lk {
		itemIds = append(itemIds, item.ItemId)
	}
	items, err := s.repository.GetItemsArray(ctx, itemIds)
	if err != nil {
		fmt.Println(err)
		return model.LookData{}, fmt.Errorf("failed to get items array. err: %w", err)
	}
	lookData := model.LookData{
		Id:        look.Id,
		UserId:    look.UserId,
		LookName:  look.LookName,
		StylistId: look.StylistId,
		Items:     items,
	}
	return lookData, nil
}

func (s *service) AddToLook(ctx context.Context, lookId, itemId int64) error {
	err := s.repository.AddItemToLook(ctx, lookId, itemId)
	if err != nil {
		return fmt.Errorf("cannot add item to look. err: %w", err)
	}
	return nil
}

func (s *service) DeleteLook(ctx context.Context, lookId int64) error {
	err := s.repository.DeleteLookItems(ctx, lookId)
	if err != nil {
		return fmt.Errorf("cannot delete item from look. err: %w", err)
	}
	err = s.repository.DeleteLook(ctx, lookId)
	if err != nil {
		return fmt.Errorf("cannot delete look. err: %w", err)
	}
	return nil
}

func (s *service) CreateLook(ctx context.Context, look model.Look) error {
	err := s.repository.CreateLook(ctx, look)
	if err != nil {
		return fmt.Errorf("cannot create look. err: %w", err)
	}
	return nil
}
