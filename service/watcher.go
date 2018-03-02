package service

// WatchNodes watches for node events and notifies addressses
func WatchNodes(
	el EventListening, ni NodeInspector, nen NodeEventNotifing) {
}

// WatchServices watches for service events and notifies addresses
func WatchServices(
	el EventListening, s *Service, n *Notification) {
	// Check if services are listening for notifcations
}
