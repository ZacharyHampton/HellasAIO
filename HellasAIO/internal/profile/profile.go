package profile

import (
	"errors"
	"github.com/lithammer/shortuuid"
	"sync"
)

var (
	profileMutex = sync.RWMutex{}

	ProfileDoesNotExistErr = errors.New("profile does not exist")
	ProfileNotAssignedErr  = errors.New("profile not assigned")
	profiles               = make(map[string]*Profile)
)

// DoesProfileExist checks if a profile exists
func DoesProfileExist(id string) bool {
	profileMutex.RLock()
	defer profileMutex.RUnlock()

	_, ok := profiles[id]
	return ok
}

// CreateProfile creates a new profile
func CreateProfile(profile *Profile) string {
	profileMutex.RLock()
	defer profileMutex.RUnlock()

	id := shortuuid.New()

	profiles[id] = profile

	return id
}

// RemoveProfile removes a profile
func RemoveProfile(id string) error {
	if !DoesProfileExist(id) {
		return ProfileDoesNotExistErr
	}

	profileMutex.RLock()
	defer profileMutex.RUnlock()

	delete(profiles, id)

	return nil
}

// GetProfileById gets a profile by id
func GetProfileById(id string) (*Profile, error) {
	if !DoesProfileExist(id) {
		return &Profile{}, ProfileDoesNotExistErr
	}

	profileMutex.RLock()
	defer profileMutex.RUnlock()

	return profiles[id], nil
}

// GetAllProfileIDs gets all profile ids
func GetAllProfileIDs() []string {
	ids := []string{}

	for id := range profiles {
		ids = append(ids, id)
	}

	return ids
}

func GetProfileIDByName(name string) (string, error) {
	for _, profileId := range GetAllProfileIDs() {
		profile, err := GetProfileById(profileId)
		if err != nil {
			return "", err
		}

		if profile.Name == name {
			return profileId, nil
		}
	}

	return "", ProfileDoesNotExistErr
}
