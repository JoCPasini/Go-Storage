package product

import (
	"context"
	"log"
	"testing"

	"github.com/JosePasiniMercadolibre/Go-storage/domain"
	"github.com/JosePasiniMercadolibre/Go-storage/util"
	"github.com/stretchr/testify/assert"
)

func Test_StoreUserDynamo_OK(t *testing.T) {
	db, err := util.InitDynamoDB()
	if err != nil {
		log.Println(err)
	}

	user := domain.User{
		Id:         "2",
		Firstname:  "firstname",
		Lastname:   "lastname",
		Username:   "username",
		Password:   "12345",
		Email:      "email.user@mercadolibre.com",
		IP:         "127.0.0.1",
		MacAddress: "FF:FF:FF:FF:FF:FF",
		Website:    "website.com",
		Image:      "image.png",
	}

	myRepo := NewDynamoRepository(db, "Users")

	err = myRepo.Store(context.Background(), &user)
	assert.NoError(t, err)
}

func Test_GetOneUserDynamo_OK(t *testing.T) {
	db, err := util.InitDynamoDB()
	if err != nil {
		log.Fatal(err)
	}

	myRepo := NewDynamoRepository(db, "Users")

	userExpected, error := myRepo.GetOne(context.Background(), "1")

	// En caso de encontrar el User con el ID indicado arriba, podremos compararlo con NotNil
	assert.NotNil(t, userExpected)
	assert.NoError(t, error)
}
