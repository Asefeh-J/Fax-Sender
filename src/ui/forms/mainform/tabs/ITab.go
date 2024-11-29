package tabs

import "fyne.io/fyne/container"

// ITab is an interface that defines methods for managing UI tabs.
type ITab interface {
	// initUI initializes the UI components of the tab.
	//
	// This method is responsible for setting up the user interface components
	// specific to the tab. It should be called when the tab is created or
	// needs to be refreshed.
	//
	// Parameters:
	//   None
	//
	// Returns:
	//   None
	initUI()

	// loadData loads data for the tab.
	//
	// This method is responsible for loading data specific to the tab.
	// Implementations should perform necessary actions to retrieve and
	// populate the data.
	//
	// Returns:
	//   - bool: True if data is loaded successfully, false otherwise.
	loadData() bool

	// GetTab returns the TabItem associated with the tab.
	//
	// This method is responsible for providing the TabItem that represents
	// the tab. The returned TabItem will be added to the TabContainer in the UI.
	//
	// Returns:
	//   - *container.TabItem: The TabItem associated with the tab.
	GetTab() *container.TabItem

	// setSignalFunc sets the signal function for the tab.
	//
	// This method sets the signal function that will be called when a signal
	// specific to the tab is emitted. It allows communication between different
	// components of the application.
	//
	// Parameters:
	//   - signalFunc: The signal function to set.
	//
	// Returns:
	//   None
	setSignalFunc(signalFunc TabManagementSignal)
}
