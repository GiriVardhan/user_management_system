package main

import (
    "net/http"
    "github.com/gorilla/mux"
    common "./common"
)


func main() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", common.WelcomePage)
    router.HandleFunc("/login", common.LoginPage)
    router.HandleFunc("/registration", common.RegistrationPage)
    router.HandleFunc("/dashboard", common.UserDashBoard)
    router.HandleFunc("/updateProfile", common.UpdateProfile)
    router.HandleFunc("/shiftUsers", common.ShiftUsers)
    router.HandleFunc("/viewManagerAndUsers", common.ViewManagerAndUsers)
    router.HandleFunc("/viewManagers", common.ViewManagers)
    router.HandleFunc("/viewUsers", common.ViewUsers)
    router.HandleFunc("/assignManagers", common.AssignManagers)
    router.HandleFunc("/logout", common.LogOutPage)
    router.HandleFunc("/clear", common.ClearSessionHandler)
    router.HandleFunc("/listUsersUnderHim", common.LitUsersUnderHim)
    router.HandleFunc("/viewListOtherManagers", common.ViewListOtherManagers)
    router.HandleFunc("/viewDeleteUsers", common.ViewDeleteUserUnderHim)
    router.HandleFunc("/sendMessage", common.SendMessage)
    router.HandleFunc("/roleChange", common.RoleChange)
    router.HandleFunc("/viewMessages", common.ViewMessages)
    router.HandleFunc("/readMessage", common.ReadMessage)
    http.ListenAndServe(":8080", router)
}



