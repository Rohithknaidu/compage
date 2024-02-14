package daos

import (
	"database/sql"
	"errors"
	"github.com/Rohithknaidu/Instagram/user-service/pkg/rest/server/daos/clients/sqls"
	"github.com/Rohithknaidu/Instagram/user-service/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
)

type UserDao struct {
	sqlClient *sqls.SQLiteClient
}

func migrateUsers(r *sqls.SQLiteClient) error {
	query := `
	CREATE TABLE IF NOT EXISTS users(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
        
		Age INTEGER NOT NULL,
		Name TEXT NOT NULL,
        CONSTRAINT id_unique_key UNIQUE (Id)
	)
	`
	_, err1 := r.DB.Exec(query)
	return err1
}

func NewUserDao() (*UserDao, error) {
	sqlClient, err := sqls.InitSqliteDB()
	if err != nil {
		return nil, err
	}
	err = migrateUsers(sqlClient)
	if err != nil {
		return nil, err
	}
	return &UserDao{
		sqlClient,
	}, nil
}

func (userDao *UserDao) CreateUser(m *models.User) (*models.User, error) {
	insertQuery := "INSERT INTO users(Age, Name)values(?, ?)"
	res, err := userDao.sqlClient.DB.Exec(insertQuery, m.Age, m.Name)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	m.Id = id

	log.Debugf("user created")
	return m, nil
}

func (userDao *UserDao) ListUsers() ([]*models.User, error) {
	selectQuery := "SELECT * FROM users"
	rows, err := userDao.sqlClient.DB.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var users []*models.User
	for rows.Next() {
		m := models.User{}
		if err = rows.Scan(&m.Id, &m.Age, &m.Name); err != nil {
			return nil, err
		}
		users = append(users, &m)
	}
	if users == nil {
		users = []*models.User{}
	}

	log.Debugf("user listed")
	return users, nil
}

func (userDao *UserDao) GetUser(id int64) (*models.User, error) {
	selectQuery := "SELECT * FROM users WHERE Id = ?"
	row := userDao.sqlClient.DB.QueryRow(selectQuery, id)
	m := models.User{}
	if err := row.Scan(&m.Id, &m.Age, &m.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}

	log.Debugf("user retrieved")
	return &m, nil
}

func (userDao *UserDao) UpdateUser(id int64, m *models.User) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}
	if id != m.Id {
		return nil, errors.New("id and payload don't match")
	}

	user, err := userDao.GetUser(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, sql.ErrNoRows
	}

	updateQuery := "UPDATE users SET Age = ?, Name = ? WHERE Id = ?"
	res, err := userDao.sqlClient.DB.Exec(updateQuery, m.Age, m.Name, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sqls.ErrUpdateFailed
	}

	log.Debugf("user updated")
	return m, nil
}

func (userDao *UserDao) DeleteUser(id int64) error {
	deleteQuery := "DELETE FROM users WHERE Id = ?"
	res, err := userDao.sqlClient.DB.Exec(deleteQuery, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sqls.ErrDeleteFailed
	}

	log.Debugf("user deleted")
	return nil
}
