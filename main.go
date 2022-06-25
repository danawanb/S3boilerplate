package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("image_1")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	//dir: temp-image
	tempFile, err := ioutil.TempFile("temp-image", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully Uploaded file\n")
}

func UploadS3Bucket(w http.ResponseWriter, r *http.Request) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-3"),
	})
	if err != nil {
		//panic(err)
		log.Fatal("tidak dapat mendapatkan session")
	}

	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("image_1")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	nameid := uuid.New()
	filename := nameid.String() + ".png"

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Bucket:      aws.String("awasbucket"),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Upload Result: %+v\n", result)
	fmt.Fprintf(w, "Successfully Uploaded file\n")
}

func SetupRoutes() {
	//fs := http.FileServer(http.Dir("/temp-image"))

	http.HandleFunc("/uploads", UploadS3Bucket)
	http.HandleFunc("/upload", UploadFile)
	//http.Handle("/files", fs)
	http.ListenAndServe("localhost:8080", nil)
}

func UploadItem(sess *session.Session) {
	f, err := os.Open("programwell.png")
	if err != nil {
		log.Fatal("UHUHUH")
	}

	defer f.Close()
	fmt.Println(f)
	nameid := uuid.New()
	filename := nameid.String() + ".png"

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Bucket:      aws.String("awasbucket"),
		Key:         aws.String(filename),
		Body:        f,
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Upload Result: %+v\n", result)
}
func ListBucket(sess *session.Session) {
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		panic(err)
	}
	for _, bucket := range result.Buckets {
		log.Printf("Bucketnya %s\n", aws.StringValue(bucket.Name))
		log.Printf("Bucketnya %s\n", aws.TimeValue(bucket.CreationDate))
	}
}

func ListItems(sess *session.Session) {
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(
		&s3.ListObjectsV2Input{
			Bucket: aws.String("awasbucket"),
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, item := range resp.Contents {
		log.Printf("Name: %s\n", *item.Key)
		log.Printf("Name: %d\n", *item.Size)
	}
}

func Downloader(sess *session.Session) {

	file, err := os.Create("asd.png")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String("awasbucket"),
			Key:    aws.String("122331811_2793238694335711_5476438353893715843_n.jpg"),
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Successfully download file")
}

func DeleteItem(sess *session.Session) {
	svc := s3.New(sess)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String("awasbucket"),
		Key:    aws.String("programwell_.png"),
	}

	result, err := svc.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Fatal(aerr.Error())
			}
		} else {
			log.Fatal(err.Error())
		}
	}
	log.Printf("result %+v\n", result)
}
func main() {
	fmt.Println("Downloads and Uploads")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-3"),
	})
	if err != nil {
		//panic(err)
		log.Fatal("tidak dapat mendapatkan session")
	}
	theMine := []string{"one1", "two2", "three3"}
	mineChan := make(chan string)

	go func(mine []string) {
		for _, item := range mine {
			mineChan <- item
		}
	}(theMine)

	go func() {
		for i := 0; i < 3; i++ {
			foundMe := <-mineChan
			fmt.Println("received " + foundMe + "from you")
		}
	}()
	<-time.After(time.Second * 5)
	//UploadItem(sess)
	//ListBucket(sess)

	//Downloader(sess)
	//DeleteItem(sess)
	ListItems(sess)
	SetupRoutes()
}
