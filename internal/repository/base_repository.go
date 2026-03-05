
package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type IBaseRepository interface {
	Query(query string) (interface{}, error)
}

type baseRepository struct {
	ctx    context.Context
	client *pgxpool.Pool
}

// NewBaseRepository initializes repository with proper context
func NewBaseRepository(ctx context.Context, client *pgxpool.Pool) (IBaseRepository, error) {
	return &baseRepository{
		ctx:    ctx,
		client: client,
	}, nil
}

func (r *baseRepository) Query(query string) (interface{}, error) {
	datas := make([]map[string]interface{}, 0)

	rows, err := r.client.Query(context.TODO(), query)
	if err != nil {
		return datas, fmt.Errorf("h.db.Query|%s", err.Error())
	}
	defer rows.Close()

	columns := rows.FieldDescriptions()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		data := map[string]interface{}{}
		for i, col := range columns {
			val := values[i]

			b, ok := val.([]byte)
			var v interface{}
			if ok {
				v = string(b)
			} else {
				v = val
			}

			data[string(col.Name)] = v
		}

		datas = append(datas, data)
	}

	return datas, nil
}

