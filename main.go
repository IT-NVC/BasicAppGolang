package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "time"
  "log"
  "fmt"
  "errors"
  "strconv"
)

type User struct {
  ID  int  `gorm:"primaryKey"`
  Name string  `gorm:"name"`
  Username string `gorm:"username"`
  Password string `gorm:"password"`
  Role string `gorm:"role"`
  CreatedAt *time.Time `gorm:"createdAt"`
  UpdatedAt *time.Time `gorm:"updatedAt"`
}

type CreateUser struct {
  Name  string
  Username string
  Password string
}

type UpdateUser struct {
  ID int
  Name  string
  Username string
  Password string
}

func (CreateUser) TableName() string {
  return "users"
}
func (UpdateUser) TableName() string {
  return "users"
}



func connectToMySQL() (*gorm.DB, error) {
  dsn := "root:password@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  if err != nil {
      return nil, err
  }

  return db, nil
}


func getListUser(db *gorm.DB) (*[]User,error) {
   var users []User
    result := db.Find(&users)
    if result.Error != nil {
        return nil, result.Error
    }
    return &users, nil

}


//create
func createUser(user *CreateUser, db *gorm.DB) error{

  //check Username
  var userCondition = User{Username: user.Username}
  userExists := db.First(&userCondition)

  if userExists !=nil {
    return errors.New("User exists")
  }

  result := db.Create(&user)
  
  if result.Error !=nil {
    return result.Error
  }

  return nil
}


//update
func updateUser(idUser int,user *UpdateUser, db *gorm.DB) (string,error){
  var userCondition = User{ID: idUser};
  userExists := db.First(&userCondition);
  fmt.Println(userExists)

  if userExists.Error != nil {
    return "User invalid.",errors.New("User invalid.")
  }

  user.ID = idUser
  result := db.Save(&user)
  if result.Error != nil {
      return "Update Faild", result.Error
  }
  return "Update successfully", nil
}


func main() {
  db, err := connectToMySQL()
  if err != nil {
      log.Fatal(err)
  }

  r := gin.Default()

  //get list user
  r.GET("/getListUser", func(c *gin.Context) {
    
    listUser,err :=getListUser(db);
    if err != nil {
      log.Fatal(err)
    }

    c.JSON(http.StatusOK, gin.H{
      "data": &listUser,
    })
    return;
  })


  //create
  r.POST("/createUser", func(c *gin.Context) {
    var user CreateUser
    
    //check body
    err := c.ShouldBind(&user)
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{
        "error": err.Error(),
      })

      return
    }

    errorCreate := createUser(&user, db);
    
    if errorCreate != nil{
      fmt.Println(errorCreate)
      c.JSON(http.StatusBadRequest, gin.H{
        "error": errorCreate,
      })

      return
    }

    c.JSON(http.StatusOK, gin.H{
      "message": "create user successfully",
    })
  })



  //update
  r.PATCH("/updateUser/:id", func(c *gin.Context) {
    var user UpdateUser
    idParam := c.Param("id")
    fmt.Println(idParam)
    if idParam == ":id" {
      c.JSON(http.StatusBadRequest, gin.H{
        "error": "Id invalid.",
      })

      return
    }
    var id, _ = strconv.Atoi(idParam)


    //check body
    err := c.ShouldBind(&user)
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{
        "error": err.Error(),
      })

      return
    }

    message,errorUpdate := updateUser(id, &user, db);
    
    if errorUpdate != nil{
      fmt.Println(errorUpdate)
      c.JSON(http.StatusBadRequest, gin.H{
        "error": message,
      })

      return
    }

    c.JSON(http.StatusOK, gin.H{
      "message": message,
    })
  })
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}