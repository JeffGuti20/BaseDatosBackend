package repository

import (
	"context"
	"database/sql"

	"rest/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// CrearRol ejecuta el procedimiento almacenado paCrearRol.
func (r *CategoryRepository) CrearCategoria(ctx context.Context, name string) (*models.Category, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paCrearCategoria @Nombre = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var category models.Category

	err = stmt.QueryRowContext(ctx, name).Scan(
		&category.CategoryID,
	)

	if err != nil {
		return nil, err
	}

	category.Name = name

	return &category, nil
}

// ObtenerRoles ejecuta el procedimiento almacenado paObtenerRoles.
func (r *CategoryRepository) ObtenerCategorias(ctx context.Context) ([]models.Category, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paObtenerCategorias",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {

		var category models.Category

		err := rows.Scan(
			&category.CategoryID,
			&category.Name,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// ObtenerRol ejecuta el procedimiento almacenado paObtenerRol.
func (r *CategoryRepository) ObtenerCategoria(ctx context.Context, categoryID int32) (*models.Category, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paObtenerCategoria @IdCategoria = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var category models.Category

	err = stmt.QueryRowContext(ctx, categoryID).Scan(
		&category.CategoryID,
		&category.Name,
	)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// ActualizarRol ejecuta el procedimiento almacenado paActualizarRol.
func (r *CategoryRepository) ActualizarCategoria(
	ctx context.Context,
	categoryID int32,
	name string,
) (*models.Category, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paActualizarCategoria @IdCategoria = @p1, @Nombre = @p2",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var category models.Category

	err = stmt.QueryRowContext(
		ctx,
		categoryID,
		name,
	).Scan(
		&category.CategoryID,
		&category.Name,
	)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// EliminarRol ejecuta el procedimiento almacenado paEliminarRol.
func (r *CategoryRepository) EliminarCategoria(ctx context.Context, category int32) error {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paEliminarCategoria @IdCategoria = @p1",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, category)

	return err
}
