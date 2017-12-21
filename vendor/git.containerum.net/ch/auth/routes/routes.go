package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"git.containerum.net/ch/auth/utils"
	"git.containerum.net/ch/grpc-proto-files/auth"
	"github.com/husobee/vestigo"
	"github.com/opentracing/opentracing-go"
)

// Headers used in REST mode
const (
	HeaderUserAgent   = "X-User-Agent"
	HeaderFingerprint = "X-User-Fingerprint"
	HeaderUserId      = "X-User-ID"
	HeaderUserIp      = "X-User-IP"
	HeaderUserRole    = "X-User-Role"
	HeaderPartTokenId = "X-User-Part-Token"
	HeaderTokenId     = "X-User-Token-ID"
)

func SetupRoutes(router *vestigo.Router, tracer opentracing.Tracer, storage auth.AuthServer) {
	// Create token
	router.Post("/token", createTokenHandler,
		newOpenTracingMiddleware(tracer, "Create Token"),
		newStorageInjectionMiddleware(storage),
		newHeaderValidationMiddleware(standardHeaderValidators),
		newBodyValidationMiddleware(resourcesAccessBodyValidator))

	// Check token
	router.Get("/token/:access_token", checkTokenHandler,
		newOpenTracingMiddleware(tracer, "Check Token"),
		newStorageInjectionMiddleware(storage),
		newHeaderValidationMiddleware(standardHeaderValidators))

	// Extend token
	router.Put("/token/:refresh_token", extendTokenHandler,
		newOpenTracingMiddleware(tracer, "Extend Token"),
		newStorageInjectionMiddleware(storage),
		newHeaderValidationMiddleware(standardHeaderValidators))

	// Get user tokens
	router.Get("/token", getUserTokensHandler,
		newOpenTracingMiddleware(tracer, "Get user tokens"),
		newStorageInjectionMiddleware(storage),
		newHeaderValidationMiddleware(standardHeaderValidators))

	// Delete token by ID
	router.Delete("/token/:token_id", deleteTokenByIdHandler,
		newOpenTracingMiddleware(tracer, "Delete token by ID"),
		newStorageInjectionMiddleware(storage),
		newParameterValidationMiddleware(validators{"token_id": uuidValidator}),
		newHeaderValidationMiddleware(standardHeaderValidators))

	// Delete user tokens
	router.Delete("/token/user/:user_id", deleteUserTokensHandler,
		newStorageInjectionMiddleware(storage),
		newParameterValidationMiddleware(validators{"user_id": uuidValidator}),
		newOpenTracingMiddleware(tracer, "Delete user tokens"))
}

const authServerContextKey = "authServer"

func authServerFromRequestContext(r *http.Request) auth.AuthServer {
	return r.Context().Value(authServerContextKey).(auth.AuthServer)
}

func createTokenHandler(w http.ResponseWriter, r *http.Request) {
	req := &auth.CreateTokenRequest{
		UserAgent:   r.Header.Get(HeaderUserAgent),
		Fingerprint: r.Header.Get(HeaderFingerprint),
		UserId:      utils.UUIDFromString(r.Header.Get(HeaderUserId)),
		UserIp:      r.Header.Get(HeaderUserIp),
		UserRole:    auth.Role(auth.Role_value[r.Header.Get(HeaderUserRole)]),
		PartTokenId: utils.UUIDFromString(r.Header.Get(HeaderPartTokenId)),
	}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(body, &req.Access)

	resp, err := authServerFromRequestContext(r).CreateToken(r.Context(), req)
	if err != nil {
		sendError(w, err)
		return
	}

	body, err = json.Marshal(resp)

	_, err = w.Write(body)
}

func checkTokenHandler(w http.ResponseWriter, r *http.Request) {
	req := &auth.CheckTokenRequest{
		AccessToken: vestigo.Param(r, "access_token"),
		UserAgent:   r.Header.Get(HeaderUserAgent),
		FingerPrint: r.Header.Get(HeaderFingerprint),
		UserIp:      r.Header.Get(HeaderUserIp),
	}

	defer r.Body.Close()

	resp, err := authServerFromRequestContext(r).CheckToken(r.Context(), req)
	if err != nil {
		sendError(w, err)
		return
	}

	var checkTokenResponseBody = struct {
		Access *auth.ResourcesAccess `json:"access"`
	}{
		Access: resp.Access,
	}

	w.Header().Add(HeaderUserId, resp.UserId.Value)
	w.Header().Add(HeaderUserRole, resp.UserRole.String())
	w.Header().Add(HeaderTokenId, resp.TokenId.Value)
	w.Header().Add(HeaderPartTokenId, resp.PartTokenId.Value)

	body, err := json.Marshal(checkTokenResponseBody)

	_, err = w.Write(body)
}

func extendTokenHandler(w http.ResponseWriter, r *http.Request) {
	req := &auth.ExtendTokenRequest{
		RefreshToken: vestigo.Param(r, "refresh_token"),
		Fingerprint:  r.Header.Get(HeaderFingerprint),
	}

	defer r.Body.Close()

	resp, err := authServerFromRequestContext(r).ExtendToken(r.Context(), req)
	if err != nil {
		sendError(w, err)
		return
	}

	body, err := json.Marshal(resp)

	_, err = w.Write(body)
}

func getUserTokensHandler(w http.ResponseWriter, r *http.Request) {
	req := &auth.GetUserTokensRequest{
		UserId: utils.UUIDFromString(r.Header.Get(HeaderUserId)),
	}

	defer r.Body.Close()

	resp, err := authServerFromRequestContext(r).GetUserTokens(r.Context(), req)
	if err != nil {
		sendError(w, err)
		return
	}

	body, err := json.Marshal(resp)

	_, err = w.Write(body)
}

func deleteTokenByIdHandler(w http.ResponseWriter, r *http.Request) {
	req := &auth.DeleteTokenRequest{
		TokenId: utils.UUIDFromString(vestigo.Param(r, "token_id")),
		UserId:  utils.UUIDFromString(r.Header.Get(HeaderUserId)),
	}

	defer r.Body.Close()

	_, err := authServerFromRequestContext(r).DeleteToken(r.Context(), req)
	if err != nil {
		sendError(w, err)
		return
	}

}

func deleteUserTokensHandler(w http.ResponseWriter, r *http.Request) {
	req := &auth.DeleteUserTokensRequest{
		UserId: utils.UUIDFromString(vestigo.Param(r, "user_id")),
	}

	defer r.Body.Close()

	_, err := authServerFromRequestContext(r).DeleteUserTokens(r.Context(), req)
	if err != nil {
		sendError(w, err)
		return
	}
}
