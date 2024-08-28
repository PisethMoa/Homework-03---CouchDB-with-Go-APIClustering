package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-kivik/couchdb/v4"
	kivik "github.com/go-kivik/kivik/v4"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "main.go/docs"
)

// @title Student API
// @version 1.0
// @description This is a simple API to interact with CouchDB and perform CRUD operations.
// @host localhost:8080
// @BasePath /
var client *kivik.Client

type Response struct {
	Message string `json:"message"`
	Rev     string `json:"rev,omitempty"`
	Error   string `json:"error,omitempty"`
}
type DeleteResponse struct {
	Message string `json:"message"`
}

func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var err error
	client, err = kivik.New("couch", "http://admin:123@localhost:5984/")
	if err != nil {
		log.Fatalf("Failed to connect to CouchDB: %v", err)
	}

	exists, err := client.DBExists(context.TODO(), "student")
	if err != nil {
		log.Fatalf("Failed to check if database exists: %v", err)
	}

	if !exists {
		err = client.CreateDB(context.TODO(), "student")
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		fmt.Println("Database created successfully.")
	} else {
		fmt.Println("Database already exists.")
	}

	r.POST("/insert", insertDocument)
	r.POST("/upload", uploadFileHandler)
	r.GET("/file/:docID/:filename", getFileHandler)
	r.GET("/documents", getAllDocumentsHandler)
	r.GET("/document/:id", getDocumentByIDHandler)
	r.GET("/changes", filterDocuments)
	r.PUT("/document/:docID", updateDocumentHandler)
	r.DELETE("/document/:docID", deleteDocumentHandler)

	r.Run(":8080")
}

// Insert Document
// @Summary Insert a document
// @Description Inserts a new document into the CouchDB
// @Tags document
// @Accept json
// @Produce json
// @Param document body map[string]interface{} true "Document"
// @Success 200 {string} string "Document inserted successfully."
// @Failure 400 {string} string "Failed to decode JSON."
// @Failure 500 {string} string "Failed to insert document."
// @Router /insert [post]
func insertDocument(c *gin.Context) {
	var doc map[string]interface{}
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode JSON."})
		return
	}
	id, ok := doc["_id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document must contain '_id' field."})
		return
	}
	_, err := client.DB("student").Put(context.TODO(), id, doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert document."})
		return
	}
	c.JSON(http.StatusOK, "Document inserted successfully.")
}

// uploadFileHandler godoc
// @Summary Uploads a file
// @Description Upload a file to CouchDB as an attachment
// @Accept  mpfd
// @Produce json
// @Param  file formData file true "File to upload"
// @Param  docID formData string true "Document ID"
// @Success 200 {string} string "File uploaded successfully."
// @Router /upload [post]
func uploadFileHandler(c *gin.Context) {
	docID := c.PostForm("docID")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := client.DB("student")

	doc := make(map[string]interface{})
	err = db.Get(context.TODO(), docID).ScanDoc(&doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get document."})
		return
	}
	rev, _ := doc["_rev"].(string)

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer openedFile.Close()

	_, err = db.PutAttachment(context.TODO(), docID, &kivik.Attachment{
		Filename:    file.Filename,
		Content:     openedFile,
		ContentType: file.Header.Get("Content-Type"),
	}, kivik.Options{"rev": rev})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "File uploaded successfully"})
}

// getFileHandler godoc
// @Summary Get a file
// @Description Retrieves an attachment from a CouchDB document
// @Tags file
// @Produce octet-stream
// @Param docID path string true "Document ID"
// @Param filename path string true "Filename"
// @Success 200 {file} file "File downloaded successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "File not found"
// @Failure 500 {string} string "Failed to retrieve file"
// @Router /file/{docID}/{filename} [get]
func getFileHandler(c *gin.Context) {
	docID := c.Param("docID")
	filename := c.Param("filename")

	db := client.DB("student")

	attachment, err := db.GetAttachment(context.TODO(), docID, filename)
	if err != nil {
		if kivik.HTTPStatus(err) == http.StatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file: " + err.Error()})
		}
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", attachment.ContentType)
	c.Stream(func(w io.Writer) bool {
		_, err := io.Copy(w, attachment.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream file content: " + err.Error()})
			return false
		}
		return false
	})
}

