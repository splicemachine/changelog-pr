package provider

//import errors to log errors when they occur
import (
	"errors"
)

// Provider = The main interface used to describe appliances
type Provider interface {
	GetChangeLogFromPR(sourcePath string, sinceTag string, releaseTag string, auth AuthToken, fileName string) (string, error)
}

type AuthToken struct {
	GithubToken string
}

// Provider Types
const (
	GITHUB = "github"
	MOCK   = "mock"
)

// GetProvider - Function to create the appliances
func GetProvider(t string) (Provider, error) {
	//Use a switch case to switch between types, if a type exist then error is nil (null)
	switch t {
	case GITHUB:
		return new(Github), nil
	case MOCK:
		return new(Mock), nil
	default:
		//if type is invalid, return an error
		return nil, errors.New("unsupported provider")
	}
}
