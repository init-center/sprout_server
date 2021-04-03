package friend

import (
	"database/sql"
	"sprout_server/common/response/code"
	"sprout_server/dao/mysql"
	"sprout_server/models"
	"sprout_server/models/queryfields"

	"go.uber.org/zap"
)

func Create(p *models.ParamsAddFriend) int {
	// check the category exist
	exist, err := mysql.CheckFriendExistByName(p.Name)
	if err != nil {
		zap.L().Error("check friend exist by name failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if exist {
		return code.CodeFriendExist
	}

	//friend does not exist, can be created
	if err := mysql.CreateFriend(p); err != nil {
		zap.L().Error("create friend failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeCreated

}

func Update(p *models.ParamsAddFriend, u *models.UriUpdateFriend) int {
	// check the friend exist

	exist, err := mysql.CheckFriendExistById(u.Id)
	if err != nil {
		zap.L().Error("check friend exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeFriendNotExist
	}

	if err := mysql.UpdateFriend(p, u); err != nil {
		zap.L().Error("update friend failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func Delete(u *models.UriDeleteFriend) int {

	exist, err := mysql.CheckFriendExistById(u.Id)
	if err != nil {
		zap.L().Error("check friend exist by id failed", zap.Error(err))
		return code.CodeServerBusy
	}

	if !exist {
		return code.CodeCategoryNotExist
	}

	if err := mysql.DeleteFriend(u.Id); err != nil {
		zap.L().Error("delete category failed", zap.Error(err))
		return code.CodeServerBusy
	}

	return code.CodeOK

}

func GetByQuery(p *queryfields.FriendQueryFields) (models.FriendList, int) {
	friends, err := mysql.GetFriendList(p)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("get friend list failed", zap.Error(err))
		return friends, code.CodeServerBusy
	}

	return friends, code.CodeOK
}
