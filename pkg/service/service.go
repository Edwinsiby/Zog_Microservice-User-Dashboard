package service

import (
	"context"
	"errors"
	"log"
	"service3/pb"
	"strconv"

	"service3/pkg/entity"
	repo "service3/pkg/repository"
)

type UserDashboard struct {
	pb.UnimplementedUserDashboardServer
}

func (s *UserDashboard) MyMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Println("Microservice1: MyMethod called")

	result := "Hello, " + req.Data
	return &pb.Response{Result: result}, nil
}

func (s *UserDashboard) AddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.AddAddressResponse, error) {
	address := &entity.Address{
		UserId:  int(req.Userid),
		House:   req.House,
		City:    req.City,
		Street:  req.Street,
		Pincode: int(req.Pincode),
		Type:    req.Type,
	}
	err := repo.CreateAddress(address)
	if err != nil {
		return nil, err
	} else {
		return &pb.AddAddressResponse{Result: "Address added succesfuly"}, nil
	}
}

func (s *UserDashboard) AddToCart(ctx context.Context, req *pb.AddToCartRequest) (*pb.AddToCartResponse, error) {
	var userCart *entity.Cart
	var cartId int
	userCart, err := repo.GetByUserID(int(req.Userid))
	if err != nil {
		cart, err1 := repo.Create(int(req.Userid))
		if err1 != nil {
			return nil, errors.New("Failed to create user cart")
		}
		userCart = cart
		cartId = int(cart.ID)
	} else {
		cartId = int(userCart.ID)
	}
	apparel, err := repo.GetApparelByID(int(req.Productid))
	if err != nil {
		return nil, errors.New("Apparel not found")
	}
	cartItem := &entity.CartItem{
		CartId:      cartId,
		ProductId:   int(req.Productid),
		Category:    "apparel",
		Quantity:    int(req.Quantity),
		ProductName: apparel.Name,
		Price:       float64(apparel.Price),
	}
	existingApparel, err := repo.GetByName(apparel.Name, cartId)
	if existingApparel == nil {
		err = repo.CreateCartItem(cartItem)
		if err != nil {
			return nil, errors.New("Adding new ticket to cart item failed")
		}
	} else {
		existingApparel.Quantity += int(req.Quantity)
		err := repo.UpdateCartItem(existingApparel)
		if err != nil {
			return nil, errors.New("error updating existing cartitem")
		}
	}
	userCart.TotalPrice += cartItem.Price * float64(req.Quantity)
	userCart.ApparelQuantity += int(req.Quantity)
	err1 := repo.UpdateCart(userCart)
	if err1 != nil {
		return nil, errors.New("Cart price updation failed")
	}
	return &pb.AddToCartResponse{Result: "Product added to cart succesfuly"}, nil
}

func (s *UserDashboard) AddToWishList(ctx context.Context, req *pb.AddToWishListRequest) (*pb.AddToWishListResponse, error) {
	apparel, err := repo.GetApparelByID(int(req.Productid))
	if err != nil {
		return nil, errors.New("Apparel not found")
	}
	exsisting, err := repo.GetApparelFromWishlist(apparel.Category, apparel.ID, int(req.Userid))
	if err != nil {
		return nil, errors.New("Error finding exsisting product")
	}
	if exsisting == true {
		return nil, errors.New("Product already exsist in wishlist")
	} else {
		wishApparel := &entity.Wishlist{
			UserId:      int(req.Userid),
			Category:    "apparel",
			ProductId:   apparel.ID,
			ProductName: apparel.Name,
			Price:       float64(apparel.Price),
		}
		err = repo.AddApparelToWishlist(wishApparel)
		if err != nil {
			return nil, errors.New("Product adding to wishlist failed")
		}
	}
	return &pb.AddToWishListResponse{Result: "Added to wishlist succesfuly"}, nil
}

