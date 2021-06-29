package datastore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	Models "github.com/tonymj76/cuvva-sample-code/models"
)

//MRepository holds the CRUD function to communicate to db
type MRepository interface {
	Create(context.Context, *Models.CreateRequest) (*Models.CreateResponse, error)
}

const merchantName = "merchants"

// Create creates a new merchant data
func (c *Connection) Create(ctx context.Context, req *Models.CreateRequest) (ms *Models.CreateResponse, errs error) {
	rows, err := c.SQLBuilder.Insert(
		merchantName,
	).SetMap(map[string]interface{}{
		"number_of_product": req.NumberOfProduct,
		"email":             req.Email,
		"business_name":     req.BusinessName,
	}).Suffix(
		`RETURNING id, number_of_product, email, business_name, created_at`,
	).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	c.Logger.Info("successfully created a merchant")
	defer func() {
		cerr := rows.Close()
		if err == nil && cerr != nil {
			logrus.WithError(err)
		}
	}()
	if rows.Next() {
		merchant, err := scanMerchantRecord(rows)
		if err != nil {
			return nil, err
		}
		return merchant, err
	}
	return
}

func scanMerchantRecord(row squirrel.RowScanner) (*Models.CreateResponse, error) {
	m := Models.CreateResponse{}
	err := row.Scan(
		&m.ID,
		&m.NumberOfProduct,
		&m.Email,
		&m.BusinessName,
		&m.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		logrus.Error("no sale record found", err.Error())
	}
	return &m, err
}
