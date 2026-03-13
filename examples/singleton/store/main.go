package main

import (
	"fmt"
	"log"

	"github.com/Fabien-Halaby/craft-go-patterns/pkg/singleton/store"
)

type User struct {
	ID   int64
	Name string
}

type Project struct {
	ID   int64
	Name string
}

type AppState struct {
	User 			*User
	IsAuthenticated bool
	CurrentProject 		*Project
	Theme 			string
}

func main() {
	s := store.GetInstance[AppState]()

	unsub := s.Subscribe(func (state AppState) {
		log.Printf("AUTH: %v, USER: %s, PROJECT: %s",
			state.IsAuthenticated,
			state.User.Name,
			state.CurrentProject.Name,
		)
	})
	defer unsub()

	s.SetState(func (st *AppState) {
		st.IsAuthenticated = true
		st.User = &User{ID: 1, Name: "Founder"}
		st.CurrentProject = &Project{ID: 42, Name: "MVP"}
		st.Theme = "dark"
	})

	final := s.GetState()
	fmt.Printf("THEME: %s\n", final.Theme)
}