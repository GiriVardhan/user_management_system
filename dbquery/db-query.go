package dbquery

import (
        "fmt"
        "time"
        Cassandra "../Cassandra"
        helpers "../helpers"
	"github.com/gocql/gocql"

)


type UserInfo struct {
    UserId          string `json:"user_id"`
    FirstName       string `json:"first_name"`
    LastName        string `json:"last_name"`
    Email           string `json:"email"`
    Password        string `json:"password"`
    Role            string `json:"role_name"`
    ManagerID       string `json:"manager_id"`
}


func GetUserByEmail(email string) (user UserInfo) {  
     m := map[string]interface{}{}
     iter := Cassandra.Session.Query("SELECT user_id, first_name, last_name,email_id,password,role_name,manager_id FROM user_role_details WHERE email_id = ? ALLOW FILTERING", email).Iter() 
     for iter.MapScan(m) {
         user.UserId = m["user_id"].(string)
         user.FirstName = m["first_name"].(string)
         user.LastName = m["last_name"].(string)
         user.Email = m["email_id"].(string)
         user.Password = m["password"].(string)
         user.Role = m["role_name"].(string)
         user.ManagerID = m["manager_id"].(string)    
     }
     return user  
  
}

func RegisterUser(reg_user helpers.User) bool {
     var flag bool
     dt := time.Now() 

     if err := Cassandra.Session.Query("INSERT INTO user_role_details(user_id, role_name, first_name, last_name, email_id, password, date_created, date_modified, manager_id) VALUES(?, ?,?, ?, ?,?,?,?, ?)",
            reg_user.UserId, reg_user.Role,reg_user.FirstName, reg_user.LastName, reg_user.Email, reg_user.Password, dt, dt, reg_user.ManagerID).Exec(); err != nil         {
                fmt.Println("Error while inserting Emp")
                fmt.Println(err)
            } else {        
                flag = true
                //w.Write([]byte("<script>alert('User Registered Successfully');window.location = '/login'</script>"))                                         
     }
     return flag
}

//****************** Begin Check Duplicate Email Code ******************************
func CheckDuplicateEmail(email string) (message string) {
            
        fmt.Println(" **** get count ****")
        var count int 
        iter := Cassandra.Session.Query("SELECT count(*) FROM user_role_details where email_id = ? allow filtering", email ).Iter();
        for iter.Scan(&count) {
        }

        if count > 0 {
            message = "Email already exists"
        }

    return message   
}
//****************** End Check Duplicate Email Code ******************************

//********Begin********** Check User Exist With User ID ********************************
func CheckUserID(userid string) bool {
     var exists bool
     var count int
     iter := Cassandra.Session.Query("SELECT count(*) FROM user_role_details WHERE user_id = ? allow filtering", userid).Iter();
     for iter.Scan(&count) {
     }
     if count > 0 {
            exists = true
     }
    return exists
}
//********End********** Check User Exist With User ID ********************************

//*********Begin******** Get  Manager List *************************************

func GetManagerList() (user []helpers.User) {
     var managerList []helpers.User
     m := map[string]interface{}{}

     iter := Cassandra.Session.Query("SELECT first_name,last_name, user_id FROM user_role_details WHERE role_name = 'Manager' ALLOW FILTERING").Iter() 
     for iter.MapScan(m) {
		managerList = append(managerList, helpers.User{
		 FirstName: m["first_name"].(string),
                 LastName: m["last_name"].(string),
                 UserId: m["user_id"].(string),
        })
	m = map[string]interface{}{}
     }
     return managerList
   
}
//*********End******** Get  Manager List *************************************

//*********Begin******** Get User List By Manager *************************************

func GetUserByMngrList(managerid string) (user []helpers.User) {
     var userList []helpers.User
     m := map[string]interface{}{}

     iter := Cassandra.Session.Query("SELECT first_name,last_name, user_id, email_id FROM user_role_details WHERE manager_id = ? ALLOW FILTERING", managerid).Iter() 
        for iter.MapScan(m) {
		    userList = append(userList, helpers.User{
                FirstName: m["first_name"].(string),
                LastName: m["last_name"].(string),
                UserId: m["user_id"].(string),
                Email: m["email_id"].(string),
            })
            m = map[string]interface{}{}
        }
        return userList
   
}
//*********End******** Get User List By Manager *************************************

//*********Begin******** Get User List By Role *************************************

