package http

import "net/http"

type Router struct {
	Todo *TodoHandler
}

func NewRouter(todo *TodoHandler) *Router {
	return &Router{Todo: todo}
}

func (r *Router) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /todos", r.Todo.CreateTodo)
	mux.HandleFunc("GET /todos/{id}", r.Todo.GetTodo)
	mux.HandleFunc("PUT /todos/{id}", r.Todo.UpdateTodo)
	mux.HandleFunc("DELETE /todos/{id}", r.Todo.DeleteTodo)
	mux.HandleFunc("GET /todos", r.Todo.GetTodoList)

	var handler http.Handler = mux
	handler = r.withLogger(handler)
	handler = r.withRecovery(handler)

	return handler
}

func (r *Router) withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, req)
	})
}

func (r *Router) withLogger(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		println(req.Method, req.URL.Path)
		nextHandler.ServeHTTP(w, req)
	})
}
