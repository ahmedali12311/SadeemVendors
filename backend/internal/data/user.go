package data

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"project/utils"
	"project/utils/validator"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID         uuid.UUID `db:"id"         json:"id"`
	Name       string    `db:"name"       json:"name"`
	Email      string    `db:"email"      json:"email"`
	Phone      string    `db:"phone"      json:"phone"`
	Img        *string   `db:"img"        json:"img"`
	Password   string    `db:"password"   json:"-"`
	Created_at time.Time `db:"created_at" json:"created_at"`
	Updated_at time.Time `db:"updated_at" json:"updated_at"`
}

<<<<<<< HEAD
type UserDB struct {
=======
type userDB struct {
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	db *sqlx.DB
}

func ValidatingUser(v *validator.Validator, user *User, fields ...string) {
	for _, field := range fields {
		switch field {
		case "name":
			if user.Name != "" {
				v.Check(len(user.Name) <= 20, "name", "Name must be less than 20")
				v.Check(len(user.Name) >= 3, "name", "Name must be more than 3")

			}
		case "phone":
			if user.Phone != "" {
				v.Check(validator.Matches(user.Phone, validator.PhoneRX), "phone", "Invalid phone number")
			}
		case "email":

			if user.Email != "" {
				v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "Invalid email format")
			}

		case "password":

			if user.Password != "" {
<<<<<<< HEAD
=======
				// Example: Validate password strength, e.g., minimum length
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
				v.Check(len(user.Password) >= 8, "password", "Password too short")
			}

		}
	}
}
<<<<<<< HEAD
func (u *UserDB) GetUsers(sortColumn, sortDirection string, page, pageSize int, searchTerm string) (*[]User, error) {
=======
func (u *userDB) GetUsers(sortColumn, sortDirection string, page, pageSize int, searchTerm string) (*[]User, error) {
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	var users []User

	// Construct query builder
	queryBuilder := QB.Select(strings.Join(user_columns, ",")).From("users")

	// Apply search filter
	if searchTerm != "" {
		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+searchTerm+"%")
	}

	// Apply sorting
	if sortDirection == "" {
		sortDirection = "ASC" // Default sort direction
	}
	validSortColumns := map[string]bool{"name": true, "created_at": true}
	if !validSortColumns[sortColumn] {
		sortColumn = "created_at" // Default sort column
	}
	queryBuilder = queryBuilder.OrderBy(sortColumn + " " + sortDirection)

	// Apply pagination
	queryBuilder = queryBuilder.Limit(uint64(pageSize)).Offset(uint64((page - 1) * pageSize))

	// Build and execute query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	err = u.db.Select(&users, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &users, nil
}
<<<<<<< HEAD
func (u *UserDB) GetUser(id uuid.UUID) (*User, error) {
	var user User

	// Construct the SQL query using squirrel
	query, args, err := QB.Select(strings.Join(user_columns, ",")).From("users").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	// Execute the query and map the result to the User struct
	err = u.db.QueryRowx(query, args...).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound // Return custom error if no record is found
		}
		return nil, fmt.Errorf("error scanning result: %w", err)
=======

func (u *userDB) GetUser(id uuid.UUID) (*User, error) {
	var user User
	query, agrs, err := QB.Select(strings.Join(user_columns, ",")).From("users").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	err = u.db.QueryRowx(query, agrs...).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}

		return nil, err
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	}

	return &user, nil
}
<<<<<<< HEAD

func (u *UserDB) Insert(user *User) error {
=======
func (u *userDB) Insert(user *User) error {
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	err := u.checkEmailExists(user.Email)
	if err != nil {
		if err == ErrDuplicatedKey {
			return ErrDuplicatedKey
		}
		return fmt.Errorf("error checking email existence: %w", err)
	}

	query, args, err := QB.
		Insert("users").
		Columns("img", "name", "phone", "email", "password").
		Values(user.Img, user.Name, user.Phone, user.Email, user.Password).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(user_columns, ", "))).
		ToSql()
	if err != nil {
		return err
	}

	err = u.db.QueryRowx(query, args...).StructScan(user)

	if err != nil {
		switch {
		case errors.Is(err, ErrDuplicatedKey):
			return ErrDuplicatedKey
		default:
			return err
		}
	}

	return nil
}
<<<<<<< HEAD
func (u *UserDB) DeleteUser(id uuid.UUID) (*User, error) {
=======
func (u *userDB) DeleteUser(id uuid.UUID) (*User, error) {
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c

	var user User
	query, args, err := QB.Delete("users").Where(squirrel.Eq{"id": id}).Suffix(fmt.Sprintf("RETURNING %s", strings.Join(user_columns, ", "))).ToSql()
	if err != nil {
		return nil, err
	}
	err = u.db.QueryRowx(query, args...).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err

	}

	if user.Img != nil {
		imgfile := strings.TrimPrefix(*user.Img, Domain+"/")
		// Check if the file exists
		if _, err := os.Stat(imgfile); err == nil {
			// File exists, attempt to delete it
			err = utils.DeleteImageFile(imgfile)
			if err != nil {
				return nil, fmt.Errorf("failed to delete file %s: %w", imgfile, err)
			}
		} else if os.IsNotExist(err) {
			// File does not exist, log a warning but do not treat it as a fatal error
			fmt.Printf("Warning: image file %s does not exist\n", imgfile)
		} else {
			// Handle other potential errors from os.Stat
			return nil, fmt.Errorf("failed to check file %s: %w", imgfile, err)
		}
	}
	return &user, nil
}
<<<<<<< HEAD
func (u *UserDB) Update(user *User) error {
=======
func (u *userDB) Update(user *User) error {
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	originalUser, err := u.GetUser(user.ID)
	if err != nil {
		return err
	}

	if user.Email != originalUser.Email {
		err := u.checkEmailExists(user.Email)
		if err != nil {
			return err
		}
	}
	query, args, err := QB.Update("users").
		Set("img", &user.Img).
		Set("name", &user.Name).
		Set("email", &user.Email).
		Set("phone", &user.Phone).
		Set("password", &user.Password).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id ": user.ID}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(user_columns, ", "))).
		ToSql()

	if err != nil {

		return err
	}

	result, err := u.db.Exec(query, args...)
	if err != nil {
		return err

	}
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		return err

	}
	if rowsaffected == 0 {
		return ErrRecordNotFound
	}

	return nil

}
<<<<<<< HEAD
func (u *UserDB) GetUserByEmail(email string) (*User, error) {
=======
func (u *userDB) GetUserByEmail(email string) (*User, error) {
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	// Construct SQL query
	var user User
	query, args, err := QB.Select(strings.Join(user_columns, ", ")).From("users").Where(squirrel.Eq{"email": email}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	err = u.db.QueryRowx(query, args...).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			// No record found
<<<<<<< HEAD
			return nil, ErrUserNotFound
=======
			return nil, ErrRecordNotFound
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
		}
		// Some other error occurred
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	// Email exists
	return &user, nil
}
<<<<<<< HEAD
func (u *UserDB) checkEmailExists(email string) error {
=======
func (u *userDB) checkEmailExists(email string) error {
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	query, args, err := QB.Select("1").From("users").Where(squirrel.Eq{"email": email}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	var exists bool
	err = u.db.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("failed to execute query: %v", err)
	}

	if exists {
		return ErrDuplicatedKey
	}
	return nil
}
