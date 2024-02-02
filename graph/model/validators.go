package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (product ProductObj) Validate() error {
	return validation.ValidateStruct(&product,
		validation.Field(&product.Name, validation.Required, validation.Length(3, 50)),
		validation.Field(&product.Description, validation.Required, validation.Length(3, 50)),
		validation.Field(&product.Model, validation.Required),
		validation.Field(&product.Price, validation.Required),
	)
}

func (customerCart CustomerCartObj) Validate() error {
	return validation.ValidateStruct(&customerCart,
		validation.Field(&customerCart.CustomerID, validation.Required),
		validation.Field(&customerCart.ProductID, validation.Required),
	)
}

func (addUser AddUserObj) Validate() error {
	return validation.ValidateStruct(&addUser,
		validation.Field(&addUser.Fullname, validation.Required),
		validation.Field(&addUser.Email, validation.Required),
		validation.Field(&addUser.City, validation.Required),
		validation.Field(&addUser.ZipCode, validation.Required),
		validation.Field(&addUser.IsActive, validation.Required),
		validation.Field(&addUser.Mobile, validation.Required),
		validation.Field(&addUser.PaymentMethod, validation.Required),
		validation.Field(&addUser.Role, validation.Required),
		validation.Field(&addUser.Username, validation.Required),
		validation.Field(&addUser.StreetNo, validation.Required),
	)
}

func (inventory InventoryObj) Validate() error {
	return validation.ValidateStruct(&inventory,
		validation.Field(&inventory.ProductID, validation.Required),
		validation.Field(&inventory.Quantity, validation.Required),
	)
}