func GetUserByRole(bldqry, role_id string) (user []helpers.User) {
     var userList []helpers.User
     m := map[string]interface{}{}
     qry:= "SELECT first_name,last_name, user_id, role_name,email_id,manager_id FROM user_role_details WHERE role_name = " + role_id 
     if bldqry == ""{
         qry = qry + " ALLOW FILTERING " 
              
     }
     if bldqry == "'unassigned'"{
            qry = qry + "AND manager_id = "+ bldqry + " ALLOW FILTERING" 
           
     }

     iter := Cassandra.Session.Query(qry).Iter()
     for iter.MapScan(m) {
		userList = append(userList, helpers.User{
			FirstName: m["first_name"].(string),
            LastName: m["last_name"].(string),
            UserId: m["user_id"].(string),
            Email: m["email_id"].(string),
            ManagerID: m["manager_id"].(string),
        })
	m = map[string]interface{}{}
    }
    return userList
   
}
//*********End******** Get User List By Role *************************************

//*********Begin******** Update Profile*************************************

func UpdUserProfile(upd_col,upd_id,upd_val, user_id string) bool {
     var flag bool
     var qstring string
     qstring = "UPDATE user_role_details SET " +upd_id + " = ? WHERE  user_id = ?"
     err := Cassandra.Session.Query(qstring, upd_val,user_id).Exec(); 
     if err != nil {
            fmt.Println("Error while updating", upd_col)
            fmt.Println(err)
     }else{
          flag = true 
     }
     return flag   
}
//*********End******** Update Profile *************************************

//*********Begin******** View/ Delete Managers/Users*************************************

func DeleteManagerUser(role ,id string) bool {
    var flag bool  
    var listLen int
    var userList []helpers.User

    err := Cassandra.Session.Query("DELETE FROM user_role_details WHERE user_id = ?", id).Exec(); 
     if err != nil {
            fmt.Println("Error while deleting "+role)
            fmt.Println(err)
     }else{
          flag = true 
     } 

     if flag == true && role == "Manager" {
         userList = GetUserByMngrList(id)
         listLen = len(userList)
        if listLen > 0 {
           for i := 0; i < listLen ; i++ {
              err = Cassandra.Session.Query("UPDATE user_role_details SET manager_id = ? WHERE user_id = ?", "unassigned",userList[i].UserId).Exec(); 
             if err != nil {
                fmt.Println("Manager deleted but user under him still remains "+userList[i].UserId)
                fmt.Println(err)
             }
           } 

        }
     }
     return flag 
}
//*********End******** View/ Delete Managers/Users *************************************

//************** Begin ***** Create Message ************************

func CreateMessage(msg_header,msg_text,msg_to,msg_from string) bool {
    var flag bool
    var uid string  
    dt := time.Now()
    id, _ := gocql.RandomUUID()
    uid = id.String()
    fmt.Println(uid)
    err := Cassandra.Session.Query("INSERT INTO messages(msg_id, msg_header, msg_text, msg_from, msg_to, date_created) VALUES(?,?,?,?,?,?)",uid, msg_header,msg_text, msg_from,msg_to,dt).Exec(); 
     if err != nil {
            fmt.Println("Error while Creating Message")
            fmt.Println(err)
     }else {
           flag = true 
     }
    
     return flag 
}

//************** End ***** Create Message **************************


//*********Begin******** Role change User To Manager *************************************
func RoleChange(role_id string,user_id string) bool {
    var flag bool
    //fmt.Println(mngr_id)
    fmt.Println(user_id)
    err := Cassandra.Session.Query("UPDATE user_role_details SET role_name = ? , manager_id = 'unassigned' WHERE  user_id = ?",role_id,user_id).Exec(); 
    if err != nil {
           fmt.Println("Error while updating user's Manager")
           fmt.Println(err)
    }else{
         flag = true 
    } 
    return flag   
}
//*********End******** Role change User To Manager *************************************


//*********Begin******** Get User List By Role *************************************

func GetMsgList(msg_id string) (msglist []helpers.Messages) {
     var msgList []helpers.Messages
     m := map[string]interface{}{}
     qry := "SELECT msg_id,msg_header,msg_from, msg_text FROM messages "

     if msg_id != "" {
         qry = "SELECT msg_id,msg_header,msg_from, msg_text FROM messages WHERE msg_id = " +"'"+msg_id+"'"  
     }

     iter := Cassandra.Session.Query(qry).Iter()
     for iter.MapScan(m) {
		msgList = append(msgList, helpers.Messages{
                        Msg_Id:m["msg_id"].(string),
                        Msg_Header: m["msg_header"].(string),
            Msg_From: m["msg_from"].(string),
            Msg_Text: m["msg_text"].(string),
        })
	m = map[string]interface{}{}
    }
    return msgList
   
}
//*********End******** Get User List By Role *************************************

