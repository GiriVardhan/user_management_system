package handlers

import (
    "html/template"
    "net/http"
    "fmt"
    dbquery "../dbquery"
    helpers "../helpers"
    "golang.org/x/crypto/bcrypt"
    "github.com/gorilla/securecookie"
)

type userDetails struct {
    Userid   string  
    FirstName string 
    LastName  string 
    Emailid    string
    Password    string
    DateCreated  string
    DateModified string
}

/*type AllUsersResponse struct {
    Users []userDetails 
    ListLen int
    SuccessMessage string
    FailedMessage string
    IssueMsg string
}*/


type AllUsersResponse struct {
    Users []helpers.User
    ListLen int
    SuccessMessage string
    FailedMessage string
    IssueMsg string
    Managers []helpers.User
    IsShow bool
    UsersList []helpers.User
    SuccessUpdated string
    SelectMessage string
    ManagerID string
    MsgList []helpers.Messages
}


type userCredentials struct {
  EmailId   string
  Password  string
}

type GetPassword struct {
    Password    string 
}

type Response struct {
  WelcomeMessage  string
  ValidateMessage string 
  LogoutMessage string   
}

/*type User struct {
    UserId          string `json:"user_id"`
    FirstName       string `json:"first_name"`
    LastName        string `json:"last_name"`
    Email           string `json:"email"`
    Password        string `json:"password"`
    Role            string `json:"role_name"`
    ManagerID       string `json:"manager_id"`
}*/

type Allissues struct {
    IssueMsg string
    SuccessFlag bool
    EmailId string
    AdminRole bool
}

type Allinfo struct {
    IssueMsg string
    SuccessFlag bool
    EmailId string
    Role string
}

type UserIds  struct {
    ManagerId string
    UserId   string
}

var MessageID string

func HashPassword(password string) (string, error) {
        bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
        return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
        err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
        return err == nil
}

func validation(a string, b string, c string, d string) bool {
  if a != b || c != d{
        return false
  }
  return true
}

//********************* Begin Session Handling Code Block ************************
var cookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
    securecookie.GenerateRandomKey(32))
  
func getSession(request *http.Request) (userDetails helpers.User) {
    if cookie, err := request.Cookie("user-data"); err == nil {
        cookieValue := make(map[string]helpers.User)
        if err = cookieHandler.Decode("user-data", cookie.Value, &cookieValue); err == nil {
            userDetails = cookieValue["user-data"]
        }
    }
    return userDetails
}
  
func setSession(userDetails helpers.User, response http.ResponseWriter) {
    value := map[string]helpers.User{
        "user-data": userDetails,
    }
    if encoded, err := cookieHandler.Encode("user-data", value); err == nil {
        cookie := &http.Cookie{
            Name:  "user-data",
            Value: encoded,
            Path:  "/",
            MaxAge: 3600,
        }
        http.SetCookie(response, cookie)
    }
}
  
func clearSession(response http.ResponseWriter) {
    cookie := &http.Cookie{
        Name:   "user-data",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    }
    http.SetCookie(response, cookie)
}
  
  
func ClearSessionHandler(response http.ResponseWriter, request *http.Request) {
    clearSession(response)
    http.Redirect(response, request, "/", 302)
}
  

//****************** End Session Handling Code ***************************

//****************** Begin Welcome Page Code *****************************
func WelcomePage(w http.ResponseWriter, r *http.Request) {  
  tmpl, err := template.ParseFiles("templates/welcomePage.html")
  if err != nil {
      fmt.Println(err)
  }

  var welcomeHomePage string
  welcomeHomePage = "Login & Registration Forms"
  
  tmpl.Execute(w, Response{WelcomeMessage: welcomeHomePage})
}
//****************** End Welcome Page Code *****************************


