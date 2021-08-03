package graph

import (
	"github.com/shien/restserver/graphql/taskstore"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Store *taskstore.TaskStore
}
