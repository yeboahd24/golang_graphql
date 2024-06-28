package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
)

// User struct to represent our data model
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func main() {
	// Connect to PostgreSQL
	var err error
	db, err = sql.Open("postgres", "postgresql://postgres:postgres@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Define GraphQL schema
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":   &graphql.Field{Type: graphql.Int},
			"name": &graphql.Field{Type: graphql.String},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: getUser,
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: createUser,
			},
			"updateUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: updateUser,
			},
			"deleteUser": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: deleteUser,
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	// Set up Gin router
	r := gin.Default()
	r.POST("/graphql", func(c *gin.Context) {
		var request struct {
			Query string `json:"query"`
		}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: request.Query,
		})

		c.JSON(http.StatusOK, result)
	})

	r.Run(":8080")
}

func getUser(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, nil
	}

	var user User
	err := db.QueryRow("SELECT id, name FROM test_users WHERE id = $1", id).Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func createUser(p graphql.ResolveParams) (interface{}, error) {
	name, _ := p.Args["name"].(string)

	var user User
	err := db.QueryRow("INSERT INTO users(name) VALUES($1) RETURNING id, name", name).Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func updateUser(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args["id"].(int)
	name, _ := p.Args["name"].(string)

	var user User
	err := db.QueryRow("UPDATE users SET name = $1 WHERE id = $2 RETURNING id, name", name, id).Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func deleteUser(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args["id"].(int)

	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	return true, nil
}
