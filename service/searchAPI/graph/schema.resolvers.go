package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.33

import (
	"YenExpress/config"
	"YenExpress/service/searchAPI/graph/model"

	// "YenExpress/service/dto"
	"context"
	"fmt"
)

// Drugs is the resolver for the Drugs field.
func (r *queryResolver) Drugs(ctx context.Context) ([]*model.Drug, error) {
	var drugs []*model.Drug
	err := config.DB.Find(&drugs).Error
	if err == nil {
		return drugs, nil
	}
	return []*model.Drug{}, err
	// panic(fmt.Errorf("not implemented: Drugs - Drugs"))
}

// DrugOrders is the resolver for the DrugOrders field.
func (r *queryResolver) DrugOrders(ctx context.Context) ([]*model.DrugOrder, error) {
	panic(fmt.Errorf("not implemented: DrugOrders - DrugOrders"))
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }