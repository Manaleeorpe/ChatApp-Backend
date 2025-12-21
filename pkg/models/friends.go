package models

import (
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
)

type Friends struct {
	gorm.Model
	//ID uint `gorm:"primaryKey"` ID already exists in gorm.Model

	Friend1UserID uint
	Friend1       User `gorm:"foreignKey:Friend1UserID;constraint:OnDelete:CASCADE;"`
	Friend2UserID uint
	Friend2       User `gorm:"foreignKey:Friend2UserID;constraint:OnDelete:CASCADE;"`
	Status        string
}

type UserFriendsResponse struct {
	UserID       uint      `json:"ID"`
	Name         string    `json:"name"`
	Email_id     string    `json:"email_id"`
	Phone_number string    `json:"phone_number"`
	FriendsSince time.Time `json:"friends_since"`
}

func (f *Friends) CreateFriend() (*Friends, error) { //method is of struct type Friend
	if err := db.Create(&f).Error; err != nil {
		return nil, err
	}
	if err := db.Preload("Friend1").Preload("Friend2").First(&f, f.ID).Error; err != nil {
		return nil, err
	}
	return f, nil
}

func GetFriendRequestById(Id int64) (*Friends, *gorm.DB) {
	var getFriendRequest Friends
	db := db.Where("ID=?", Id).Find(&getFriendRequest)
	db.Preload("Friend1").Preload("Friend2").First(&getFriendRequest, getFriendRequest.ID)

	return &getFriendRequest, db
}

func GetAFriendRequestByOFriendIDs(Friend1ID int64, Friend2ID int64) (ID int64) {
	var getFriendRequest1 Friends
	var getFriendRequest2 Friends

	db1 := db.Where("friend1_user_id=? AND friend2_user_id=?", Friend1ID, Friend2ID).Find((&getFriendRequest1))
	db2 := db.Where("friend1_user_id=? AND friend2_user_id=?", Friend2ID, Friend1ID).Find((&getFriendRequest2))

	if db1.RecordNotFound() && db2.RecordNotFound() {
		return 0
	}

	if reflect.DeepEqual(getFriendRequest1, Friends{}) {
		return int64(getFriendRequest2.ID)
	} else if reflect.DeepEqual(getFriendRequest2, Friends{}) {
		return int64(getFriendRequest1.ID)
	} else {
		return 0
	}

}

func GetAllUsersFriendsModel(ID int64, status string) ([]UserFriendsResponse, *gorm.DB) {
	var AllFriends []Friends

	dbResult := db.Preload("Friend1").Preload("Friend2").
		Where("(friend1_user_id = ? OR friend2_user_id = ?) AND status = ?", ID, ID, status).
		Find(&AllFriends)

	//decide which to select friend1 or friend2 and get the required response
	var response []UserFriendsResponse

	for _, f := range AllFriends {
		var Friend User

		if int64(f.Friend1UserID) == ID {
			Friend = f.Friend2
		} else {
			Friend = f.Friend1
		}

		r := UserFriendsResponse{
			UserID:       Friend.ID,
			Name:         Friend.Name,
			Email_id:     Friend.Email_id,
			Phone_number: Friend.Phone_number,
			FriendsSince: Friend.CreatedAt,
		}
		response = append(response, r)
	}

	return response, dbResult

}

func DeleteFriendRequestModel(ID int64) *gorm.DB {
	var FriendRequest Friends
	result := db.Where("ID=?", ID).Unscoped().Delete(&FriendRequest)

	return result

}
func FriendRequestUserCanAcceptModel(ID int64) ([]UserFriendsResponse, *gorm.DB) {
	var AllFriends []Friends

	dbResult := db.Preload("Friend1").Preload("Friend2").
		Where("friend2_user_id = ? AND status = 'Pending'", ID).
		Find(&AllFriends)

	//decide which to select friend1 or friend2 and get the required response
	var response []UserFriendsResponse

	for _, f := range AllFriends {

		r := UserFriendsResponse{
			UserID:       f.Friend1.ID,
			Name:         f.Friend1.Name,
			Email_id:     f.Friend1.Email_id,
			Phone_number: f.Friend1.Phone_number,
			FriendsSince: f.Friend1.CreatedAt,
		}
		response = append(response, r)
	}

	return response, dbResult

}
