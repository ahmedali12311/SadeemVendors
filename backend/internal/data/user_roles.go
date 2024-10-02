package data

import (
	"database/sql"
	"errors"
	"fmt"
	"project/utils/validator"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User_role struct {
	UserID uuid.UUID `db:"user_id" json:"userID"`
	RoleID int       `db:"role_id" json:"roleID"`
}

<<<<<<< HEAD
type UserRoleDB struct {
	db *sqlx.DB
}

func ValidatingUserRole(v *validator.Validator, roles int) {
	v.Check(roles >= 1 && roles <= 3, "role", "Role not found")

	v.Check(roles == 1, "role", "You don't not have the required role to perform this operation")
}
func (r *UserRoleDB) GrantRole(user uuid.UUID, role int) (*User_role, error) {
=======
type user_roleDB struct {
	db *sqlx.DB
}

func ValidatingUserRole(v *validator.Validator, roles *User_role) {
	v.Check(roles.RoleID >= 1, "role", "Role not found")
	v.Check(roles.RoleID <= 3, "role", "Role not found")

	v.Check(roles.RoleID == 1, "role", "You don't not have the required role to perform this operation")
}
func (r *user_roleDB) GrantRole(user uuid.UUID, role int) (*User_role, error) {
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c

	existingRole, err := r.GetUserRole(user)

	if err != nil && !errors.Is(err, ErrRecordNotFound) {
		return nil, err
	}

	if existingRole != nil {
		return nil, ErrHasRole
	}

	// Step 2: Insert the new role
	var user_role User_role
	query, args, err := QB.Insert("user_roles").Columns(strings.Join(user_roles, ",")).Values(user, role).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(user_roles, ","))).
		ToSql()
	if err != nil {

		return nil, err
	}

	// Execute the insert query
	err = r.db.QueryRowx(query, args...).StructScan(&user_role)
	if err != nil {
		return nil, err
	}

	return &user_role, nil
}
<<<<<<< HEAD
func (r *UserRoleDB) UpdateRole(userID uuid.UUID, newRoleID int) (*User_role, error) {
=======
func (r *user_roleDB) UpdateRole(userID uuid.UUID, newRoleID int) (*User_role, error) {
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	updatedUserRole := User_role{
		UserID: userID,
		RoleID: newRoleID,
	}

	// Check if the user already has a role
	existingRole, err := r.GetUserRole(userID)
	if err != nil && errors.Is(err, ErrRecordNotFound) {
		MakeRole, err := r.GrantRole(userID, 3)
		updatedUserRole = *MakeRole
		if err != nil {
			return nil, err
		}
		return &updatedUserRole, nil // Return the newly granted role
	}
	if err != nil {
		return nil, err
	}

	// If the existing role is the same as the new role, return a duplication error
	if existingRole.RoleID == newRoleID {
		return nil, ErrDuplicatedRole
	}

	// Update the role if it's different
	query, args, err := QB.Update("user_roles").
		Set("role_id", newRoleID).
		Where(squirrel.Eq{"user_id": userID}).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowx(query, args...).StructScan(&updatedUserRole)
	if err != nil {
		return nil, err
	}

	return &updatedUserRole, nil
}
<<<<<<< HEAD
func (r *UserRoleDB) RevokeRole(user uuid.UUID, role int) error {

	_, err := r.GetUserRole(user)
=======
func (r *user_roleDB) RevokeRole(user uuid.UUID, role int) error {

	// Step 2: Check if the user has any remaining roles
	_, err := r.GetUserRole(user) // Assuming GetUserRoles returns all roles for a user
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	if err != nil {
		return err
	}
	query, args, err := QB.Delete("user_roles").Where(squirrel.And{
		squirrel.Eq{"user_id": user},
		squirrel.Eq{"role_id": role},
	}).ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

<<<<<<< HEAD
func (r *UserRoleDB) GetUserRole(id uuid.UUID) (*User_role, error) {
=======
func (r *user_roleDB) GetUserRole(id uuid.UUID) (*User_role, error) {
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	var userRole User_role
	query, args, err := QB.Select("user_id", "role_id").From("user_roles").Where(squirrel.Eq{"user_id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.Get(&userRole, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &userRole, nil
}

<<<<<<< HEAD
func (r *UserRoleDB) GetUserRoles() (*[]User_role, error) {
=======
func (r *user_roleDB) GetUserRoles() (*[]User_role, error) {
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	var users_roles []User_role
	query, args, err := QB.Select(strings.Join(user_roles, ",")).From("user_roles").ToSql()
	if err != nil {

		return nil, err
	}
	err = r.db.Select(&users_roles, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}

	return &users_roles, nil
}