//****************** Begin User Login Page *****************************
func LoginPage(w http.ResponseWriter, r *http.Request) {  
  tmpl, err := template.ParseFiles("templates/loginPage.html")
  if err != nil {
      fmt.Println(err)
  }
  var user helpers.User
  credentials := userCredentials{
     EmailId:   r.FormValue("emailId"),
    Password:   r.FormValue("password"), 
  }

  login_info := dbquery.GetUserByEmail(credentials.EmailId)
  user = helpers.User{
        UserId: login_info.UserId,
        FirstName: login_info.FirstName,
        LastName: login_info.LastName, 
        Role: login_info.Role,
        Email: login_info.Email,
  
    }

  var emailValidation string

  _userIsValid := CheckPasswordHash(credentials.Password, login_info.Password)

  if !validation(login_info.Email, credentials.EmailId, login_info.Password, credentials.Password) {
    emailValidation = "Please enter valid Email ID/Password"
  }

  if _userIsValid {
    setSession(user, w)
    http.Redirect(w, r, "/dashboard", http.StatusFound)
  }

  var welcomeLoginPage string
  welcomeLoginPage = "Login Page"

  tmpl.Execute(w, Response{WelcomeMessage: welcomeLoginPage, ValidateMessage: emailValidation})   
  
}
//***************************** End User Login Page Code *******************************************


//*********************** Begin Registration Page Code *********************************************
func RegistrationPage(w http.ResponseWriter, r *http.Request) {
    var flag bool
    var details helpers.User
    var targettmpl string
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    userDetails := getSession(r)
    
    if  (len(userDetails.FirstName) <= 0 && userDetails.Role == "Admin" ) {
       w.Write([]byte("<script>alert('Unauthorized Access!!,please login');window.location = '/login'</script>"))
    }

    if (userDetails.Role == "Admin" || userDetails.Role == "admin"){ 
        targettmpl = "templates/registrationPage.html"       
    }else{
        targettmpl = "templates/registrationUser.html"  
    }

    t := template.Must(template.ParseFiles(targettmpl))

    if r.Method != http.MethodPost {
        t.Execute(w, nil)
        return
    }
     

    details = helpers.User{
        UserId:     r.FormValue("userid"),
        FirstName:  r.FormValue("fname"),
        LastName:   r.FormValue("lname"),
        Email:      r.FormValue("email"),
        Password:   r.FormValue("pwd"),
        Role:   r.FormValue("role"),
        ManagerID : "unassigned",
    }
    
    if details.Role =="" {
       details.Role = "User" 
    }
     
    msg := dbquery.CheckDuplicateEmail(details.Email)

    details.Password, _ = HashPassword(details.Password)

    if msg == ""{
        fmt.Println(" **** Inserting a record ****")
        flag = dbquery.RegisterUser(details)
    }   

    t.Execute(w, Allinfo{EmailId: details.Email, IssueMsg: msg, SuccessFlag: flag} )
}
//*********************** End Registration Page Code ****************************************

//*********************** Start User Dashboard Code *****************************************
func UserDashBoard(w http.ResponseWriter, r *http.Request) {
     AuthorizePages(w,r) // Restrict Unauthorized User
     var targettmpl string
     w.Header().Set("Content-Type", "text/html; charset=utf-8")
     userDetails := getSession(r) 

     if (userDetails.Role == "Admin" || userDetails.Role == "admin"){ 
        targettmpl = "templates/admin-dashboard.html"       
     }else if (userDetails.Role == "Manager" || userDetails.Role == "manager"){
        targettmpl = "templates/manager-dashboard.html"  
     }else {
        targettmpl = "templates/user-dashboard.html"  
     }

     t, err := template.ParseFiles(targettmpl)
      if err != nil {
          fmt.Println(err)
      }
 
      items := struct {
          FirstName string
          LastName string
          UserID string
          Email string
          Homepage string
          
      }{
        FirstName: userDetails.FirstName,
        LastName: userDetails.LastName,
        UserID : userDetails.UserId,
        Email : userDetails.Email,
        Homepage: "Your Dashboard",
        
      }
      t.Execute(w, items)
  }
//*********************** End User Dashboard Code *****************************************



//****************** Begin Logout Page Code *****************************
func LogOutPage(w http.ResponseWriter, r *http.Request) {  
  tmpl, err := template.ParseFiles("templates/logoutPage.html")
  if err != nil {
      fmt.Println(err)
  }
  clearSession(w) 
  var logoutMsg string
  logoutMsg = "You are successfully Logged Out"
  
  tmpl.Execute(w, Response{LogoutMessage: logoutMsg})
}
//****************** End Logout Page Code *****************************

