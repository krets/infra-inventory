package inventory

import (
	log "github.com/golang/glog"
	"net/http"
)

/* Auth handler to wrap any handler function */
func AuthHandler(hFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "POST", "PUT", "DELETE":
			token := r.Header.Get("Auth-Token")
			if len(token) < 1 {
				user, _, ok := r.BasicAuth()
				if !ok {
					WriteAndLogResponse(w, r, 401, map[string]string{"Content-Type": "text/plain"}, []byte(`Unauthorized`))
					return
				}
				log.V(8).Infof("Requesting auth for: %s\n", user)
			} else {
				log.V(8).Infof("Checking token validity: %s\n", token)
			}
			break
		default:
			break
		}

		hFunc(w, r)
	}
}

/* Helper function to write http data */
func WriteAndLogResponse(w http.ResponseWriter, r *http.Request, code int, headers map[string]string, data []byte) {
	w.WriteHeader(code)

	if headers != nil {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}

	w.Write(data)
	log.Infof("%s %d %s %d\n", r.Method, code, r.RequestURI, len(data))
}
