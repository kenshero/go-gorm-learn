package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kenshero/conn_db/config"
	"github.com/kenshero/conn_db/model"
)

type UserModel struct {
	ID        int    `gorm:"primary_key";"AUTO_INCREMENT"`
	Name      string `gorm:"size:255"`
	Address   string `gorm:"type:varchar(100)”`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Customer struct {
	CustomerID   int `gorm:"primary_key"`
	CustomerName string
	Contacts     []Contact `gorm:"ForeignKey:CustId"` //you need to do like this
}

type Contact struct {
	ContactID   int `gorm:"primary_key"`
	CountryCode int
	MobileNo    uint
	CustId      int
}

type UserL struct {
	ID        int `gorm:"primary_key"`
	Uname     string
	Languages []Language `gorm:"many2many:user_languages";"ForeignKey:UserId"`
	//Based on this 3rd table user_languages will be created
}

type Language struct {
	ID   int `gorm:"primary_key"`
	Name string
}

type UserLanguages struct {
	UserLId    int
	LanguageId int
}

func main() {
	config := config.GetConfig()
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Name,
		config.DB.Password)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	defer db.Close()
	if err != nil {
		panic("failed to connect database")
	}
	log.Println("Connection Established")
	// initDb(db)
	// insertDB(db)
	// updateDB(db)
	// deleteRecorddB(db)
	// hasOne(db)
	// hasMany(db)
	manyTomany(db)
}

func initDb(db *gorm.DB) {
	db.Debug().DropTableIfExists(&UserModel{})
	//Drops table if already exists
	db.Debug().AutoMigrate(&UserModel{})
	//Auto create table based on Model
}

func insertDB(db *gorm.DB) {
	user := &UserModel{Name: "John", Address: "New York"}
	tx := db.Begin()
	err := tx.Create(user).Error
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
	//You can insert multiple records too
	users := []UserModel{
		UserModel{Name: "Ricky", Address: "Sydney"},
		UserModel{Name: "Adam", Address: "Brisbane"},
		UserModel{Name: "Justin", Address: "California"},
	}

	for _, user := range users {
		db.Create(&user)
	}

}

func updateDB(db *gorm.DB) {
	user := &UserModel{}
	db.Where("name = ?", "John").Find(&user)
	fmt.Println(user)
	user.Address = "Brisbane"
	db.Save(&user)

	// // Update with column names, not attribute names
	// db.Model(&user).Update("Name", "Jack")

	// db.Model(&user).Updates(
	// 	map[string]interface{}{
	// 		"Name":    "Amy",
	// 		"Address": "Boston",
	// 	})

	// // UpdateColumn()
	// db.Model(&user).UpdateColumn("Address", "Phoenix")
	// db.Model(&user).UpdateColumns(
	// 	map[string]interface{}{
	// 		"Name":    "Taylor",
	// 		"Address": "Houston",
	// 	})
	// // Using Find()
	// db.Find(&user).Update("Address", "San Diego")

	// // Batch Update
	// db.Table("user_models").Where("address = ?", "california").Update("name", "Walker")
}

func deleteRecorddB(db *gorm.DB) {
	db.Table("user_models").Where("address= ?", "California").Delete(&UserModel{})
	// db.Where("address=?", "Sydney").Delete(&UserModel{})
	// db.Model(&UserModel{}).Delete(&UserModel{})
}

func hasOne(db *gorm.DB) {
	db.DropTableIfExists(&model.Place{}, &model.Town{})

	db.AutoMigrate(&model.Place{}, &model.Town{})
	db.Model(&model.Place{}).AddForeignKey("town_id", "towns(id)", "CASCADE", "CASCADE")

	t1 := model.Town{
		Name: "Pune",
	}
	t2 := model.Town{
		Name: "Mumbai",
	}
	t3 := model.Town{
		Name: "Hyderabad",
	}

	p1 := model.Place{
		Name: "Katraj",
		Town: t1,
	}
	p2 := model.Place{
		Name: "Thane",
		Town: t2,
	}
	p3 := model.Place{
		Name: "Secundarabad",
		Town: t3,
	}

	db.Save(&p1) //Saving one to one relationship
	db.Save(&p2)
	db.Save(&p3)

	fmt.Println("t1==>", t1, "p1==>", p1)
	fmt.Println("t2==>", t2, "p2s==>", p2)
	fmt.Println("t2==>", t3, "p2s==>", p3)

	// //Delete
	db.Where("name=?", "Hyderabad").Delete(&model.Town{})

	// //Update
	db.Model(&model.Place{}).Where("id=?", 1).Update("name", "Shivaji Nagar")

	// //Select
	places := model.Place{}
	towns := model.Town{}
	fmt.Println("Before Association", places)
	db.Where("name=?", "Shivaji Nagar").Find(&places)
	fmt.Println("After Association", places)
	err := db.Model(&places).Association("town").Find(&places.Town).Error
	fmt.Println("After Association", towns, places)
	fmt.Println("After Association", towns, places, err)
}

func hasMany(db *gorm.DB) {
	db.DropTableIfExists(&Contact{}, &Customer{})
	db.AutoMigrate(&Customer{}, &Contact{})
	db.Model(&Contact{}).AddForeignKey("cust_id", "customers(customer_id)", "CASCADE", "CASCADE") // Foreign key need to define manually

	Custs1 := Customer{CustomerName: "John", Contacts: []Contact{
		{CountryCode: 91, MobileNo: 956112},
		{CountryCode: 91, MobileNo: 997555}}}

	Custs2 := Customer{CustomerName: "Martin", Contacts: []Contact{
		{CountryCode: 90, MobileNo: 808988},
		{CountryCode: 90, MobileNo: 909699}}}

	Custs3 := Customer{CustomerName: "Raym", Contacts: []Contact{
		{CountryCode: 75, MobileNo: 798088},
		{CountryCode: 75, MobileNo: 965755}}}

	Custs4 := Customer{CustomerName: "Stoke", Contacts: []Contact{
		{CountryCode: 80, MobileNo: 805510},
		{CountryCode: 80, MobileNo: 758863}}}

	db.Create(&Custs1)
	db.Create(&Custs2)
	db.Create(&Custs3)
	db.Create(&Custs4)

	customers := &Customer{}
	contacts := &Contact{}

	db.Debug().Where("customer_name=?", "Martin").Preload("Contacts").Find(&customers) //db.Debug().Where("customer_name=?","John").Preload("Contacts").Find(&customers)
	fmt.Println("Customers", customers)
	fmt.Println("Contacts", contacts)

	// //Update
	db.Debug().Model(&Contact{}).Where("cust_id=?", 3).Update("country_code", 77)
	// //Delete
	// db.Debug().Where("customer_name=?", customers.CustomerName).Delete(&customers)
	// fmt.Println("After Delete", customers)
}

func manyTomany(db *gorm.DB) {
	db.DropTableIfExists(&UserLanguages{}, &Language{}, &UserL{})
	db.AutoMigrate(&UserL{}, &Language{}, &UserLanguages{})

	//All foreign keys need to define here
	db.Model(UserLanguages{}).AddForeignKey("user_l_id", "user_ls(id)", "CASCADE", "CASCADE")
	db.Model(UserLanguages{}).AddForeignKey("language_id", "languages(id)", "CASCADE", "CASCADE")

	langs := []Language{{Name: "English"}, {Name: "French"}}
	//log.Println(langs)

	user1 := UserL{Uname: "John", Languages: langs}
	user2 := UserL{Uname: "Martin", Languages: langs}
	user3 := UserL{Uname: "Ray", Languages: langs}
	db.Save(&user1) //save is happening
	db.Save(&user2)
	db.Save(&user3)

	fmt.Println("After Saving Records")
	fmt.Println("User1", &user1)
	fmt.Println("User2", &user2)
	fmt.Println("User3", &user3)

	// //Fetching
	user := &UserL{}
	db.Debug().Where("uname=?", "Ray").Find(&user)
	err := db.Debug().Model(&user).Association("Languages").Find(&user.Languages).Error
	fmt.Println("User is now coming", user, err)

	// //Deletion
	fmt.Println(user, "to delete")
	db.Debug().Where("uname=?", "John").Delete(&user)

	// //Updation
	db.Debug().Model(&UserL{}).Where("uname=?", "Ray").Update("uname", "Martin")

}
