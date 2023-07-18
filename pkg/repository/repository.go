package repository

import (
	"errors"
	"log"
	"service3/pkg/db"
	"service3/pkg/entity"
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func init() {
	DB, err = db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
}

func CreateAddress(address *entity.Address) error {
	return DB.Create(address).Error
}

func Create(userid int) (*entity.Cart, error) {
	cart := &entity.Cart{
		UserId: userid,
	}
	if err := DB.Create(cart).Error; err != nil {
		return nil, err
	}
	return cart, nil

}

func UpdateCart(cart *entity.Cart) error {
	return DB.Where("user_id = ?", cart.UserId).Save(&cart).Error
}

func GetByUserID(userid int) (*entity.Cart, error) {
	var cart entity.Cart
	result := DB.Where("user_id=?", userid).First(&cart)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart not found")
		}
		return nil, errors.New("cart not found")
	}
	return &cart, nil
}

func GetCartById(userId int) (*entity.Cart, error) {
	var cart entity.Cart
	result := DB.Where("user_id=?", userId).First(&cart)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &cart, nil
}

func CreateCartItem(cartItem *entity.CartItem) error {
	if err := DB.Create(cartItem).Error; err != nil {
		return err
	}
	return nil
}

func UpdateCartItem(cartItem *entity.CartItem) error {
	return DB.Save(cartItem).Error
}

func RemoveCartItem(cartItem *entity.CartItem) error {
	return DB.Where("product_name=?", cartItem.ProductName).Delete(&cartItem).Error
}

func GetByName(productName string, cartId int) (*entity.CartItem, error) {
	var cartItem entity.CartItem
	result := DB.Where("product_name = ? AND cart_id = ?", productName, cartId).First(&cartItem)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &cartItem, nil
}

func GetAllCartItems(cartId int) ([]entity.CartItem, error) {
	var cartItems []entity.CartItem
	result := DB.Where("cart_id=?", cartId).Find(&cartItems)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return cartItems, nil
}

func RemoveCartItems(cartId int) error {
	var cartItems entity.CartItem
	result := DB.Where("cart_id=?", cartId).Delete(&cartItems)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		return result.Error
	}
	return nil
}

func GetApparelByID(id int) (*entity.Apparel, error) {
	var apparel entity.Apparel
	result := DB.First(&apparel, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &apparel, nil
}

func AddApparelToWishlist(apparel *entity.Wishlist) error {
	if err := DB.Create(apparel).Error; err != nil {
		return err
	}
	return nil
}

func GetApparelFromWishlist(category string, id int, userId int) (bool, error) {
	var apparel entity.Wishlist
	result := DB.Where(&entity.Wishlist{UserId: userId, Category: category, ProductId: id}).First(&apparel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errors.New("Error finding apparel")
	}
	return true, nil
}

func GetWishlist(userId int) (*[]entity.Wishlist, error) {
	var wishlist []entity.Wishlist
	result := DB.Where("user_id=?", userId).Find(&wishlist)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &wishlist, nil
}

func RemoveFromWishlist(category string, id, userId int) error {
	product := entity.Wishlist{
		ProductId: id,
		UserId:    userId,
		Category:  category,
	}
	return DB.Where("user_id=?", userId).Delete(&product).Error
}

func GetCouponByCode(code string) (*entity.Coupon, error) {
	coupon := &entity.Coupon{}
	err := DB.Where("code = ?", code).First(coupon).Error
	if err != nil {
		return nil, err
	}
	return coupon, nil
}

func UpdateCouponCount(coupon *entity.Coupon) error {
	return DB.Save(coupon).Error
}

func UpdateCouponUsage(usedCoupon *entity.UsedCoupon) error {
	if err := DB.Create(usedCoupon).Error; err != nil {
		return err
	}
	return nil
}

func GetAllCoupons() ([]entity.Coupon, error) {
	var coupns []entity.Coupon
	currentTime := time.Now()
	err := DB.Where("valid_until >= ?", currentTime).Find(&coupns).Error
	if err != nil {
		return nil, err
	}
	return coupns, nil
}

func GetAllApparels(offset, limit int) ([]entity.Apparel, error) {
	var apparels []entity.Apparel
	err := DB.Offset(offset).Limit(limit).Where("removed = ?", false).Find(&apparels).Error
	if err != nil {
		return nil, err
	}
	return apparels, nil
}

func GetAllApparelsBySearch(offset, limit int, search string) ([]entity.Apparel, error) {
	var apparels []entity.Apparel
	err := DB.Where("name LIKE ?", search+"%").Offset(offset).Limit(limit).Find(&apparels).Error
	if err != nil {
		return nil, err
	}
	return apparels, nil
}

func GetAllApparelsByCategory(offset, limit int, category string) ([]entity.Apparel, error) {
	var apparels []entity.Apparel
	err := DB.Offset(offset).Limit(limit).Where("removed = ? AND sub_category = ?", false, category).Find(&apparels).Error
	if err != nil {
		return nil, err
	}
	return apparels, nil
}

func GetOfferByPrice(price int) ([]entity.Offer, error) {
	offers := []entity.Offer{}
	err := DB.Where("min_price <= ?", price).Find(&offers).Error
	if err != nil {
		return nil, err
	} else if offers == nil {
		return nil, err
	}
	return offers, nil
}

func GetByID(id int) (*entity.User, error) {
	var user entity.User
	result := DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func GetAddressByUserId(userid int) (*entity.Address, error) {
	var address entity.Address
	result := DB.Where(&entity.Address{UserId: userid}).Find(&address)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &address, nil
}
