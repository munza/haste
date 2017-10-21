package authrepo

import (
	"errors"
	auth "haste/lib/auth"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type userRepo struct {
	BaseRepo
}

func UserRepo() *userRepo {
	user := &userRepo{}
	return user
}

func (*userRepo) Find(id int) (auth.User, error) {
	db := DB()
	defer db.Close()

	var user auth.User
	err := db.First(&user, "id = ?", id)
	if err.RecordNotFound() {
		return user, errors.New("Not Found!")
	}

	return user, nil
}
