Cji
===

pronounced `chi`, is an inline middleware chain for Goji web apps.

chi allows you to use middlewares for a single route without creating a new handler with a SubRouter.

For instance

```
m.Use(someMiddleware)
m.Get("/one", handlerOne)
admin:= web.New()
m.Handle("/admin", admin)

admin.Use(middleware.SubRouter)
admin.Use(PasswordMiddleware)
admin.Use(HttpsOnlyMiddleware)
admin.Get("/", AdminRoot)
```

Becomes:

```
m.Use(someMiddleware)
m.Get("/one", handlerOne)
admin.Get("/", cji.Use(PasswordMiddleware, HttpsOnlyMiddleware).On(AdminRoot))
```

We find this useful for middlewares that lookup objects in the database and handle authorization

```
m.Get("/posts/:postId", cji.Use(PostContext).On(GetPost))

func HubContext(c *web.C, h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        user := c.Env["user"].(*data.User)
        postId = c.URLParams["postId"])
        post, err := //look up post in db, and make sure user has permissions
        if err != nil {
            w.WriteHeader(403)
            w.Write([]byte("Unauthorized"))
            return
        }
        c.Env["post"] = post
        h.ServeHTTP(w, r)
    }
    return http.HandlerFunc(fn)
}

func (h *FeedHandler) GetFeed(c web.C, w http.ResponseWriter, r *http.Request) {
    feed := c.Env["feed"].(*data.Feed)
    h.JSON(w, 200, feed)
}
```

Authors: @pkieltyka & @mveytsman
