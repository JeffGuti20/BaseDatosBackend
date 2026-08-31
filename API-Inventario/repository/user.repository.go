package repository

import (
	"context"
	"database/sql"

	"rest/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Login(ctx context.Context, email string) (*models.User, error) {

	// Preparar el procedimiento almacenado.
	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paIniciarSesion @Correo = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ejecutar el procedimiento.
	row := stmt.QueryRowContext(ctx, email)

	// Modelo que recibirá los datos del usuario.
	var user models.User

	err = row.Scan(
		&user.UserID,
		&user.RoleID,
		&user.Name,
		&user.Password,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}