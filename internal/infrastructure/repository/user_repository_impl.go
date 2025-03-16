package repository

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/masahiro0000/BizMatch/internal/domain"
)

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepositoryImpl(db *sqlx.DB) domain.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

// CreateUser inserts a new user into the database.
func (r *userRepositoryImpl) CreateUser(user *domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, display_name, password) VALUES ($1, $2, $3)",
						user.Username, user.DisplayName, user.Password)
	if err != nil {
		log.Printf("Failed to insert user. username:%v, display name:%v, error:%v", user.Username, user.DisplayName, err)
		var pqErr *pq.Error
		// Check if the error is a PostgreSQL error for duplicate entries.
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

// GetAllUsers retrieves all users from the database.
func (r *userRepositoryImpl) GetAllUsers() ([]*domain.User, error) {
	var users []*domain.User
	query := "SELECT * FROM users"

	if err := r.db.Select(&users, query); err != nil {
		log.Printf("Failed to get all users. error:%v", err)
		return nil, err
	}
	return users, nil
}

// GetUserByID retrieve a user from the database based on user ID.
func (r *userRepositoryImpl) GetUserByID(userID int64) (*domain.User, error) {
	var user domain.User
	query := "SELECT * FROM users WHERE id = $1"

	// Execute the SQL query using userID and map the result to the user variable.
	if err := r.db.Get(&user, query, userID); err != nil {
		log.Printf("Failed to get user by ID. userID:%v, error:%v", userID, err)
		return nil, err
	}

	// If successful, return a pointer to the user.
	return &user, nil
}

// GetUserByUsername retrieve a user from the database based on username.
func (r *userRepositoryImpl) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	query := "SELECT * FROM users WHERE username = $1"

	// Execute the SQL query using username and map the result to the user variable.
	if err := r.db.Get(&user, query, username); err != nil {
		log.Printf("Failed to get user by username. username:%v, error:%v", username, err)
		// If no rows are returned, return an error indicated the user was not found.
		return nil, err
	}

	// If successful, return a pointer to the user.
	return &user, nil
}

// RegisterInfo updates the user's information in the database.
func (r *userRepositoryImpl) RegisterInfo(user *domain.User) error {
	// SQL query that update the user's information.
	query := `
		UPDATE users
		SET
			username 	 		= :username,
			display_name		= :display_name,
			prefecture_id		= :prefecture_id,
			industry_id 		= :industry_id,
			job_id				= :job_id,
			position_id 		= :position_id,
			age					= :age,
			gender				= :gender,
			photo				= :photo,
			profile_description	= :profile_description
		WHERE
			id = :id
		`

	// Execute the SQL update query.
	_, err := r.db.NamedExec(query, user)
	if err != nil {
		log.Printf("Failed to update user's information. query:%v, error:%v", query, err)
		return err
	}

	return nil
}

// UpdatePassword updates the user's password in the database.
func (r *userRepositoryImpl) UpdatePassword(userID int64, hashedPassword string) error {
	query := "UPDATE users SET password = $1 WHERE id = $2"
	_, err := r.db.Exec(query, hashedPassword, userID)

	if err != nil {
		log.Printf("Fail to update password. userID:%v, error:%v", userID, err)
		return err
	}

	return nil
}

// SearchUsers retrieves users from the database based on the provided filter criteria.
func (r *userRepositoryImpl) SearchUsers(filter *domain.UserSearchFilter) ([]*domain.User, error) {
	queries := []string{
		`SELECT
			u.id,
			u.username,
			u.display_name,
			u.password,
			u.prefecture_id,
			u.industry_id,
			u.job_id,
			u.position_id,
			u.age,
			u.gender,
			u.photo,
			u.profile_description,
			p.prefecture,
			i.industry,
			j.job,
			pos.position
		FROM users AS u
		LEFT JOIN prefectures AS p
			ON u.prefecture_id = p.id
		LEFT JOIN industries AS i
			ON u.industry_id = i.id
		LEFT JOIN jobs AS j
			ON u.job_id = j.id
		LEFT JOIN positions AS pos
			ON u.position_id = pos.id
		WHERE 1=1`,
	}

	// args will collect the parameters to be substituted into the query.
	var args []interface{}

	// Add conditions if provided in the filter.
	if len(filter.Prefectures) > 0 {
		queries = append(queries, " AND u.prefecture_id IN (?)")
		args = append(args, filter.Prefectures)
	}

	if len(filter.Industries) > 0 {
		queries = append(queries, " AND u.industry_id IN (?)")
		args = append(args, filter.Industries)
	}

	if len(filter.Jobs) > 0 {
		queries = append(queries, " AND u.job_id IN (?)")
		args = append(args, filter.Jobs)
	}

	if len(filter.Positions) > 0 {
		queries = append(queries, " AND u.position_id IN (?)")
		args = append(args, filter.Positions)
	}

	// If age groups are specified, generate OR conditions for each age range(start to start+9).
	if len(filter.AgeGroups) > 0 {
		var orConds []string
		for _, start := range filter.AgeGroups {
			end := start + 9
			cond := fmt.Sprintf("(u.age >= %d AND u.age <= %d)", start, end)
			orConds = append(orConds, cond)
		}
		// Combine all the OR conditions into a single string and appear to the queries.
		queries = append(queries, " AND (" + strings.Join(orConds, " OR ") + ")")
	}

	// Add conditions if provided in the filter.
	if len(filter.Genders) > 0 {
		queries = append(queries, " AND u.gender IN (?)")
		args = append(args, filter.Genders)
	}

	// Exclude a specific user ID if ExcludeUserID is set.
	if filter.ExcludeUserID != nil {
		queries = append(queries, " AND u.id <> ?")
		args = append(args, filter.ExcludeUserID)
	}

	query := strings.Join(queries, "")

	// Use sqlx.In to expand slice parameters (like IN clauses) into the query correctly.
	expandedQuery, expandedArgs, err := sqlx.In(query, args...)
	if err != nil {
		log.Printf("Failed to expand query. query:%v, error:%v", query, err)
		return nil, err
	}

	expandedQuery = r.db.Rebind(expandedQuery)

	// Execute the query with the expanded arguments.
	rows, err := r.db.Queryx(expandedQuery, expandedArgs...)
	if err != nil {
		log.Printf("Failed to search users. query:%v, error:%v", expandedQuery, err)
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	// Iterate over the result rows and map each row to a User object.
	for rows.Next() {
		var u domain.User
		if err := rows.StructScan(&u); err != nil {
			log.Printf("Failed to scan user. error:%v", err)
			return nil, err
		}
		users = append(users, &u)
	}
	// Return the list of matching users.
	return users, nil
}
