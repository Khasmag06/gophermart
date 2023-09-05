package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Khasmag06/gophermart/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

var (
	ErrExistsData               = errors.New("exists data")
	ErrOrderInvalidNum          = errors.New("invalid order number")
	ErrUserCredentials          = errors.New("incorrect login or password")
	ErrOrderUploadedByUser      = errors.New("order uploaded user")
	ErrOrderUploadedByOtherUser = errors.New("order uploaded by other user")
	ErrBalanceNotEnoughPoints   = errors.New("not enough points")
)

type Storage interface {
	AddUser(user *models.User) error
	Login(user *models.User) error
	AddOrder(userID int, orderNum string) error
	GetOrders(userID int) ([]*models.Order, error)
	GetSumAccruals(userID int) (float64, error)
	GetSumWithdrawals(userID int) (float64, error)
	NewWithdrawal(userID int, withdraws *models.Withdraws) error
	GetWithdrawals(userID int) ([]*models.Withdraws, error)
	GetBalance(userID int) (*models.JSONBalance, error)
	UpdateAccrual(accrual models.AccrualData) error
}

type DBStorage struct {
	db *sql.DB
}

func (dbs *DBStorage) AddUser(user *models.User) error {
	_, err := dbs.db.Exec("INSERT INTO users (login, password) VALUES ($1, $2) ", user.Login, user.Password)
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("user already exists: %w", ErrExistsData)
		}
		return err
	}
	row := dbs.db.QueryRow("SELECT id FROM users WHERE login = $1 and password = $2", user.Login, user.Password)
	if err := row.Scan(&user.ID); err != nil {
		return err
	}
	return nil
}

func (dbs *DBStorage) Login(user *models.User) error {
	row := dbs.db.QueryRow("SELECT id, password  FROM users WHERE login = $1", user.Login)
	if err := row.Scan(&user.ID, &user.Password); err != nil {
		return ErrUserCredentials
	}
	return nil
}

func (dbs *DBStorage) AddOrder(userID int, orderNum string) error {
	query := "INSERT INTO orders (order_num, user_id, status) VALUES ($1, $2, $3)"
	_, err := dbs.db.Exec(query, orderNum, userID, "NEW")
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			var id int
			query = "SELECT user_id FROM orders WHERE order_num = $1"
			if err := dbs.db.QueryRow(query, orderNum).Scan(&id); err != nil {
				return err
			}
			if userID != id {
				return ErrOrderUploadedByOtherUser
			} else {
				return ErrOrderUploadedByUser
			}
		}
		return err
	}
	return nil
}

func (dbs *DBStorage) GetOrders(userID int) ([]*models.Order, error) {
	query := "SELECT order_num, status, accrual, uploaded_at FROM orders WHERE user_id = $1  ORDER BY uploaded_at DESC"
	rows, err := dbs.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := make([]*models.Order, 0)

	for rows.Next() {
		var order models.Order
		var uploadedAt time.Time
		err = rows.Scan(&order.OrderNum, &order.Status, &order.Accrual, &uploadedAt)
		if err != nil {
			return nil, err
		}
		order.UploadedAt = uploadedAt.Format(time.RFC3339)
		order.Accrual = order.Accrual / 100
		orders = append(orders, &order)

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (dbs *DBStorage) GetSumAccruals(userID int) (float64, error) {
	row := dbs.db.QueryRow("SELECT coalesce(SUM(accrual), 0) AS points FROM orders WHERE user_id = $1", userID)
	var sumAccruals float64
	if err := row.Scan(&sumAccruals); err != nil {
		return 0, err
	}
	return sumAccruals / 100, nil
}

func (dbs *DBStorage) GetSumWithdrawals(userID int) (float64, error) {
	row := dbs.db.QueryRow("SELECT coalesce(SUM(sum), 0) FROM withdrawals WHERE user_id = $1", userID)
	var sumWithdrawals float64
	if err := row.Scan(&sumWithdrawals); err != nil {
		return 0, err
	}
	return sumWithdrawals / 100, nil
}

func (dbs *DBStorage) NewWithdrawal(userID int, withdraws *models.Withdraws) error {
	accruals, err := dbs.GetSumAccruals(userID)
	if err != nil {
		return err
	}
	withdrawals, err := dbs.GetSumWithdrawals(userID)
	if err != nil {
		return err
	}
	points := accruals - withdrawals

	if withdraws.Sum > points {
		return ErrBalanceNotEnoughPoints
	}

	query := "INSERT INTO withdrawals (user_id, order_num, sum) VALUES ($1, $2, $3)"
	_, err = dbs.db.Exec(query, userID, withdraws.Order, withdraws.Sum*100)
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DBStorage) GetWithdrawals(userID int) ([]*models.Withdraws, error) {
	query := "SELECT order_num, sum, processed_at FROM withdrawals WHERE user_id = $1 ORDER BY processed_at DESC"
	rows, err := dbs.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	withdrawals := make([]*models.Withdraws, 0)
	for rows.Next() {
		var withdrawal models.Withdraws
		var processedAt time.Time
		err := rows.Scan(&withdrawal.Order, &withdrawal.Sum, &processedAt)
		if err != nil {
			return nil, err
		}
		withdrawal.ProcessedAt = processedAt.Format(time.RFC3339)
		withdrawal.Sum = withdrawal.Sum / 100
		withdrawals = append(withdrawals, &withdrawal)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (dbs *DBStorage) GetBalance(userID int) (*models.JSONBalance, error) {
	accruals, err := dbs.GetSumAccruals(userID)
	if err != nil {
		return nil, err
	}
	withdrawals, err := dbs.GetSumWithdrawals(userID)
	if err != nil {
		return nil, err
	}
	currentBalance := accruals - withdrawals
	var balance models.JSONBalance
	balance.Current = currentBalance
	balance.Withdrawn = withdrawals
	return &balance, nil

}

func (dbs *DBStorage) UpdateAccrual(accrual models.AccrualData) error {
	_, err := dbs.db.Exec(
		"UPDATE orders SET status=$1, accrual=$2 WHERE order_num = $3",
		accrual.Status, accrual.Accrual*100, accrual.Order,
	)
	return err
}

func NewDB(dsn string) (Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to open sql connection: %w", err)
	}

	for _, table := range tables {
		_, err = db.Exec(table)
		if err != nil {
			return nil, err
		}
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DBStorage{
		db: db,
	}, nil
}
