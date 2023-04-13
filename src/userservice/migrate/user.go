package migrate

import (
	"gorm.io/gorm"

	"userservice/config"
)

// MigrateUserTable Migrate user table
func MigrateUserTable(db *gorm.DB, conf *config.Config) {

	m := db.Migrator()
	if !m.HasTable(&User{}) {
		if err := m.CreateTable(&User{}); err != nil {
			panic("migrate Failed.[ERROR]=>create user table failed.")
		}
		db.Exec("ALTER TABLE `user` COMMENT 'user table'")
	}

}

// User User Table
type User struct {
	// primary id
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:int(10);unique;autoIncrement;comment:primary id"`
	// username
	Username string `json:"username" gorm:"column:username;uniqueIndex:idx_username;type:varchar(20);default:'';not null;comment:username"`
	// password
	Password string `json:"password" gorm:"column:password;type:varchar(255);default:'';not null;comment:password"`
	// sex
	Sex int64 `json:"sex" gorm:"column:sex;type:tinyint(2);default:0;comment:sex[1=male2=female]"`
	// id number
	IDNumber string `json:"id_number" gorm:"column:id_number;type:varchar(30);default:'';comment:id number"`
	// email
	Email string `json:"email" gorm:"column:email;type:varchar(255);default:'';comment:email"`
	// phone
	Phone string `json:"phone" gorm:"column:phone;type:varchar(20);default:'';comment:phone"`
	// is_disable
	IsDisable int64 `json:"is_disable" gorm:"column:is_disable;type:tinyint(1);default:2;not null;comment:is_disable[1=enable2=disable]"`
	// access_token
	AccessToken string `json:"access_token" gorm:"column:access_token;type:varchar(255);default:'';comment:access_token"`
	// access_token_expire_time
	AccessTokenExpireTime int64 `json:"access_token_expire_time" gorm:"column:access_token_expire_time;type:int(10);default:1;not null;comment:access_token_expire_time"`
	// nick_name
	NickName string `json:"nick_name" gorm:"column:nick_name;type:varchar(20);default:'';comment:nick_name"`
	// real_name
	RealName string `json:"real_name" gorm:"column:real_name;type:varchar(10);default:'';comment:real_name"`
	//create_time / update_time
	CreateTime int64 `json:"create_time" gorm:"column:create_time;type:int(10);default:0;comment:create time'"`
	UpdateTime int64 `json:"update_time" gorm:"column:update_time;type:int(10);default:0;comment:update time"`
}

// TableName Table Name
func (*User) TableName() string {
	return "user"
}
