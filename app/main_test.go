package main

import (
	"bytes"
	"encoding/json"
	"go-test/dto"
	"go-test/infra"
	"go-test/models"

	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Test user
var TestUser = dto.SignUpInput{
	Email:    "admin@example.com",
	Password: "admin_password",
}

func TestMain(m *testing.M) {
	// open test env
	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func setUp() *gin.Engine {
	db := infra.SetUpDB()
	err := db.AutoMigrate(&models.Item{}, &models.User{})
	if err != nil {
		log.Fatal(err)
	}

	router := setUpRouter(db)
	return router
}

// _signUp method for Test User
//
// request: POST, "/auth/signup", dto.SignUpInput
// response: {"user": dto.SignUpInput}
func _signUp(t *testing.T, router *gin.Engine, signUpTestUser dto.SignUpInput) dto.SignUpInput {
	// sign up request body
	reqBody, err := json.Marshal(signUpTestUser)
	assert.Equal(t, err, nil)

	// http request
	req, err := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(reqBody))
	assert.Equal(t, err, nil)

	// http response
	w := httptest.NewRecorder()

	// do API http request
	router.ServeHTTP(w, req)

	// ? signup succeeded
	assert.Equal(t, http.StatusCreated, w.Code)

	// put http response result in resBody
	var resBody map[string]dto.SignUpInput
	err = json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, err, nil)

	return resBody["user"]
}

// _login method for Test User to login
//
// request: POST, "/auth/login", dto.LoginInput
// response: {"token": "xxx"}
func _login(t *testing.T, router *gin.Engine, loginUser dto.LoginInput) string {
	// request body
	reqBody, err := json.Marshal(loginUser)
	assert.Equal(t, err, nil)

	// http request
	req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
	assert.Equal(t, err, nil)

	// http response
	w := httptest.NewRecorder()

	// do API http Request
	router.ServeHTTP(w, req)

	// ? login succeeded
	assert.Equal(t, http.StatusOK, w.Code)

	// put http response result in resBody
	var resBody map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, err, nil)

	return resBody["token"]
}

// _createTestItems make count numbers item
//
// request: POST, "/items", dto.CreateItemInput
// response: {"data": models.Item}
func _createTestItems(t *testing.T, router *gin.Engine, token string, count int) []uint {
	var createdItemIDs []uint

	for i := 1; i <= count; i++ {
		newItem := dto.CreateItemInput{
			Name:        "test" + strconv.Itoa(i),
			Price:       uint(i * 100),
			Description: "No." + strconv.Itoa(i),
		}

		// http request body
		reqBody, err := json.Marshal(newItem)
		assert.Equal(t, err, nil)

		// http request
		req, err := http.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(reqBody))
		assert.Equal(t, err, nil)

		// set token in header
		req.Header.Set("Authorization", "Bearer "+token)

		// http response
		w := httptest.NewRecorder()

		// do API request
		router.ServeHTTP(w, req)

		// ? new item is created successfully
		assert.Equal(t, http.StatusCreated, w.Code)

		// put http response result in resBody
		var resBody map[string]models.Item
		err = json.Unmarshal(w.Body.Bytes(), &resBody)
		assert.Equal(t, err, nil)

		createdItemIDs = append(createdItemIDs, resBody["data"].ID)
	}

	return createdItemIDs
}

// _findMyAll
//
// repository package内に、FindAll APIを参考に、
// 下記 FindMyAll APIを独自に作成。（services、controllersも同様）
// [内容]
// サインインしたユーザー自身が作成したitemのみ取得
//
//	func (ir *ItemDBRepository) FindMyAll(userId uint) (*[]models.Item, error) {
//	   items := make([]models.Item, 0)
//
//	   result := ir.db.Find(&items, "user_id = ?", userId)
//	   if result.Error != nil {
//	      if result.Error.Error() == "record not found" {
//	         return nil, errors.New("item not found")
//	      }
//	      return nil, result.Error
//	   }
//
//	   return &items, nil
//	}
//
// request: GET, "/items/mine", dto.CreateItemInput
// response: {"data": []models.Item}
func _findMyAll(t *testing.T, router *gin.Engine, token string) []models.Item {
	// http request
	req, err := http.NewRequest(http.MethodGet, "/items/mine", nil)
	assert.Equal(t, err, nil)

	// http response
	w := httptest.NewRecorder()

	// set Authorization token in header
	req.Header.Set("Authorization", "Bearer "+token)

	// do API FindMyAll request
	router.ServeHTTP(w, req)

	// put http response result in resBody
	var resBody map[string][]models.Item
	err = json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, err, nil)

	return resBody["data"]
}

// _update
//
// request: PUT, "/items/:id", dto.UpdateItemInput
// response: {"data": models.Item}
func _update(t *testing.T, router *gin.Engine, token string, updateId string, updateItem dto.UpdateItemInput) models.Item {
	// update request body
	reqBody, err := json.Marshal(updateItem)
	assert.Equal(t, err, nil)

	// update http request
	req, err := http.NewRequest(http.MethodPut, "/items/"+updateId, bytes.NewBuffer(reqBody))
	assert.Equal(t, err, nil)

	// set token in header
	req.Header.Set("Authorization", "Bearer "+token)

	// http response
	w := httptest.NewRecorder()

	// do API request
	router.ServeHTTP(w, req)

	// ? update succeeded
	assert.Equal(t, http.StatusOK, w.Code)

	// put http response result in resBody
	var resBody map[string]models.Item
	err = json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, err, nil)

	return resBody["data"]
}

