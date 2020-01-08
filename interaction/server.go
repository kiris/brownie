package interaction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"net/http"
	"net/url"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	verificationToken string
	listenAddr        string
	server            *http.Server
	handlers          map[string]Handler
}


type Handler interface {
	ServInteraction(w http.ResponseWriter, message *slack.InteractionCallback) error
}

func CreateServer(verificationToken string, listenAddr string) *Server {
	server := &http.Server{
		Addr:         listenAddr,
		// Handler:      tracing(nextRequestID)(logging(logger)(router)),
		// ErrorLog:     log,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		verificationToken: verificationToken,
		listenAddr       : listenAddr,
		server           : server,
		handlers         : map[string]Handler{},
	}
}


func (s *Server) ListenAndServ() error {
	http.HandleFunc("/interaction", s.handleInteraction)

	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	log.Info("Server is ready to handle requests at", s.listenAddr)
	return nil
}

func (s *Server) Handle(actionName string, handler Handler) {
	s.handlers[actionName] = handler
}

func (s *Server) handleInteraction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.WithField("method", r.Method).Warn("invalid method request.")
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


	if len(message.Message.Attachments) > 0 {
		for i, a := range message.Message.Attachments[0].Actions {
			println(fmt.Sprintf("message actions[%d] name = %s, value = %s", i, a.Name, a.Value))

			for si, sa := range a.SelectedOptions {
				println(fmt.Sprintf("message selected[%d] value = %s, text = %s", si, sa.Value, sa.Text))
			}
		}
	}

	for i, a := range message.OriginalMessage.Attachments[0].Actions {
		println(fmt.Sprintf("original actions[%d] name = %s, value = %s", i, a.Name, a.Value))

		for si, sa := range a.SelectedOptions {
			println(fmt.Sprintf("original selected[%d] value = %s, text = %s", si, sa.Value, sa.Text))
		}
	}

	for i, a := range message.ActionCallback.AttachmentActions {
		println(fmt.Sprintf("ActionCallback actions[%d] name = %s, value = %s", i, a.Name, a.Value))

		for si, sa := range a.SelectedOptions {
			println(fmt.Sprintf("ActionCallback selected[%d] value = %s, text = %s", si, sa.Value, sa.Text))
		}
	}

	actionName := message.ActionCallback.AttachmentActions[0].Name
	handler, ok := s.handlers[actionName]
	if !ok {
		// unknown command.
		return
	}

	if err := handler.ServInteraction(w, &message); err != nil {
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

