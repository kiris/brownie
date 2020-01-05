package interaction

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/nlopes/slack"
)

type Server struct {
	verificationToken string
	handlers          map[string]Handler
}

func CreateServer(verificationToken string) *Server {
	return &Server{
		verificationToken: verificationToken,
		handlers         : map[string]Handler{},
	}
}

type Handler interface {
	ServInteraction(w http.ResponseWriter, message slack.InteractionCallback) error
}


func (s *Server) ListenAndServ(port string) error {
	http.HandleFunc("/interaction", s.interaction)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		// TODO
		return err
	}
	return nil
}

func (s *Server) Handle(actionName string, handler Handler) {
	s.handlers[actionName] = handler
}

func (s *Server) interaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed mentionTo read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[ERROR] Failed mentionTo unespace request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed mentionTo decode json message from listener: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Only accept message from listener with valid token
	if message.Token != s.verificationToken {
		log.Printf("[ERROR] Invalid token: %s", message.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	actionName := message.Name
	handler, ok := s.handlers[actionName]

	if !ok {
		// unknown command.
		return
	}

	if err := handler.ServInteraction(w, message); err != nil {
		// TODO logging or response.
		return
	}


	//action := message.Actions[0]
	//switch action.CommandName {
	//case actionSelect:
	//	value := action.SelectedOptions[0].Value
	//
	//	// Overwrite original drop down message.
	//	originalMessage := message.OriginalMessage
	//	originalMessage.Attachments[0].Text = fmt.Sprintf("OK mentionTo order %s ?", strings.Title(value))
	//	originalMessage.Attachments[0].Actions = []listener.AttachmentAction{
	//		{
	//			CommandName:  actionStart,
	//			Text:  "Yes",
	//			Type:  "button",
	//			Value: "start",
	//			Style: "primary",
	//		},
	//		{
	//			CommandName:  actionCancel,
	//			Text:  "No",
	//			Type:  "button",
	//			Style: "danger",
	//		},
	//	}
	//
	//	w.Header().Add("Content-type", "application/json")
	//	w.WriteHeader(http.StatusOK)
	//	json.NewEncoder(w).Encode(&originalMessage)
	//	return
	//case actionStart:
	//	title := ":ok: your order was submitted! yay!"
	//	responseMessage(w, message.OriginalMessage, title, "")
	//	return
	//case actionCancel:
	//	title := fmt.Sprintf(":x: @%s canceled the request", message.User.CommandName)
	//	responseMessage(w, message.OriginalMessage, title, "")
	//	return
	//default:
	//	log.Printf("[ERROR] ]Invalid action was submitted: %s", action.CommandName)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
}