// _delete
//
// request: DELETE, "/items/:id", nil
// response:
func _delete(t *testing.T, router *gin.Engine, token string, deleteId string) {
	// http request
	req, err := http.NewRequest(http.MethodDelete, "/items/"+deleteId, nil)
	assert.Equal(t, err, nil)

	// http response
	w := httptest.NewRecorder()

	// set Authorization token in header
	req.Header.Set("Authorization", "Bearer "+token)

	// do API request
	router.ServeHTTP(w, req)

	// assert status check
	assert.Equal(t, w.Code, http.StatusOK)
}

// _findAll
//
// request: GET, "/items", nil
// response: {"data": []models.Item}
func _findAll(t *testing.T, router *gin.Engine) []models.Item {
	// http request
	req, err := http.NewRequest(http.MethodGet, "/items", nil)
	assert.Equal(t, err, nil)

	// http response
	w := httptest.NewRecorder()

	// do API request
	router.ServeHTTP(w, req)

	// assert status check
	assert.Equal(t, w.Code, http.StatusOK)

	// put http response result in resBody
	var resBody map[string][]models.Item
	err = json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, err, nil)

	return resBody["data"]
}

func TestFindAll(t *testing.T) {
	// set up DB
	db := infra.SetUpDB()
	err := db.AutoMigrate(&models.Item{}, &models.User{})
	if err != nil {
		log.Fatal(err)
	}

	/* 1. add item and user test data */
	items := []models.Item{
		{Name: "test1", Price: 100, Description: "test1", SoldOut: false},
		{Name: "test2", Price: 200, Description: "test2", SoldOut: true},
		{Name: "test3", Price: 300, Description: "test3", SoldOut: false},
	}
	users := []models.User{
		{Email: "test1@example.com", Password: "test1 password"},
		{Email: "test2@example.com", Password: "test2 password"},
		{Email: "test3@example.com", Password: "test3 password"},
	}
	for _, item := range items {
		db.Create(&item)
	}
	for _, user := range users {
		db.Create(&user)
	}

	/* 2. set up router */
	router := setUpRouter(db)

	/* 3. Test FindAll */
	allItems := _findAll(t, router)

	/* 4. assert check */
	assert.Equal(t, 3, len(allItems))
}

func TestFindMyAll(t *testing.T) {
	testRouter := setUp()

	/* 1. Sign up by TestUser */
	signUpUser := _signUp(t, testRouter, TestUser)

	/* 2. Login */
	token := _login(t, testRouter, dto.LoginInput{Email: signUpUser.Email, Password: signUpUser.Password})

	/* 3. Create 10 items */
	newItemIDs := _createTestItems(t, testRouter, token, 10)

	/* 4. call _findMyAll method */
	myItems := _findMyAll(t, testRouter, token)

	/* 5. assert check */
	assert.Equal(t, len(newItemIDs), len(myItems))
}

func TestCreate(t *testing.T) {
	testRouter := setUp()

	/* 1. Sign up by TestUser */
	signUpUser := _signUp(t, testRouter, TestUser)

	/* 2. Login */
	token := _login(t, testRouter, dto.LoginInput{Email: signUpUser.Email, Password: signUpUser.Password})

	/* 3. asset check (not yet created) */
	// get all items
	currentItems := _findMyAll(t, testRouter, token)
	// ? is empty at first
	assert.Equal(t, 0, len(currentItems))

	/* 4. create 10 items */
	newItemIDs := _createTestItems(t, testRouter, token, 10)

	/* 5. assert check (created 10 items) */
	assert.Equal(t, len(newItemIDs), 10)
}

func TestDelete(t *testing.T) {
	testRouter := setUp()

	/* 1. Sign up by TestUser */
	signUpUser := _signUp(t, testRouter, TestUser)

	/* 2. Login */
	token := _login(t, testRouter, dto.LoginInput{Email: signUpUser.Email, Password: signUpUser.Password})

	/* 3. Create 10 items by login user */
	newItemIDs := _createTestItems(t, testRouter, token, 10)

	/* 4. Delete one item */
	// delete first index
	deleteId := strconv.Itoa(int(newItemIDs[0]))
	_delete(t, testRouter, token, deleteId)

	/* 5. get items */
	myItems := _findMyAll(t, testRouter, token)

	/* 6. assert check */
	assert.Equal(t, 9, len(myItems))
}

func TestUpdate(t *testing.T) {
	testRouter := setUp()

	/* 1. Signup by test user */
	newUser := _signUp(t, testRouter, TestUser)

	/* 2. Login */
	token := _login(t, testRouter, dto.LoginInput{Email: newUser.Email, Password: newUser.Password})

	/* 3. Create 10 items */
	itemIDs := _createTestItems(t, testRouter, token, 10)

	/* 4. Update one of the created items */
	updateName := "updated name"
	updatePrice := uint(9999)
	updateDescription := "updated"

	updateItem := dto.UpdateItemInput{
		Name:        &updateName,
		Price:       &updatePrice,
		Description: &updateDescription,
	}
	// update first index item
	updateId := strconv.Itoa(int(itemIDs[0]))
	updatedItem := _update(t, testRouter, token, updateId, updateItem)

	/* 5. assert check */
	assert.Equal(t, updatedItem.Name, updateName)
	assert.Equal(t, updatedItem.Price, updatePrice)
	assert.Equal(t, updatedItem.Description, updateDescription)
}
