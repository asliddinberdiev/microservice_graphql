package catalog

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ErrNotFound = errors.New("Entity not found")
)

type Repository interface {
	Close()
	UpdateProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip, take uint64) ([]*Product, error)
	GetProductByIDs(ctx context.Context, ids []string) ([]*Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]*Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	Quantity    uint64 `json:"quantity"`
}

func NewElasticRepository(url string) (*elasticRepository, error) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	return &elasticRepository{client: client}, nil
}

func (r *elasticRepository) Close() {
	r.client.Stop()
}

func (r *elasticRepository) UpdateProduct(ctx context.Context, p Product) error {
	_, err := r.client.Index().Index("catalog").Type("product").Id(p.ID).
		BodyJson(productDocument{Name: p.Name, Description: p.Description, Price: p.Price, Quantity: p.Quantity}).
		Do(ctx)
	return err
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get().Index("catalog").Type("product").Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}

	if !res.Found {
		return nil, ErrNotFound
	}

	p := productDocument{}
	if err := json.Unmarshal(*res.Source, &p); err != nil {
		return nil, err
	}

	return &Product{ID: id, Name: p.Name, Description: p.Description, Quantity: p.Quantity}, nil
}

func (r *elasticRepository) GetProducts(ctx context.Context, skip, take uint64) ([]*Product, error) {
	res, err := r.client.Search().Index("catalog").Type("product").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*Product, 0)
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err := json.Unmarshal(*hit.Source, &p); err != nil {
			return nil, err
		}
		products = append(products, &Product{ID: hit.Id, Name: p.Name, Description: p.Description, Quantity: p.Quantity})
	}

	return products, nil
}

func (r *elasticRepository) GetProductByIDs(ctx context.Context, ids []string) ([]*Product, error) {
	items := make([]*elastic.MultiGetItem, 0)
	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().Index("catalog").Type("product").Id(id))
	}

	res, err := r.client.MultiGet().Add(items...).Do(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*Product, 0)
	for _, doc := range res.Docs {
		p := productDocument{}
		if err := json.Unmarshal(*doc.Source, &p); err != nil {
			return nil, err
		}
		products = append(products, &Product{ID: doc.Id, Name: p.Name, Description: p.Description, Quantity: p.Quantity})
	}
	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]*Product, error) {
	res, err := r.client.Search().Index("catalog").Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description", "quantity")).
		From(int(skip)).Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*Product, 0)
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err := json.Unmarshal(*hit.Source, &p); err != nil {
			return nil, err
		}
		products = append(products, &Product{ID: hit.Id, Name: p.Name, Description: p.Description, Quantity: p.Quantity})
	}

	return products, nil
}