//****************** Begin Assign Managers to Users Code *****************************
func AssignManagers(w http.ResponseWriter, r *http.Request) {
    AuthorizePages(w,r) // Restrict Unauthorized User
    var selectMessage string
    var successUpdated string
    tmpl, err := template.ParseFiles("templates/assign-managers.html")
    if err != nil {
        fmt.Println(err)
    }
    //userDetails := getSession(r)
    var managerList []helpers.User
    var usersList []helpers.User
    fmt.Println("Getting all users and managers")
    
    managerList = dbquery.GetUserByRole("","'Manager'")
    usersList = dbquery.GetUserByRole("'unassigned'","'User'")

    
    if(len(usersList) == 0) {
        successUpdated = "Currently there are no unassigned Users"
    }
    if(len(managerList) == 0) {
        successUpdated = "Currently there are no managers to assign"
    }
   
    details:= helpers.User{
        UserId  :r.FormValue("users"),
        ManagerID : r.FormValue("managers"),
    }
   
    if ((details.ManagerID != "Select" && details.UserId != "Select") && (details.ManagerID != "" && details.UserId != ""))  {
          if (dbquery.UpdUserProfile("Manager ID","manager_id",details.ManagerID,details.UserId)){
              successUpdated = "Manager Assigned to User Successfully"
              //w.Write([]byte("<script>alert('Manager assigned to user Successfully')"))
          }
          
    }else { 
          selectMessage = "Please select Manager and User"
    }
    tmpl.Execute(w,AllUsersResponse{Managers: managerList,UsersList:usersList,SuccessUpdated:successUpdated,SelectMessage:selectMessage})
}
//****************** End Assign Managers to Users Code *****************************

//****************** Begin Authorize Page Code *****************************
func AuthorizePages(w http.ResponseWriter, r *http.Request) {
    userData := getSession(r)
    if  (len(userData.UserId) <= 0 || len(userData.FirstName) <= 0 ) {
        w.Write([]byte("<script>alert('Unauthorized Access!!,please login');window.location = '/login'</script>"))
   }
}
//****************** End Authorize Page Code *****************************


//****************** Begin Shift User Page Code *****************************
func ShiftUsers(w http.ResponseWriter, r *http.Request) {
    AuthorizePages(w,r) // Restrict Unauthorized User  
    tmpl, err := template.ParseFiles("templates/shiftUsers.html")
    if err != nil {
        fmt.Println(err)
    }
    
    var managerList []helpers.User
    var userList []helpers.User
    var tmpUserList []helpers.User
    var listLen int

    managerList = dbquery.GetUserByRole("","'Manager'")
    tmpUserList = dbquery.GetUserByRole("","'User'")
    listLen = len(tmpUserList);
    for i := 0;i < listLen; i++ {
        if tmpUserList[i].ManagerID != "unassigned"{
          userList = append(userList,helpers.User{
                FirstName: tmpUserList[i].FirstName,
                LastName: tmpUserList[i].LastName,
                UserId: tmpUserList[i].UserId,

         })  
           
        }
    }

    
    userId := UserIds{
        ManagerId: r.FormValue("managerId"),
        UserId: r.FormValue("userId"),
    }

    var successMessage string
    var isShow bool = false

    if ((userId.ManagerId != "Select" && userId.UserId != "Select") && (userId.ManagerId != "" && userId.UserId != ""))  {
        if (dbquery.UpdUserProfile("Manager to User","manager_id",userId.ManagerId,userId.UserId)){
            isShow = true
            successMessage = "User Shifted to New Manager Successfully"
        }
        listLen = len(userList);
    } 

    if (listLen == 0 && (userId.ManagerId == "Select" || userId.UserId == "Select" || userId.ManagerId == "" || userId.UserId == "")) {
        isShow = true
        successMessage = "Please select Manager & User"
    }

  
   tmpl.Execute(w, AllUsersResponse{Managers: managerList, Users: userList, SuccessMessage: successMessage, IsShow: isShow})

  }
//****************** End Shift User Page Code *****************************

