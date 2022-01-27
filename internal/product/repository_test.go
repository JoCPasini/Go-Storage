package product

import (
	"database/sql"
	"log"
	"testing"

	"github.com/JosePasiniMercadolibre/Go-storage/domain"
	"github.com/JosePasiniMercadolibre/Go-storage/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var (
	StorageDB *sql.DB
)

func Init() {
	dataSource := "root:@tcp(localhost:3306)/storage"
	var err error
	StorageDB, err = sql.Open("mysql", dataSource)
	if err != nil {
		log.Println(err)
	}

	if err = StorageDB.Ping(); err != nil {
		log.Println(err)
	}
	log.Println("Database Configured")
}

func Close() {
	StorageDB.Close()
}

func Test_GetByNameProduct_NIL(t *testing.T) {
	Init()
	defer Close()
	myRepo := NewRepository(StorageDB)
	// Le Pasamos un nombre que no exista en la base de datos.
	products, _ := myRepo.GetByName("Coca Cola 2")
	assert.Nil(t, products)
}

func Test_GetByNameProduct_OK(t *testing.T) {
	Init()
	defer Close()
	myRepo := NewRepository(StorageDB)

	// Busco por el nombre que se le seteó en el método Test_StoreProduct_OK
	products, _ := myRepo.GetByName("Coca Cola")
	assert.NotNil(t, products)
}

func Test_GetAllProduct_OK(t *testing.T) {
	Init()
	defer Close()
	myRepo := NewRepository(StorageDB)
	products, _ := myRepo.GetAll()

	assert.NotNil(t, products)
}

func Test_StoreProduct_OK(t *testing.T) {
	Init()
	defer Close()
	myRepo := NewRepository(StorageDB)

	product := domain.Product{
		Name:     "Coca Cola",
		Category: "Bebidas",
		Count:    10,
		Price:    10.0,
	}

	productSave, _ := myRepo.Store(product)
	assert.Equal(t, "Coca Cola", productSave.Name)
}

// TEST CON LIBRERÍA TXDB
func Test_StoreProduct_OK_TXDB(t *testing.T) {
	db, err := util.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	myRepo := NewRepository(db)

	// Creamos un objeto producto para guardarlo
	product := domain.Product{
		Name:     "Queso Cheddar",
		Category: "Lacteos",
		Count:    1,
		Price:    150,
	}
	// Guardamos el objeto producto
	productSave, err := myRepo.Store(product)
	// Verificamos que no haya error
	assert.NoError(t, err)
	// Verificamos que se haya guardado correctamente comparando el nombre.
	assert.Equal(t, "Queso Cheddar", productSave.Name)

	// Llamamos a todos los productos que se llamen "Queso Cheddar"
	productByName, err := myRepo.GetByName("Queso Cheddar")

	// Verificamos que no haya error
	assert.NoError(t, err)
	// Verificamos que nos devuelva una lista NO VACÍA
	assert.NotNil(t, productByName)
}

func Test_UpdateProduct_OK_TXDB(t *testing.T) {
	db, err := util.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	myRepo := NewRepository(db)

	// Obtenemos con GetOne un objeto para actualizar
	product, _ := myRepo.GetOne(4)

	// Creamos un objeto producto para updatear el objeto que recién obtuvimos
	productParaActualizar := domain.Product{
		ID:       4,
		Name:     "Queso Cheddar",
		Category: "Lacteos",
		Count:    1,
		Price:    150,
	}
	// Llamamos al método Update.
	product, _ = myRepo.Update(productParaActualizar)
	// Verificamos que no haya error
	assert.NoError(t, err)
	// Verificamos que se haya modificado correctamente
	assert.Equal(t, "Queso Cheddar", product.Name)
}

func Test_GetById_Product_OK(t *testing.T) {
	db, err := util.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	myRepo := NewRepository(db)

	product, err := myRepo.GetOne(4)
	assert.NoError(t, err)
	assert.Equal(t, "Coca Cola", product.Name)
}

func Test_DeleteById_Product_OK(t *testing.T) {
	db, err := util.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	myRepo := NewRepository(db)

	// Obtenemos un producto de la base de datos
	product, _ := myRepo.GetOne(4)
	// Verificamos que lo devuelva no nil
	assert.NotNil(t, product)

	// Eliminamos el producto de la base de datos
	err = myRepo.Delete(4)
	assert.NoError(t, err)

	// Intentamos obtener le mismo producto de la base de datos
	product, _ = myRepo.GetOne(4)
	// Verificamos que esta vez lo devuelva nil
	assert.Equal(t, 0, product.ID)
	assert.Equal(t, "", product.Category)
}

func Test_DeleteById_Product_ERROR(t *testing.T) {
	db, err := util.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	myRepo := NewRepository(db)

	// Obtenemos un producto de la base de datos
	product, _ := myRepo.GetOne(0)
	// Verificamos NO tener un producto con el ID 0, para luego intentar eliminar
	assert.Equal(t, "", product.Category)

	// Intentamos eliminar el producto de la base de datos
	err = myRepo.Delete(0)
	// Verificamos que pase por el error ErrNotFound al no encontrar producto con el id 0
	assert.Error(t, err)

}
