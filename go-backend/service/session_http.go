package service

import (
	// "encoding/json"
	// "io/ioutil"
	ae "joshsoftware/peerly/apperrors"
	"joshsoftware/peerly/config"
	"joshsoftware/peerly/db"
	log "joshsoftware/peerly/util/log"
	"net/http"
	// "net/url"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	logger "github.com/sirupsen/logrus"
)

// OAuthUser - a struct that represents the "user" we'll get back from Google's /userinfo query
type OAuthUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Domain        string `json:"hd"`
	VerifiedEmail bool   `json:"verified_email"`
	PictureURL    string `json:"picture"`
}

// OAuthToken - a struct used to json.Unmarshal the response body from an oauth provider
type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

// authBody - a struct we use for marshalling into JSON to send down as the response body after a user has been
// successfully authenticated and they need their token for using the app in subsequent API requests
// (see the 'handleAuth' function below).
type authBody struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}



// newJWT() - Creates and returns a new JSON Web Token to be sent to an API consumer on valid
// authentication, so they can re-use it by sending it in the Authorization header on subsequent
// requests.
func newJWT(userID, orgID int) (newToken string, err error) {
	signingKey := config.JWTKey()
	if signingKey == nil {
		log.Error(ae.ErrNoSigningKey, "Application error: No signing key configured", err)
		return
	}

	expiryTime := time.Now().Add(time.Duration(config.JWTExpiryDurationHours()) * time.Hour).Unix()
	claims := &jwt.MapClaims{
		"exp": expiryTime,
		"iss": "joshsoftware.com",
		"iat": time.Now().Unix(),
		"sub": strconv.Itoa(userID),
		"org": strconv.Itoa(orgID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newToken, err = token.SignedString(signingKey)
	if err != nil {
		log.Error(ae.ErrSignedString, "Failed to get signed string", err)
		return
	}
	return
}

func handleLogout(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		parsedToken, ok := req.Context().Value("user").(*jwt.Token)
		if !ok {
			logger.Error("Error parsing JSON for token response")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		claims := parsedToken.Claims.(jwt.MapClaims)

		userID, err := strconv.Atoi(claims["sub"].(string))
		if err != nil {
			log.Error(ae.ErrJSONParseFail, "Error parsing JSON for token response", err)
			ae.JSONError(rw, http.StatusInternalServerError, err)
			return
		}

		expirationTimeStamp := int64(claims["exp"].(float64))

		expirationDate := time.Unix(expirationTimeStamp, 0)
		userBlackListedToken := db.UserBlacklistedToken{
			UserID:         userID,
			ExpirationDate: expirationDate,
			Token:          parsedToken.Raw,
		}

		err = deps.Store.CreateUserBlacklistedToken(req.Context(), userBlackListedToken)
		if err != nil {
			log.Error(ae.ErrFailedToCreate, "Error creating blaclisted token record", err)
			rw.Header().Add("Content-Type", "application/json")
			ae.JSONError(rw, http.StatusInternalServerError, err)
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		return
	})
}