//****************** Begin View Manager and Associated Users Page Code *****************************
func ViewManagerAndUsers(w http.ResponseWriter, r *http.Request) {
    AuthorizePages(w,r) // Restrict Unauthorized User   
    tmpl, err := template.ParseFiles("templates/viewManagersAndUsers.html")
    if err != nil {
        fmt.Println(err)
    }

    var managerList []helpers.User
    var userList []helpers.User

    managerList = dbquery.GetManagerList()

    userId := UserIds{
        ManagerId: r.FormValue("managerId"),
    }
    
    var isShow bool = false
    var noDataMessage string
    var listLen int
    
    if userId.ManagerId != "Select" && userId.ManagerId != "" {
         userList = dbquery.GetUserByMngrList(userId.ManagerId)
         listLen = len(userList);
    } else {
        isShow = true
        noDataMessage = "Please select Manager"
    }

    if (listLen == 0 && (userId.ManagerId != "Select" && userId.ManagerId != "")) {
        isShow = true
        noDataMessage = "There are no users for this Manager"
    }

    AuthorizePages(w,r) // Restrict Unauthorized User
    
    tmpl.Execute(w, AllUsersResponse{ListLen: listLen, Managers: managerList, Users: userList, IsShow: isShow, FailedMessage: noDataMessage})
}

//****************** End View Manager and Associated Users Page Code *****************************

// **************** Begin View/Delete Managers List *********************************
func ViewManagers(w http.ResponseWriter, r *http.Request) {  
    AuthorizePages(w,r)
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    t, err := template.ParseFiles("templates/viewManagers.html")
    
    if err != nil {
        fmt.Println(err) // Ugly debug output
        w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
        return
    }
    
    userId := UserIds{
        ManagerId: r.FormValue("managerId"),
    }

    var successMessage string
    var isShow bool 

    if (userId.ManagerId != "" ) {
        if (dbquery.DeleteManagerUser("Manager",userId.ManagerId)){
            isShow = true
            successMessage = "Manager Deleted Successfully"
        }
    }

    var managerList []helpers.User
    managerList = dbquery.GetUserByRole("","'Manager'")
    t.Execute(w, AllUsersResponse{Users: managerList, SuccessMessage: successMessage, IsShow: isShow})  
}
// **************** End View/Delete Managers List *********************************


// **************** Begin View/Delete User List *********************************

func ViewUsers(w http.ResponseWriter, r *http.Request) { 
    AuthorizePages(w,r) 
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    t, err := template.ParseFiles("templates/viewUsers.html")
    if err != nil {
        fmt.Println(err) // Ugly debug output
        w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
        return
    }
    
    userId := UserIds{
        UserId: r.FormValue("userId"),
    }

    var successMessage string
    var isShow bool 

    if (userId.UserId != "" ) {
        if (dbquery.DeleteManagerUser("User",userId.UserId)){
            isShow = true
            successMessage = "User Deleted Successfully"
        }
    }

    var userList []helpers.User     
    userList = dbquery.GetUserByRole("","'User'")
    t.Execute(w, AllUsersResponse{Users: userList, SuccessMessage: successMessage, IsShow: isShow})  
}
// **************** End View/Delete User List *********************************

//**************** Begin Update Profile *****************************************
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
    AuthorizePages(w,r) // Restrict Unauthorized User
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    t, err := template.ParseFiles("templates/update-profile.html")
    userData := getSession(r)
    if err != nil {
        fmt.Println(err) // Ugly debug output
        w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
        return
      }
    
      if r.Method != http.MethodPost {
        t.Execute(w, nil)
        return
    }
    
    var successMessage string
    var failedMessage string
    var msg string
    msg = ""
  
    details := userDetails{
        Userid: userData.UserId,
        FirstName: r.FormValue("fname"),
        LastName: r.FormValue("lname"),
        Emailid: r.FormValue("emailid"),
        Password: r.FormValue("pwd2"),
     }
    if (details.FirstName != "" || len(details.FirstName) > 0) {
       if (dbquery.UpdUserProfile("First Name","first_name",details.FirstName,details.Userid)){
             successMessage = "User First Name Updated Successfully"
             //w.Write([]byte("<script>alert('Email Id  Updated Successfully,please login');window.location = '/login'</script>"))
       }
    }
    if (details.LastName != "" || len(details.LastName) > 0){
       if (dbquery.UpdUserProfile("Last Name","last_name",details.LastName,details.Userid)){
             successMessage = "User Last Name Updated Successfully"
             //w.Write([]byte("<script>alert('Email Id  Updated Successfully,please login');window.location = '/login'</script>"))
       }
    }

    if (details.Emailid != "" || len(details.Emailid) > 0){
       msg = dbquery.CheckDuplicateEmail(details.Emailid)
       if(dbquery.CheckUserID(details.Userid)) {
         if msg == ""{
            if (dbquery.UpdUserProfile("Email ID","email_id",details.Emailid,details.Userid)){
               successMessage = "User Email Id Updated Successfully"
               w.Write([]byte("<script>alert('Email Id  Updated Successfully,please login');window.location = '/logout'</script>"))
            }
         }else {
               //failedMessage = "Email Already Exist"  
         }

    }else {
          failedMessage = "There is no User with that User Id"
    }
       
    }

    if (details.Password != "" || len(details.Password) > 0){
       password := details.Password
       hash, _ := HashPassword(password) 
       if dbquery.CheckUserID(details.Userid) {
          if dbquery.UpdUserProfile("Password","password",hash, details.Userid) {
             successMessage = "Password Updated Successfully"
             w.Write([]byte("<script>alert('Password Updated Successfully,please login');window.location = '/logout'</script>"))
          }
       }else {
        failedMessage = "There is no User with that User Id"
       }  
    }
    
     
    t.Execute(w, AllUsersResponse{SuccessMessage: successMessage, FailedMessage: failedMessage,IssueMsg: msg})  
}
//**************** End Update Profile *****************************************

