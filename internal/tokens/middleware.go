package tokens

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/tofutf/tofutf/internal"
	otfapi "github.com/tofutf/tofutf/internal/api"
	"github.com/tofutf/tofutf/internal/http/html"
	"github.com/tofutf/tofutf/internal/http/html/paths"
	"github.com/tofutf/tofutf/internal/tfeapi"
	"google.golang.org/api/idtoken"
)

const (
	// HTTP header in Google Cloud IAP request containing JWT
	googleIAPHeader string = "x-goog-iap-jwt-assertion"
)

// AuthenticatedPrefixes are those URL path prefixes requiring authentication.
var AuthenticatedPrefixes = []string{
	tfeapi.APIPrefixV2,
	tfeapi.APIPrefixV1,
	tfeapi.ModuleV1Prefix,
	otfapi.DefaultBasePath,
	paths.UIPrefix,
}

type (
	middlewareOptions struct {
		GoogleIAPConfig

		logger *slog.Logger
		key    jwk.Key

		*registry
	}

	GoogleIAPConfig struct {
		Audience string
	}

	middleware struct {
		middlewareOptions
	}
)

// newMiddleware constructs middleware that verifies that all requests
// to protected endpoints possess a valid token, applying the following logic:
//
// 1. Skip authentication for non-protected paths and allow request.
// 2. If Google IAP header is present then authenticate its token and allow or deny
// accordingly.
// 3. If Bearer token is present then authenticate it and allow or deny accordingly.
// 4. If requested path is for a UI endpoint then check for session cookie. If
// present then authenticate its token. If cookie is missing or authentication fails
// then redirect user to login page.
// 5. Otherwise, return 401
//
// Where authentication succeeds, the authenticated subject is attached to the request
// context and the upstream handler is called. If the authenticated subject is a
// user and the user does not exist the user is first created.
func newMiddleware(opts middlewareOptions) mux.MiddlewareFunc {
	mw := middleware{middlewareOptions: opts}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				subject internal.Subject
				err     error
			)
			// Until request is authenticated, call service endpoints using
			// superuser privileges. Once authenticated, the authenticated user
			// replaces the superuser in the context.
			ctx := internal.AddSubjectToContext(r.Context(), &internal.Superuser{
				Username: "auth",
			})

			if !isProtectedPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			if token := r.Header.Get(googleIAPHeader); token != "" {
				subject, err = mw.validateIAPToken(ctx, token)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
			} else if bearer := r.Header.Get("Authorization"); bearer != "" {
				subject, err = mw.validateBearer(ctx, bearer)
				if err != nil {
					mw.logger.Error("validating bearer token", "err", err)
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
			} else if strings.HasPrefix(r.URL.Path, paths.UIPrefix) {
				var ok bool
				subject, ok = mw.validateUIRequest(ctx, w, r)
				if !ok {
					html.SendUserToLoginPage(w, r)
					return
				}
			} else {
				http.Error(w, "no authentication token found", http.StatusUnauthorized)
				return
			}
			ctx = internal.AddSubjectToContext(r.Context(), subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *middleware) validateIAPToken(ctx context.Context, token string) (internal.Subject, error) {
	payload, err := idtoken.Validate(ctx, token, m.Audience)
	if err != nil {
		return nil, err
	}
	email, ok := payload.Claims["email"]
	if !ok {
		return nil, fmt.Errorf("IAP token is missing email claim")
	}
	return m.GetOrCreateUISubject(ctx, email.(string))
}

func (m *middleware) validateBearer(ctx context.Context, bearer string) (internal.Subject, error) {
	splitToken := strings.Split(bearer, "Bearer ")
	if len(splitToken) != 2 {
		return nil, fmt.Errorf("malformed bearer token")
	}
	token := splitToken[1]

	if m.SiteToken != "" && m.SiteToken == token {
		return m.SiteAdmin, nil
	}
	//
	// parse jwt and verify signature
	parsed, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256, m.key))
	if err != nil {
		return nil, err
	}
	kindClaim, ok := parsed.Get("kind")
	if !ok {
		return nil, fmt.Errorf("missing claim: kind")
	}
	kind := Kind(kindClaim.(string))
	return m.GetSubject(ctx, kind, parsed.Subject())
}

func (m *middleware) validateUIRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (internal.Subject, bool) {
	cookie, err := r.Cookie(SessionCookie)
	if errors.Is(err, http.ErrNoCookie) {
		html.FlashSuccess(w, "you need to login to access the requested page")
		return nil, false
	}

	// parse jwt from cookie and verify signature
	token, err := jwt.Parse([]byte(cookie.Value), jwt.WithKey(jwa.HS256, m.key))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired()) {
			html.FlashError(w, "session expired")
		} else {
			html.FlashError(w, "unable to verify session token: "+err.Error())
		}
		return nil, false
	}
	user, err := m.GetOrCreateUISubject(ctx, token.Subject())
	if err != nil {
		html.FlashError(w, "unable to find user: "+err.Error())
		return nil, false
	}
	return user, true
}

func isProtectedPath(path string) bool {
	for _, prefix := range AuthenticatedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