func (s *UserDashboard) ApparelDetails(ctx context.Context, req *pb.ApparelDetailsRequest) (*pb.ApparelDetailsResponse, error) {
	apparel, err := repo.GetApparelByID(int(req.Id))
	if err != nil {
		return nil, err
	}
	resp := &pb.ApparelDetailsResponse{
		Apparel: &pb.Apparel{
			Id:          int32(apparel.ID),
			Name:        apparel.Name,
			Price:       int32(apparel.Price),
			Image:       apparel.ImageURL,
			Subcategory: apparel.SubCategory,
		},
	}
	return resp, nil
}
func (s *UserDashboard) Apparels(ctx context.Context, req *pb.ApparelsRequest) (*pb.ApparelsResponse, error) {
	offset := (req.Page - 1) * req.Limit
	var apparelList []entity.Apparel
	var err error
	if req.Category == "" {
		apparelList, err = repo.GetAllApparels(int(offset), int(req.Limit))
		if err != nil {
			return nil, err
		}
	} else {
		apparelList, err = repo.GetAllApparelsByCategory(int(offset), int(req.Limit), req.Category)
		if err != nil {
			return nil, err
		}
	}
	var pbApparels []*pb.Apparel
	for _, apparel := range apparelList {
		pbApparel := &pb.Apparel{
			Id:          int32(apparel.ID),
			Name:        apparel.Name,
			Price:       int32(apparel.Price),
			Image:       apparel.ImageURL,
			Subcategory: apparel.SubCategory,
		}
		pbApparels = append(pbApparels, pbApparel)
	}

	response := &pb.ApparelsResponse{
		Apparels: pbApparels,
	}
	return response, nil
}
func (s *UserDashboard) ApplyCoupon(ctx context.Context, req *pb.ApplyCouponRequest) (*pb.ApplyCouponResponse, error) {
	var totalOffer, totalPrice int
	userCart, err := repo.GetByUserID(int(req.Userid))
	if err != nil {
		return nil, errors.New("Failed to find user cart")
	}
	coupon, err := repo.GetCouponByCode(req.Code)
	if err != nil {
		return nil, errors.New("Sorry coupon not found")
	}
	cartItems, err := repo.GetAllCartItems(int(userCart.ID))
	if err != nil {
		return nil, errors.New("User Cart Items not found")
	}
	for _, cartItem := range cartItems {
		if cartItem.Category == coupon.Category {
			totalPrice += int(cartItem.Price) * cartItem.Quantity
		}
	}
	if totalPrice > 0 {
		if coupon.Type == "percentage" {
			totalOffer = totalPrice / coupon.Amount
		} else {
			totalOffer = coupon.Amount
		}
	} else {
		return nil, errors.New("Add more product from different category")
	}
	if userCart.OfferPrice != 0 {
		return nil, errors.New("User Cart offer already applied")
	} else {
		userCart.OfferPrice = totalOffer
		err = repo.UpdateCart(userCart)
		if err != nil {
			return nil, errors.New("User Cart updation failed")
		}
		var usedCoupon = entity.UsedCoupon{
			UserId:     int(req.Userid),
			CouponCode: req.Code,
		}
		err = repo.UpdateCouponUsage(&usedCoupon)
		if err != nil {
			return nil, errors.New("User Coupon Usage updation failed")
		}
		var coupons = entity.Coupon{
			UsedCount: coupon.UsedCount + 1,
		}
		err = repo.UpdateCouponCount(&coupons)
		if err != nil {
			return nil, errors.New("User Coupon Usage updation failed")
		}
	}
	return &pb.ApplyCouponResponse{Result: strconv.Itoa(totalOffer)}, nil
}
func (s *UserDashboard) AvailableCoupons(ctx context.Context, req *pb.AvailableCouponsRequest) (*pb.AvailableCouponsResponse, error) {
	coupons, err := repo.GetAllCoupons()
	if err != nil {
		return nil, errors.New(err.Error())
	}
	var couponList []*pb.Coupon
	for _, u := range coupons {
		pbCoupon := &pb.Coupon{
			Code:     u.Code,
			Amount:   int32(u.Amount),
			Type:     u.Type,
			Category: u.Category,
		}
		couponList = append(couponList, pbCoupon)
	}
	response := &pb.AvailableCouponsResponse{
		Coupons: couponList,
	}
	return response, nil
}
func (s *UserDashboard) Cart(ctx context.Context, req *pb.CartRequest) (*pb.CartResponse, error) {
	userCart, err := repo.GetByUserID(int(req.Userid))
	if err != nil {
		return nil, errors.New("Failed to find user cart")
	} else {
		resp := &pb.CartResponse{
			Userid:     int32(userCart.UserId),
			Quantity:   int32(userCart.ApparelQuantity),
			Totalprice: int32(userCart.TotalPrice),
			Offerprice: int32(userCart.OfferPrice),
		}

		return resp, nil
	}
}
func (s *UserDashboard) CartList(ctx context.Context, req *pb.CartListRequest) (*pb.CartListResponse, error) {
	userCart, err := repo.GetByUserID(int(req.Userid))
	if err != nil {
		return nil, errors.New("Failed to find user cart")
	}
	cartItems, err := repo.GetAllCartItems(int(userCart.ID))
	if err != nil {
		return nil, err
	}
	var pbApparels []*pb.Apparel
	for _, apparel := range cartItems {
		pbApparel := &pb.Apparel{
			Id:          int32(apparel.ID),
			Name:        apparel.ProductName,
			Price:       int32(apparel.Price),
			Subcategory: apparel.Category,
		}
		pbApparels = append(pbApparels, pbApparel)
	}

	response := &pb.CartListResponse{
		Appaers: pbApparels,
	}
	return response, nil
}