// **************** Begin List users under Manager Page *********************************
func LitUsersUnderHim(w http.ResponseWriter, r *http.Request) {  
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    t, err := template.ParseFiles("templates/listUsersUnderHim.html")

    userDetails := getSession(r)

    AuthorizePages(w,r)
    if err != nil {
        fmt.Println(err) // Ugly debug output
        w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
        return
      }
    
      if err != nil {
          fmt.Println(err)
      }
      var userList []helpers.User
      var listLen int
      var failedMessage string
      var isShow bool = false

      userList = dbquery.GetUserByMngrList(userDetails.UserId)
      listLen = len(userList);

      if listLen == 0 {
        isShow = true
        failedMessage =  "Currently you are not assigned for any User"
      }     

    t.Execute(w, AllUsersResponse{Users: userList, ListLen: listLen, FailedMessage: failedMessage, IsShow: isShow})  
}
// **************** End List users under Manager Page *********************************

// **************** Begin List of Other Manager Page *********************************
func ViewListOtherManagers(w http.ResponseWriter, r *http.Request) {  
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    t, err := template.ParseFiles("templates/viewListOtherManagers.html")

    userDetails := getSession(r)

    AuthorizePages(w,r)
    if err != nil {
        fmt.Println(err) // Ugly debug output
        w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
        return
      }
    
      if err != nil {
          fmt.Println(err)
      }
      var managerList []helpers.User
      var listLen int
      var failedMessage string
      var isShow bool = false

      managerList = dbquery.GetManagerList()
      listLen = len(managerList);

      var managerList1 []helpers.User

      for i := 0; i < listLen; i++ {
        if managerList[i].UserId != userDetails.UserId {
            managerList1 = append(managerList1, helpers.User{
                FirstName: managerList[i].FirstName,
                LastName: managerList[i].LastName,
                UserId: managerList[i].UserId,
            })
        }
      }
      if listLen == 0 {
        isShow = true
        failedMessage =  "Currently you are not assigned for any User"
      }     

    t.Execute(w, AllUsersResponse{Users: managerList1, ListLen: listLen, FailedMessage: failedMessage, IsShow: isShow})  
}
// **************** End List of Other Manager Page *********************************

// **************** Begin List users under Manager Page *********************************
func ViewDeleteUserUnderHim(w http.ResponseWriter, r *http.Request) {  
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    t, err := template.ParseFiles("templates/viewUsers.html")

    userDetails := getSession(r)

    if err != nil {
        fmt.Println(err) // Ugly debug output
        w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
        return
    }
    
    if err != nil {
        fmt.Println(err)
    }

    userId := UserIds{
        UserId: r.FormValue("userId"),
    }

    var userList []helpers.User
    var successMessage string
    var isShow bool 

    if (userId.UserId != "" ) {
        if (dbquery.DeleteManagerUser("User",userId.UserId)){
            isShow = true
            successMessage = "User Deleted Successfully"
        }
    }
     
    userList = dbquery.GetUserByMngrList(userDetails.UserId)
    t.Execute(w, AllUsersResponse{Users: userList, SuccessMessage: successMessage, IsShow: isShow}) 
}
// **************** End List users under Manager Page *********************************

