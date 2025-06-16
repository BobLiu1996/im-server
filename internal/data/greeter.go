package data

import (
	"context"
	"im-server/internal/biz"
	"im-server/internal/biz/do"
)

type greeterImpl struct {
	data *Data
}

var _ biz.GreeterRepo = (*greeterImpl)(nil)

// NewGreeterRepo .
func NewGreeterRepo(data *Data) biz.GreeterRepo {
	return &greeterImpl{
		data: data,
	}
}

func (r *greeterImpl) ListAll(ctx context.Context) ([]*do.Greeter, error) {
	tGreeter := r.data.Query().TGreeter
	res := make([]*do.Greeter, 0)
	if err := tGreeter.WithContext(ctx).Select(tGreeter.Name, tGreeter.Age).Scan(&res); err != nil {
		return nil, err
	}
	return res, nil
}
