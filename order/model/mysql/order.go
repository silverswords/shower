package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Order struct {
	ID         uint32
	OrderCode  string    `json:"ordercode"`
	UserID     uint64    `json:"userid"`
	ShipCode   string    `json:"shipcode"`
	AddressID  string    `json:"addressid"`
	TotalPrice uint32    `json:"totalprice"`
	PayWay     uint8     `json:"payway"`
	Promotion  bool      `json:"promotion"`
	Freight    uint32    `json:"freight"`
	Status     uint8     `json:"status"`
	Created    time.Time `json:"created"`
	Closed     time.Time `json:"closed"`
	Updated    time.Time `json:"updated"`
}

type Item struct {
	ProductId uint32 `json:"productid"`
	OrderID   uint32 `json:"orderid"`
	Count     uint32 `json:"count"`
	Price     uint32 `json:"price"`
	Discount  uint32 `json:"discount"`
}

type OrmOrder struct {
	*Order
	Orm []*Item
}

const (
	orderDB = iota
	orderTable
	itemTable
	orderInsert
	itemInsert
	orderIdByOrderCode
	orderByOrderID
	itemsByOrderID
	orderListByUserID
	payByOrderID
	consignByOrderID
	statusByOrderID
)

var categorySQLFormatStr = []string{
	`CREATE DATABASE IF NOT EXISTS %s`,
	`CREATE TABLE IF NOT EXISTS %s(
				id INT UNSIGNED NOT NULL AUTO_INCREMENT ,
				orderCode VARCHAR(50) NOT NULL,
				userID BIGINT UNSIGNED NOT NULL,
				shipCode VARCHAR(50) NOT NULL DEFAULT '100000',
				addressID VARCHAR(20) NOT NULL,
				totalPrice INT UNSIGNED NOT NULL,
				payWay TINYINT UNSIGNED DEFAULT '0',
				promotion TINYINT(1) UNSIGNED DEFAULT '0',
				freight INT UNSIGNED NOT NULL,
				status TINYINT UNSIGNED DEFAULT '0',
				created DATETIME DEFAULT NOW(),
				closed DATETIME DEFAULT '8012-12-31 00:00:00',
				updated DATETIME DEFAULT NOW(),
				PRIMARY KEY (id),
				UNIQUE KEY orderCode (orderCode) USING BTREE,
				KEY created (created),
				KEY updated (updated),
				KEY status (status), 
				KEY payWay (payWay)
			)ENGINE=InnoDB AUTO_INCREMENT = 10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='order info'`,
	`CREATE TABLE IF NOT EXISTS %s(
				productID INT UNSIGNED NOT NULL,
				orderID VARCHAR(50) NOT NULL,
				count INT UNSIGNED NOT NULL,
				price INT UNSIGNED NOT NULL,
				discount TINYINT UNSIGNED NOT NULL,
				KEY orderID (orderID)
			)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='orderitem info'`,
	`INSERT INTO %s (orderCode,userID,addressID,totalPrice,promotion,freight,closed) VALUES(?,?,?,?,?,?,?)`,
	`INSERT INTO %s (productID,orderID,count,price,discount) VALUES(?,?,?,?,?)`,
	`SELECT id FROM %s WHERE orderCode = ? LOCK IN SHARE MODE`,
	`SELECT * FROM %s WHERE id = ? LOCK IN SHARE MODE`,
	`SELECT * FROM %s WHERE orderID = ? LOCK IN SHARE MODE`,
	`SELECT * FROM %s WHERE userID = ? AND status = ? LOCK IN SHARE MODE`,
	`UPDATE %s.%s SET payWay = ? , updated = ? , status = 2 WHERE id = ? LIMIT 1 `,
	`UPDATE %s.%s SET shipCode = ? , updated = ? , status = 3 WHERE id = ? LIMIT 1 `,
	`UPDATE %s.%s SET status = ? , updated = ? WHERE id = ? LIMIT 1 `,
}

// CreateDB -
func CreateDB(db *sql.DB, createDB string) error {
	sql := fmt.Sprintf(categorySQLFormatStr[orderDB], createDB)
	_, err := db.Exec(sql)
	return err
}

// CreateTable -
func CreateTable(db *sql.DB, ostore string) error {
	sql := fmt.Sprintf(categorySQLFormatStr[orderTable], ostore)
	_, err := db.Exec(sql)
	return err
}