// Get All Documents Handler
// @Summary Get all documents
// @Description Retrieves all documents from the CouchDB student database
// @Tags document
// @Produce json
// @Success 200 {array} map[string]interface{} "Documents retrieved successfully."
// @Failure 500 {string} string "Failed to retrieve documents."
// @Router /documents [get]
func getAllDocumentsHandler(c *gin.Context) {
	url := "http://admin:123@localhost:5984/student/_all_docs?include_docs=true"

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents."})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body."})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON response."})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Get Document by ID Handler
// @Summary Get a document by ID
// @Description Retrieves a specific document from the CouchDB student database by its ID
// @Tags document
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} map[string]interface{} "Document retrieved successfully."
// @Failure 404 {string} string "Document not found."
// @Failure 500 {string} string "Failed to retrieve document."
// @Router /document/{id} [get]
func getDocumentByIDHandler(c *gin.Context) {
	id := c.Param("id")

	db := client.DB("student")

	doc := make(map[string]interface{})
	err := db.Get(context.TODO(), id).ScanDoc(&doc)
	if err != nil {
		if kivik.HTTPStatus(err) == http.StatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document"})
		}
		return
	}

	c.JSON(http.StatusOK, doc)
}

// Get Changes with Filter
// @Summary Get changes from CouchDB with a filter
// @Description Retrieves changes from CouchDB using a specified filter
// @Tags changes
// @Accept json
// @Produce json
// @Param address query string false "Address"
// @Param age query int false "Age"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} string "Failed to retrieve changes."
// @Router /changes [get]
func filterDocuments(c *gin.Context) {
	address := c.Query("address")
	age := c.Query("age")

	if address == "" || age == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address and age are required."})
		return
	}

	ageInt, err := strconv.Atoi(age)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age format."})
		return
	}

	viewURL := fmt.Sprintf("http://admin:123@localhost:5984/student/_changes?filter=filters/by_address_and_age&address=%s&age=%d",
		url.QueryEscape(address), ageInt)

	resp, err := http.Get(viewURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve changes."})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response."})
		return
	}

	c.JSON(http.StatusOK, result)
}

// updateDocumentHandler godoc
// @Summary Update an existing document
// @Description Updates an existing document in the CouchDB database
// @Tags document
// @Accept json
// @Produce json
// @Param docID path string true "Document ID"
// @Param document body map[string]interface{} true "Document Data"
// @Success 200 {object} Response "Document updated successfully"
// @Failure 400 {object} Response "Invalid request"
// @Failure 404 {object} Response "Document not found"
// @Failure 500 {object} Response "Failed to update document"
// @Router /document/{docID} [put]
func updateDocumentHandler(c *gin.Context) {
	docID := c.Param("docID")

	db := client.DB("student")
	existingDoc := make(map[string]interface{})
	err := db.Get(context.TODO(), docID).ScanDoc(&existingDoc)
	if err != nil {
		if kivik.HTTPStatus(err) == http.StatusNotFound {
			c.JSON(http.StatusNotFound, Response{Error: "Document not found"})
		} else {
			c.JSON(http.StatusInternalServerError, Response{Error: "Failed to retrieve document: " + err.Error()})
		}
		return
	}

	updatedData := make(map[string]interface{})
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "Invalid request body"})
		return
	}

	for key, value := range updatedData {
		existingDoc[key] = value
	}

	rev, err := db.Put(context.TODO(), docID, existingDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: "Failed to update document: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Message: "Document updated successfully", Rev: rev})
}

// deleteDocumentHandler godoc
// @Summary Delete a document
// @Description Deletes a document from the CouchDB database
// @Tags document
// @Param docID path string true "Document ID"
// @Success 200 {object} DeleteResponse "Document deleted successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Document not found"
// @Failure 500 {string} string "Failed to delete document"
// @Router /document/{docID} [delete]
func deleteDocumentHandler(c *gin.Context) {
	docID := c.Param("docID")

	db := client.DB("student")
	row := db.Get(context.TODO(), docID)
	var doc map[string]interface{}
	if err := row.ScanDoc(&doc); err != nil {
		if kivik.HTTPStatus(err) == http.StatusNotFound {
			c.JSON(http.StatusNotFound, DeleteResponse{Message: "Document not found"})
		} else {
			c.JSON(http.StatusInternalServerError, DeleteResponse{Message: "Failed to retrieve document: " + err.Error()})
		}
		return
	}

	rev := doc["_rev"].(string)
	_, err := db.Delete(context.TODO(), docID, rev)
	if err != nil {
		c.JSON(http.StatusInternalServerError, DeleteResponse{Message: "Failed to delete document: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, DeleteResponse{Message: "Document deleted successfully"})
}
