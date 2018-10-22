package httpsignatures

//KeyID define ID of key
type KeyID string

//Permission list KeyID that have permission
type Permission []KeyID

//isPermit check if keyID in permissions list
func (permissions Permission) isPermit(keyID KeyID) bool {
	for _, p := range permissions {
		if keyID == p {
			return true
		}
	}
	return false
}
