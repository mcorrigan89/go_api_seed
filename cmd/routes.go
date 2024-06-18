package main

import (
	"context"
	"net/http"
	"time"

	"corrigan.io/go_api_seed/graph"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (app *application) graphqlHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(app.services, app.logger)}))

		h.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
			err := graphql.DefaultErrorPresenter(ctx, e)
			app.logger.Error().Ctx(ctx).AnErr("error", err).Msg("GRAPHQL ERROR")

			return err
		})

		h.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
			oc := graphql.GetOperationContext(ctx)

			res := next(ctx)
			if oc.Operation.Name == "__ApolloGetServiceDefinition__" {
				return res
			}
			if oc.OperationName != "IntrospectionQuery" {
				app.logger.Info().
					Ctx(ctx).
					Str("operation", oc.OperationName).
					Str("operationTime", time.Since(oc.Stats.OperationStart).String()).
					Msg("GRAPHQL")
			}

			return res
		})

		h.ServeHTTP(w, r)
	})
}

func (app *application) playgroundHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := playground.Handler("GraphQL", "/graphql")
		h.ServeHTTP(w, r)
	})
}

func (app *application) routes() http.Handler {

	r := chi.NewRouter()

	r.Get("/graphql", app.playgroundHandler())
	r.Post("/graphql", app.graphqlHandler())
	r.Get("/ping", ping)

	return app.recoverPanic(app.enabledCORS(app.contextBuilder(r)))
}
