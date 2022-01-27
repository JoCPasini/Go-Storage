package product

import (
	"database/sql"
	"errors"
	"log"

	"github.com/JosePasiniMercadolibre/Go-storage/domain"
)

var (
	QueryGetAllProduct     = "SELECT id,name,category,count, price FROM products;"
	QueryGetOneProduct     = "SELECT id,name,category,count, price FROM products WHERE id = ?;"
	QueryInsertIntoProduct = "INSERT INTO products(name, category, count, price) VALUES (?,?,?,?)"
	QueryUpdateProduct     = "UPDATE products SET name=?, category=?, count=?, price=? WHERE id=?;"
	QueryDeleteProduct     = "DELETE FROM products WHERE id=?;"
	ErrNotFound            = errors.New("product not found")
)

type Repository interface {
	GetByName(name string) ([]domain.Product, error)
	Store(product domain.Product) (domain.Product, error)
	GetAll() ([]domain.Product, error)
	Update(product domain.Product) (domain.Product, error)
	GetOne(id int) (domain.Product, error)
	Delete(id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetOne(id int) (domain.Product, error) {
	var product domain.Product
	row := r.db.QueryRow(QueryGetOneProduct, id)
	err := row.Scan(&product.ID, &product.Name, &product.Category, &product.Count, &product.Price)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (r *repository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	rows, err := r.db.Query(QueryGetAllProduct)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.Count); err != nil {
			log.Println(err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *repository) GetByName(name string) ([]domain.Product, error) {
	var products []domain.Product
	rows, err := r.db.Query("SELECT * FROM products WHERE name = ?", name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Count, &product.Price); err != nil {
			log.Println(err)
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *repository) Store(product domain.Product) (domain.Product, error) {
	stmt, err := r.db.Prepare(QueryInsertIntoProduct)
	if err != nil {
		return domain.Product{}, err
	}

	var res sql.Result
	res, err = stmt.Exec(product.Name, product.Category, product.Count, product.Price)
	if err != nil {
		return domain.Product{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.Product{}, err
	}

	product.ID = int(id)

	return product, nil
}

func (r *repository) Update(product domain.Product) (domain.Product, error) {
	stmt, err := r.db.Prepare(QueryUpdateProduct)
	if err != nil {
		log.Fatal(err)
		return product, err
	}

	res, err := stmt.Exec(&product.Name, &product.Category, &product.Count, &product.Price, &product.ID)

	if err != nil {
		log.Fatal(err)
		return product, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Delete(id int) error {
	stmt, err := r.db.Prepare(QueryDeleteProduct)
	if err != nil {
		log.Fatal(err)
		return err
	}
	res, err := stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return err
	}
	if affected < 1 {
		return ErrNotFound
	}
	return nil
}
