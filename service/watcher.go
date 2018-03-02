package service

// WatchNodes watches for node events and notifies addressses
func WatchNodes(
	el EventListening, ni NodeInspector, nen EventNodeNotifing) {
}

// WatchServices watches for service events and notifies addresses
func WatchServices(
	el EventListening, s Servicer, nsn EventServiceNotifier) {
	// Check if services are listening for notifcations
}
