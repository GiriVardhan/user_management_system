package globalvar
 

type User struct {
    UserId          string `json:"user_id"`
    FirstName       string `json:"first_name"`
    LastName        string `json:"last_name"`
    Email           string `json:"email"`
    Password        string `json:"password"`
    Role            string `json:"role_name"`
    ManagerID       string `json:"manager_id"`
}

type Messages struct {
    Msg_Id            string `json:"msg_id"`
    Date_Created      string `json:"date_created"`
    Msg_From          string `json:"msg_from"`
    Msg_Header        string `json:"msg_header"`
    Msg_Text          string `json:"msg_text"`
    Msg_To            string `json:"mag_to"`
}
