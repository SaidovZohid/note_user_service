package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/SaidovZohid/note_user_service/pkg/utils"
	"github.com/SaidovZohid/note_user_service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) Create(u *repo.User) (*repo.User, error) {
	var updatedAt sql.NullTime
	query := `
		INSERT INTO users (
			first_name,
			last_name,
			phone_number,
			email,
			password,
			username,
			image_url,
			type
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at
	`
	err := ur.db.QueryRow(
		query,
		u.FirstName,
		u.LastName,
		utils.NullString(u.PhoneNumber),
		u.Email,
		u.Password,
		utils.NullString(u.Username),
		utils.NullString(u.ImageUrl),
		u.Type,
	).Scan(
		&u.ID,
		&u.CreatedAt,
		&updatedAt,
	)
	u.UpdatedAt = updatedAt.Time

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *userRepo) Get(user_id int64) (*repo.User, error) {
	var (
		u                               repo.User
		phoneNumber, username, imageUrl sql.NullString
		updatedAt                       sql.NullTime
	)
	query := `
		SELECT
			id,
			first_name,
			last_name,
			phone_number,
			email,
			username,
			image_url,
			type,
			created_at,
			updated_at
		FROM users WHERE id = $1 and deleted_at is null
	`
	err := ur.db.QueryRow(
		query,
		user_id,
	).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&phoneNumber,
		&u.Email,
		&username,
		&imageUrl,
		&u.Type,
		&u.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	u.PhoneNumber = phoneNumber.String
	u.ImageUrl = imageUrl.String
	u.Username = username.String
	u.UpdatedAt = updatedAt.Time

	return &u, nil
}

func (ur *userRepo) Update(u *repo.User) (*repo.User, error) {
	var phoneNumber, username, imageUrl sql.NullString
	query := `
		UPDATE users SET
			first_name = $1,
			last_name = $2,
			phone_number = $3,
			email = $4,
			username = $5,
			image_url = $6,
			updated_at = $7
		WHERE id = $8
		RETURNING
			id,
			first_name,
			last_name,
			phone_number,
			email,
			username,
			type,
			image_url,
			created_at,
			updated_at
	`
	var res repo.User
	err := ur.db.QueryRow(
		query,
		u.FirstName,
		u.LastName,
		utils.NullString(u.PhoneNumber),
		u.Email,
		utils.NullString(u.Username),
		utils.NullString(u.ImageUrl),
		time.Now(),
		u.ID,
	).Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&phoneNumber,
		&res.Email,
		&username,
		&imageUrl,
		&res.Type,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	res.PhoneNumber = phoneNumber.String
	res.ImageUrl = imageUrl.String
	res.Username = username.String

	return &res, nil
}

func (ur *userRepo) Delete(user_id int64) error {
	query := "UPDATE users SET deleted_at = $1 WHERE id = $2"

	result, err := ur.db.Exec(query, time.Now(), user_id)
	if err != nil {
		return err
	}
	if res, _ := result.RowsAffected(); res == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (ur *userRepo) GetAll(params *repo.GetAllUsersParams) (*repo.GetAllUsersResult, error) {
	var phoneNumber, username, imageUrl sql.NullString
	var updatedAt sql.NullTime
	var res repo.GetAllUsersResult
	res.Users = make([]*repo.User, 0)
	filter := " WHERE deleted_at IS NULL "
	offset := (params.Page - 1) * params.Limit
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(` 
		 AND first_name ILIKE '%s' OR 
		last_name ILIKE '%s' OR 
		phone_number ILIKE '%s' OR 
		email ILIKE '%s' OR username ILIKE '%s' `, str, str, str, str, str)
	}

	orderBy := " ORDER BY created_at DESC "
	if params.SortBy != "" {
		orderBy = fmt.Sprintf(" ORDER BY created_at %s ", params.SortBy)
	}

	query := `
		SELECT 
			id,
			first_name,
			last_name,
			phone_number,
			email,
			username,
			image_url,
			type,
			created_at,
			updated_at
		FROM users 
	` + filter + orderBy + limit

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var u repo.User
		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&phoneNumber,
			&u.Email,
			&username,
			&imageUrl,
			&u.Type,
			&u.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		u.PhoneNumber = phoneNumber.String
		u.ImageUrl = imageUrl.String
		u.Username = username.String
		u.UpdatedAt = updatedAt.Time

		res.Users = append(res.Users, &u)
	}
	queryCount := "SELECT count(*) FROM users " + filter

	err = ur.db.QueryRow(queryCount).Scan(&res.Count)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (ur *userRepo) GetByEmail(email string) (*repo.User, error) {
	var phoneNumber, username, imageUrl sql.NullString
	var updatedAt sql.NullTime
	query := `
		SELECT 
			id,
			first_name,
			last_name,
			phone_number,
			email,
			password,
			username,
			image_url,
			type,
			created_at,
			updated_at
		FROM users WHERE email = $1
	`
	var res repo.User
	err := ur.db.QueryRow(
		query,
		email,
	).Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&phoneNumber,
		&res.Email,
		&res.Password,
		&username,
		&imageUrl,
		&res.Type,
		&res.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	res.PhoneNumber = phoneNumber.String
	res.ImageUrl = imageUrl.String
	res.Username = username.String
	res.UpdatedAt = updatedAt.Time

	return &res, nil
}