//********* Begin ******* Message Functionality *********************
func SendMessage(w http.ResponseWriter, r *http.Request) {  
    AuthorizePages(w,r)
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    t, err := template.ParseFiles("templates/createMessage.html")
    if err != nil {
        fmt.Println(err) // Ugly debug output
        w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
        return
    }
    var msg string
    var msgfrm string
    var successMessage string
    var selectMessage string
    userDetails := getSession(r)
    var managerList []helpers.User

    if userDetails.Role == "Admin" {
       msgfrm = userDetails.Role
    }else{
       msgfrm = userDetails.UserId 
    }
   
    managerList = dbquery.GetUserByRole("","'Manager'")
    details := struct {
          msg_header string
          msg_text string
          msg_to string
          msg_from string
                   
      }{
        msg_header: r.FormValue("msg_header"),
        msg_text: r.FormValue("msg_text"),
        msg_to : r.FormValue("managers"),
        msg_from : msgfrm,
        
      }

    if (details.msg_to != "Select" && details.msg_to != "" ) {
          if (dbquery.CreateMessage(details.msg_header, details.msg_text, details.msg_to, details.msg_from)){
              successMessage = "Your Message Sent Successfully!!"
          }else {
             msg = "Message Sent Failed!!"
          }
          
    }else { 
          selectMessage = "Please All Users OR Specific Manager"
    }
   t.Execute(w, AllUsersResponse{Managers: managerList,IssueMsg: msg,SuccessMessage: successMessage,SelectMessage:selectMessage})  
}

//********* End ******* Message Functionality *********************

func RoleChange(w http.ResponseWriter, r *http.Request) {
    AuthorizePages(w,r) // Restrict Unauthorized User
    var selectMessage string
    var successUpdated string
    tmpl, err := template.ParseFiles("templates/role-change.html")
    if err != nil {
        fmt.Println(err)
    }
    //userDetails := getSession(r)
    var usersList []helpers.User
    fmt.Println("Getting all users")
    
    usersList = dbquery.GetUserByRole("","'User'")

    if(len(usersList) == 0) {
        successUpdated = "Currently there are no Users"
    }
    details:= helpers.User{
        UserId  :r.FormValue("users"),
    }
    if details.UserId != "Select" && details.UserId != ""   {
          if (dbquery.RoleChange("Manager",details.UserId)){
              successUpdated = "Role changed to Manager Successfully"
          }
          
    }else {
              selectMessage = "Please select User"
    }
    tmpl.Execute(w,AllUsersResponse{UsersList:usersList,SuccessUpdated:successUpdated,SelectMessage:selectMessage})
}
// **************** End Role change *********************************


// **************** Begin Message List *********************************

func ViewMessages(w http.ResponseWriter, r *http.Request) { 
    AuthorizePages(w,r) 
    MessageID = ""
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    t, err := template.ParseFiles("templates/viewMessages.html")
    if err != nil {
        fmt.Println(err) // Ugly debug output
        w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
        return
    }
    
    msgId := helpers.Messages{
        Msg_Id: r.FormValue("messageId"),
    }

    if msgId.Msg_Id != ""{
       MessageID = msgId.Msg_Id
       http.Redirect(w, r, "/readMessage", http.StatusFound)
    }

    var msgList []helpers.Messages     
    msgList = dbquery.GetMsgList("")
    t.Execute(w, AllUsersResponse{MsgList:msgList})  
}
// **************** End Message List *********************************

// **************** Begin Read Message *********************************

func ReadMessage(w http.ResponseWriter, r *http.Request) { 
    AuthorizePages(w,r) 
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    t, err := template.ParseFiles("templates/readMessage.html")
    if err != nil {
        fmt.Println(err) // Ugly debug output
        w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
        return
    }
    
    var msgList []helpers.Messages     
    msgList = dbquery.GetMsgList(MessageID)
    t.Execute(w, AllUsersResponse{MsgList:msgList})  
}
// **************** End Read Message *********************************