func (s *UserDashboard) OfferCheck(ctx context.Context, req *pb.OfferCheckRequest) (*pb.OfferCheckResponse, error) {
	userCart, err := repo.GetByUserID(int(req.Userid))
	if err != nil {
		return nil, errors.New("Failed to find user cart")
	}
	offers, err := repo.GetOfferByPrice(int(userCart.TotalPrice))
	if err != nil {
		return nil, errors.New("No valid offers, Add few more products worth of 500")
	}
	var pbOffers []*pb.Coupon
	for _, offer := range offers {
		pbOffer := &pb.Coupon{
			Code:     offer.Name,
			Type:     offer.Type,
			Category: offer.Category,
			Amount:   int32(offer.Amount),
		}
		pbOffers = append(pbOffers, pbOffer)
	}
	response := &pb.OfferCheckResponse{
		Offers: pbOffers,
	}
	return response, nil
}
func (s *UserDashboard) RemoveFromCart(ctx context.Context, req *pb.RemoveFromCartRequest) (*pb.RemoveFromCartResponse, error) {
	apparel, err := repo.GetApparelByID(int(req.Productid))
	if err != nil {
		return nil, errors.New("Apparel not found")
	}
	userCart, err := repo.GetByUserID(int(req.Userid))
	if err != nil {
		return nil, errors.New("Failed to find user cart")
	}
	existingApparel, err1 := repo.GetByName(apparel.Name, int(userCart.ID))
	if err1 != nil {
		return nil, errors.New("Removing apparel from cart failed")
	}
	if existingApparel.Quantity == 1 {
		err := repo.RemoveCartItem(existingApparel)
		if err != nil {
			return nil, errors.New("Removin apparel from cart failed")
		}
	} else {
		existingApparel.Quantity -= 1
		err := repo.UpdateCartItem(existingApparel)
		if err != nil {
			return nil, errors.New("error updating existing cartitem")
		}
	}
	userCart.TotalPrice -= float64(apparel.Price)
	userCart.ApparelQuantity -= 1

	if userCart.OfferPrice > 0 {
		userCart.OfferPrice = 0
	}
	err = repo.UpdateCart(userCart)
	if err != nil {
		return nil, errors.New("Remove from cart failed")
	}

	return &pb.RemoveFromCartResponse{Result: "Product Removed from cart"}, nil
}
func (s *UserDashboard) RemoveFromWishlist(ctx context.Context, req *pb.RemoveFromWishlistRequest) (*pb.RemoveFromWishlistResponse, error) {
	exsisting, err := repo.GetApparelFromWishlist("apparel", int(req.Productid), int(req.Userid))
	if err != nil {
		return nil, errors.New("Error finding exsisting product")
	}
	if exsisting == true {
		err = repo.RemoveFromWishlist("apparel", int(req.Productid), int(req.Userid))
		if err != nil {
			return nil, errors.New("Can't remove from wishlist")
		}
	} else {
		return nil, errors.New("Apparel not found")
	}
	return &pb.RemoveFromWishlistResponse{Result: "Removed from wishlist succesfuly"}, nil
}
func (s *UserDashboard) SearchApparels(ctx context.Context, req *pb.SearchApparelsRequest) (*pb.SearchApparelsResponse, error) {
	offset := (req.Page - 1) * req.Limit
	apparelList, err := repo.GetAllApparelsBySearch(int(offset), int(req.Limit), req.Search)
	if err != nil {
		return nil, err
	} else {
		var pbApparels []*pb.Apparel
		for _, apparel := range apparelList {
			pbApparel := &pb.Apparel{
				Id:          int32(apparel.ID),
				Name:        apparel.Name,
				Price:       int32(apparel.Price),
				Image:       apparel.ImageURL,
				Subcategory: apparel.SubCategory,
			}
			pbApparels = append(pbApparels, pbApparel)
		}

		response := &pb.SearchApparelsResponse{
			Apparels: pbApparels,
		}
		return response, nil
	}
}

func (s *UserDashboard) UserDetails(ctx context.Context, req *pb.UserDetailsRequset) (*pb.UserDetailsResponse, error) {
	user, err := repo.GetByID(int(req.Userid))
	if err != nil {
		return nil, err
	}
	address, err := repo.GetAddressByUserId(int(req.Userid))
	if err != nil {
		return nil, err
	}
	if user != nil && address != nil {
		userResp := &pb.User{
			Id:        int32(user.ID),
			Firstname: user.FirstName,
			Lastname:  user.LastName,
			Phone:     user.Phone,
			Email:     user.Email,
			Wallet:    int32(user.Wallet),
		}
		addressResp := &pb.Address{
			House:     address.House,
			City:      address.City,
			Pincode:   strconv.Itoa(address.Pincode),
			Type:      address.Type,
			Addressid: int32(address.ID),
		}
		return &pb.UserDetailsResponse{User: userResp, Address: addressResp}, nil
	} else {
		return nil, errors.New("user with this id not found")
	}
}
func (s *UserDashboard) Wishlist(ctx context.Context, req *pb.WishlistRequest) (*pb.WishlistResponse, error) {
	wishlist, err := repo.GetWishlist(int(req.Userid))
	if err != nil {
		return nil, err
	}
	var pbApparels []*pb.Apparel
	for _, apparel := range *wishlist {
		pbApparel := &pb.Apparel{
			Id:          int32(apparel.ID),
			Name:        apparel.ProductName,
			Price:       int32(apparel.Price),
			Subcategory: apparel.Category,
		}
		pbApparels = append(pbApparels, pbApparel)
	}

	response := &pb.WishlistResponse{
		Appares: pbApparels,
	}
	return response, nil
}
