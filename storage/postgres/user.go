package postgres

import (
	"context"
	"database/sql"
	"errors"
	models "exam3/api/model"
	"exam3/pkg/check"
	"exam3/pkg/logger"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

type userRepo struct {
	db     *pgxpool.Pool
	logger logger.ILogger
}

func NewUser(db *pgxpool.Pool, log logger.ILogger) userRepo {
	return userRepo{
		db:     db,
		logger: log,
	}
}

func (c *userRepo) Create(ctx context.Context, user models.Users) (string, error) {
	// Validate the password
	if err := check.ValidatePassword(user.Password); err != nil {
		fmt.Println(user.Password)
		return "", err
	}

	// Validate the email
	if err := check.Validategmail(user.Mail); err != nil {
		return "", err
	}

	// Validate the phone number
	if err := check.ValidatePhone(user.Phone); err != nil {
		return "", err
	}

	// Generate hashed password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.logger.Error("failed to generate users new password", logger.Error(err))
		return "", err
	}

	// Generate new UUID for users ID
	id := uuid.New()

	// Prepare and execute SQL query to insert new user
	query := `INSERT INTO users (
        id,
        first_name,
        last_name,
        mail,
        phone,
        sex,
        active,
        password,
        created_at
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err = c.db.Exec(context.Background(), query,
		id.String(),
		user.First_name,
		user.Last_name,
		user.Mail,
		user.Phone,
		user.Sex,
		true,
		newHashedPassword,
		time.Now(),
	)

	if err != nil {
		fmt.Println(err)
		c.logger.Error("error while creating user:", logger.Error(err))
		return "", err
	}

	return id.String(), nil
}

func (c *userRepo) Update(ctx context.Context, user models.Users) (string, error) {
	if err := check.ValidatePassword(user.Password); err != nil {
		fmt.Println(user.Password)
		return "", err
	}

	if err := check.Validategmail(user.Mail); err != nil {
		return "", err
	}

	if err := check.ValidatePhone(user.Phone); err != nil {
		return "", err
	}

	query := `UPDATE users SET
        first_name = $1,
        last_name = $2,
        mail = $3,
        phone = $4,
        sex = $5,
        active = $6,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = $7 AND deleted_at =0`

	// Execute the query
	_, err := c.db.Exec(ctx, query,
		user.First_name,
		user.Last_name,
		user.Mail,
		user.Phone,
		user.Sex,
		user.Active,
		user.Id)
	if err != nil {
		c.logger.Error("Error while updating user:", logger.Error(err))
		return "", err
	}

	return user.Id, nil
}

func (c *userRepo) GetAlluser(ctx context.Context, req models.GetAllUsersRequest) (models.GetAllusersResponse, error) {
	resp := models.GetAllusersResponse{}
	filter := ""

	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND first_name ILIKE '%%%v%%' `, req.Search)
	}

	fmt.Println("filter: ", filter)

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter: ", filter)

	query := `
        SELECT 
            count(id) OVER(),
            id,
            first_name,
            last_name,
            mail,
            phone,
            sex,
            active,
            created_at,
            updated_at
        FROM 
		users 
        WHERE 
            deleted_at = 0 ` + filter

	rows, err := c.db.Query(context.Background(), query)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	var errScan error

	for rows.Next() {
		var user models.Users
		var createdAt, updatedAt time.Time

		errScan = rows.Scan(
			&resp.Count,
			&user.Id,
			&user.First_name,
			&user.Last_name,
			&user.Mail,
			&user.Phone,
			&user.Sex,
			&user.Active,
			&createdAt,
			&updatedAt,
		)
		if errScan != nil {
			c.logger.Error("error while scanning user info: ", logger.Error(errScan))
			return resp, errScan
		}

		user.Created_at = createdAt.Format("2006-01-02 15:04:05")
		user.Updated_at = updatedAt.Format("2006-01-02 15:04:05")

		fmt.Printf("Scanned user: %+v\n", user)

		resp.Users = append(resp.Users, user)
	}

	if err := rows.Err(); err != nil {
		c.logger.Error("error while iterating over rows: ", logger.Error(err))
		return models.GetAllusersResponse{}, err
	}

	fmt.Printf("Final response: %+v\n", resp)

	return resp, nil
}

func (c *userRepo) GetByID(ctx context.Context, id string) (models.Users, error) {
	query := `SELECT 
        id,
        first_name,
        last_name,
        mail,
        phone,
        sex,
        active
        FROM users WHERE id=$1 AND deleted_at=0`

	row := c.db.QueryRow(ctx, query, id)

	user := models.Users{}

	err := row.Scan(
		&user.Id,
		&user.First_name,
		&user.Last_name,
		&user.Mail,
		&user.Phone,
		&user.Sex,
		&user.Active,
	)
	if err != nil {
		fmt.Println(err)
		c.logger.Error("error while getting user by ID: ", logger.Error(err))
		return user, err
	}

	return user, nil
}

func (c *userRepo) Delete(ctx context.Context, id string) (string, error) {
	query := `UPDATE users SET 
        deleted_at = date_part('epoch', CURRENT_TIMESTAMP)::int
        WHERE id = $1 AND deleted_at IS NULL`

	_, err := c.db.Exec(ctx, query, id)

	if err != nil {
		c.logger.Error("failed to delete user from database", logger.Error(err))
		return id, err
	}

	return id, nil
}

