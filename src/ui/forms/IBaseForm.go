package forms

// IBaseForm is an interface that defines methods for basic form operations.
type IBaseForm interface {

	// InitControls initializes the controls of the form.
	//
	// Parameters:
	//   None
	//
	// Returns:
	//   None
	InitControls()

	// Show displays the form.
	//
	// Parameters:
	//   None
	//
	// Returns:
	//   None
	Show()
}
