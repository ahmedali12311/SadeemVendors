package main

import (
	"errors"
	"net/http"
	"project/internal/data"
	"project/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

// CreateCartHandler handles the creation of a new cart.
func (app *application) CreateCartHandler(w http.ResponseWriter, r *http.Request) {

	err := app.Model.CartDB.InsertCart(&data.Cart{ID: uuid.MustParse(r.Context().Value(UserIDKey).(string))})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, utils.Envelope{"created cart for user ": r.Context().Value(UserIDKey).(string)})
}

// DeleteCartHandler handles the deletion of a cart by its ID.
func (app *application) DeleteCartHandler(w http.ResponseWriter, r *http.Request) {
	cartID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid cart ID"))
		return
	}

	err = app.Model.CartDB.DeleteCart(cartID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"message": "cart deleted successfully"})
}
func (app *application) UpdateCartHandler(w http.ResponseWriter, r *http.Request) {

	// Get cart ID from request
	cartIDStr := r.PathValue("id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid cart ID"))
		return
	}

	// Fetch the current cart
	cart, err := app.Model.CartDB.GetCart(cartID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if cart.Quantity != 0 {
		cart.Description = r.FormValue("description")
		if len(cart.Description) > 100 {
			app.errorResponse(w, r, http.StatusBadRequest, "description too long!")
		}
		err = app.Model.CartDB.UpdateCart(cart)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

	}
	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"cart": cart})
}

func (app *application) GetCartHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from the context (assuming you have middleware that sets this)
	userID := uuid.MustParse(r.Context().Value(UserIDKey).(string))

	// Retrieve all items in the cart
	cartItems, err := app.Model.CartDB.GetCart(userID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"cart": nil})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"cart": cartItems})
}
func (app *application) CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse(r.Context().Value(UserIDKey).(string))

	// Check if the customer has a table
	_, err := app.Model.TableDB.GetCustomertable(r.Context(), userID)
	if err != nil {
		app.handleRetrievalError(w, r, err)
		return
	}

	// Fetch the cart for the user
	cart, err := app.Model.CartDB.GetCart(userID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Fetch all cart items
	cartItems, err := app.Model.CartItemDB.GetCartItems(cart.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Check if cart is empty
	if len(cartItems) == 0 {
		app.errorResponse(w, r, http.StatusBadRequest, "Cart is empty.")
		return
	}

	// Ensure all items are from the same vendor
	vendorID := cart.VendorID
	for _, item := range cartItems {
		itemData, err := app.Model.ItemDB.GetItem(item.ItemID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if itemData.VendorID != vendorID {
			app.errorResponse(w, r, http.StatusBadRequest, "All items in the cart must be from the same vendor.")
			return
		}
	}

	// Begin a transaction
	tx, err := app.Model.BeginTransaction()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback() // Rollback on error
		} else {
			err = tx.Commit() // Commit if no errors
		}
	}()

	// Create the order
	order := &data.Order{
		ID:             uuid.New(),
		TotalOrderCost: cart.TotalPrice,
		CustomerID:     userID,
		VendorID:       vendorID,
		Status:         "preparing",
		Description:    cart.Description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	// Insert the order into the database
	if err = tx.InsertOrder(order); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Insert order items
	for _, item := range cartItems {
		// Fetch the item to get its price
		itemData, err := app.Model.ItemDB.GetItem(item.ItemID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		var price float64
		if itemData.Discount != 0 {
			// Calculate the discounted price
			price = itemData.Discount
		} else {
			// Use the original price if there's no discount
			price = itemData.Price
		}

		// Create order item with the current price
		orderItem := &data.OrderItem{
			ID:       uuid.New(),
			OrderID:  order.ID,
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
			Price:    price, // Use the current price of the item
		}

		// Insert order item into the database
		if err = tx.InsertOrderItem(orderItem); err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Update item quantity in the inventory
		newQuantity := itemData.Quantity - item.Quantity
		if newQuantity < 0 {
			app.serverErrorResponse(w, r, errors.New("not enough quantity"))
			return
		}
		itemData.Quantity = newQuantity

		// Remove domain prefix from image URL if present
		if itemData.Img != nil {
			*itemData.Img = strings.TrimPrefix(*itemData.Img, data.Domain+"/")
		}

		// Update the item in the database with the new quantity
		if err = app.Model.ItemDB.UpdateItem(itemData); err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// Delete the cart and cart items
	if err = tx.DeleteCart(cart.ID); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err = tx.DeleteCartItems(cart.ID); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Respond with a success message and the order ID
	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"message": "checkout successful", "order_id": order.ID})
}