func (c *userRepo) Login(ctx context.Context, login models.Changepasswor) (string, error) {
	var hashedPass string

	query := `SELECT password
	FROM users
	WHERE mail = $1 AND deleted_at = 0`

	err := c.db.QueryRow(ctx, query,
		login.Mail,
	).Scan(&hashedPass)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("incorrect login")
		}
		c.logger.Error("failed to get users password from database", logger.Error(err))
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(login.OldPassword))
	if err != nil {
		return "", errors.New("password mismatch")
	}

	return "Logged in successfully", nil
}

func (c *userRepo) GetPassword(ctx context.Context, phone string) (string, error) {
	var hashedPass string

	query := `SELECT password
	FROM users
	WHERE phone = $1 AND deleted_at = 0`

	err := c.db.QueryRow(ctx, query, phone).Scan(&hashedPass)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("incorrect phone")
		} else {
			c.logger.Error("failed to get users password from database", logger.Error(err))
			return "", err
		}
	}

	return hashedPass, nil
}

func (c *userRepo) ChangePassword(ctx context.Context, password models.Changepasswor) (string, error) {
	// Validate the new password
	if err := check.ValidatePassword(password.NewPassword); err != nil {
		fmt.Println(err, "")
		return "", err

	}

	query := `SELECT 
        password
        FROM users WHERE mail=$1 AND deleted_at=0`

	row := c.db.QueryRow(context.Background(), query, password.Mail)

	users := models.Changepasswor{}

	err := row.Scan(
		&users.OldPassword,
	)
	if err != nil {
		fmt.Println(err, "_===")
		c.logger.Error("error while getting users's old password: ", logger.Error(err))
		return users.OldPassword, err
	}

	// Compare the hash of the old password stored in the database with the provided old password

	// Generate a new hashed password for the new password provided
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.logger.Error("failed to generate new hashed password", logger.Error(err))
		return "", err
	}

	fmt.Println(newHashedPassword) // Print hashed password only when there's no error

	err = bcrypt.CompareHashAndPassword([]byte(newHashedPassword), []byte(password.OldPassword))
	if err == nil {
		fmt.Println(err, "")
		fmt.Println(newHashedPassword)
		return "", errors.New("password mismatch")
	}

	if len(password.NewPassword) < 8 {
		return "", errors.New("new password must be at least 8 characters long")
	}
	// Update the password in the database
	query = `UPDATE users SET 
        password = $1, 
        updated_at = CURRENT_TIMESTAMP 
    WHERE mail = $2 AND deleted_at = 0`

	_, err = c.db.Exec(ctx, query, newHashedPassword, password.Mail)
	if err != nil {
		fmt.Println(err, "===================")
		c.logger.Error("failed to change user password in database", logger.Error(err))
		return "", err
	}

	return "Password changed successfully", nil
}

func (c *userRepo) GetByMail(ctx context.Context, mail string) (models.Users, error) {
	var (
		firstname sql.NullString
		lastname  sql.NullString
		phone     sql.NullString
		email     sql.NullString
		createdat sql.NullString
		updatedat sql.NullString
	)

	query := `SELECT 
        id, 
        first_name, 
        last_name, 
        phone,
        mail,
        created_at, 
        updated_at,
        password
        FROM users WHERE phone = $1 AND deleted_at IS NULL`

	row := c.db.QueryRow(ctx, query, mail)

	user := models.Users{}
	err := row.Scan(
		&user.Id,
		&firstname,
		&lastname,
		&phone,
		&mail,
		&createdat,
		&updatedat,
		&user.Password)

	if err != nil {
		c.logger.Error("failed to scan user by login from database", logger.Error(err))
		return models.Users{}, err
	}

	user.First_name = firstname.String
	user.Last_name = lastname.String
	user.Phone = phone.String
	user.Mail = email.String
	user.Created_at = createdat.String
	user.Updated_at = updatedat.String

	return user, nil
}

func (c *userRepo) Checklogin(ctx context.Context, gmail string) (models.UserRegisterRequest, error) {
	fmt.Println("Input Gmail:", gmail)

	query := `SELECT 
        gmail
       FROM users WHERE gmail=$1 AND deleted_at=0`

	row := c.db.QueryRow(ctx, query, gmail)

	users := models.UserRegisterRequest{}

	err := row.Scan(
		&users.Mail,
	)
	if err != nil {
		users.Mail = "iikkasdkifdfh@gmail.com"
		fmt.Println("Returning fixed email:", users.Mail)
		return users, nil
	}

	fmt.Println("Returning email from database:", users.Mail)
	return users, nil
}

func (c *userRepo) Createconf(ctx context.Context, userss models.Users) (string, error) {
	id := uuid.New()
	query := `INSERT INTO users (
        id,
        first_name,
        last_name,
        mail,
        phone,
        sex,
        active,
        password,
        created_at
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := c.db.Exec(context.Background(), query,
		id.String(),
		userss.First_name,
		userss.Last_name,
		userss.Mail,
		userss.Phone,
		userss.Sex,
		true,
		userss.Password,
		time.Now(),
	)

	if err != nil {
		fmt.Println(err)
		c.logger.Error("error while creating user:", logger.Error(err))
		return "", err
	}

	return id.String(), nil
}

func (c *userRepo) Updatestatus(ctx context.Context, user models.Updatestatus) (string, error) {
	query := `UPDATE users SET
       active = $1,
       updated_at = CURRENT_TIMESTAMP
    WHERE id = $2 AND deleted_at = 0`

	_, err := c.db.Exec(ctx, query)
	if err != nil {
		c.logger.Error("Error while updating user:", logger.Error(err))
		return "", err
	}

	return user.Id, nil
}