// Insert -
func Insert(order Order, items []Item, db *sql.DB, closedInterval int, orderDB string, orderTable string, itemTable string) (id uint32, err error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	order.Closed = order.Created.Add(time.Duration(closedInterval * int(time.Hour)))
	ostore := orderDB + "." + orderTable
	ostoresql := fmt.Sprintf(categorySQLFormatStr[orderInsert], ostore)
	result, err := tx.Exec(ostoresql, order.OrderCode, order.UserID, order.AddressID, order.TotalPrice, order.Promotion, order.Freight, order.Closed)

	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[insert order] : insert order affected 0 rows")
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	order.ID = uint32(ID)
	istore := orderDB + "." + itemTable
	istoresql := fmt.Sprintf(categorySQLFormatStr[itemInsert], istore)
	for _, x := range items {
		result, err := tx.Exec(istoresql, x.ProductId, order.ID, x.Count, x.Price, x.Discount)
		if err != nil {
			return 0, err
		}
		if affected, _ := result.RowsAffected(); affected == 0 {
			return 0, errors.New("insert item: insert affected 0 rows")
		}

		// err = os.Cnf.User.UserCheck(tx, order.UserID, x.ProductId)
		// if err != nil {
		// 	return 0, err
		// }

		// err = os.Cnf.Stock.ModifyProductStock(tx, x.ProductId, int(x.Count))
		// if err != nil {
		// 	return 0, err
		// }
	}

	return order.ID, err
}

// OrderIDByOrderCode -
func OrderIDByOrderCode(db *sql.DB, ostore string, ordercode string) (uint32, error) {
	var (
		orderid uint32
		err     error
	)
	sql := fmt.Sprintf(categorySQLFormatStr[orderIdByOrderCode], ostore)
	rows, err := db.Query(sql, ordercode)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&orderid); err != nil {
			return 0, err
		}
	}

	return orderid, nil
}

func SelectByOrderKey(db *sql.DB, ostore, istore string, orderid uint32) (*OrmOrder, error) {
	var (
		oo OrmOrder
		o  Order
	)
	sql := fmt.Sprintf(categorySQLFormatStr[orderByOrderID], ostore)
	rows, err := db.Query(sql, orderid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&o.ID, &o.OrderCode, &o.UserID, &o.ShipCode, &o.AddressID, &o.TotalPrice, &o.PayWay, &o.Promotion, &o.Freight, &o.Status, &o.Created, &o.Closed, &o.Updated); err != nil {
			return nil, err
		}
	}
	lisitItemByOrderIdsql := fmt.Sprintf(categorySQLFormatStr[itemsByOrderID], istore)
	oo.Order = &o
	oo.Orm, err = LisitItemByOrderId(db, lisitItemByOrderIdsql, orderid)
	if err != nil {
		return nil, err
	}

	return &oo, nil
}

func LisitItemByOrderId(db *sql.DB, query string, orderid uint32) ([]*Item, error) {
	var (
		ProductId uint32
		OrderID   uint32
		Count     uint32
		Price     uint32
		Discount  uint32

		items []*Item
	)

	rows, err := db.Query(query, orderid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ProductId, &OrderID, &Count, &Price, &Discount); err != nil {
			return nil, err
		}

		item := &Item{
			ProductId: ProductId,
			OrderID:   OrderID,
			Count:     Count,
			Price:     Price,
			Discount:  Discount,
		}
		items = append(items, item)
	}

	return items, nil
}

func LisitOrderByUserID(db *sql.DB, ostore, istore string, userid uint64, mode uint8) ([]*OrmOrder, error) {
	var OOs []*OrmOrder
	orderListByUserIDSql := fmt.Sprintf(categorySQLFormatStr[orderListByUserID], ostore)
	rows, err := db.Query(orderListByUserIDSql, userid, mode)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var oo OrmOrder
		var o Order
		if err := rows.Scan(&o.ID, &o.OrderCode, &o.UserID, &o.ShipCode, &o.AddressID, &o.TotalPrice, &o.PayWay, &o.Promotion, &o.Freight, &o.Status, &o.Created, &o.Closed, &o.Updated); err != nil {
			return nil, err
		}
		lisitItemByOrderIDSql := fmt.Sprintf(categorySQLFormatStr[itemsByOrderID], istore)
		oo.Order = &o
		oo.Orm, err = LisitItemByOrderId(db, lisitItemByOrderIDSql, oo.ID)
		if err != nil {
			return nil, err
		}
		OOs = append(OOs, &oo)
	}

	return OOs, nil
}
