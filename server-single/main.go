package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"

	"connectrpc.com/connect"
	"github.com/quic-go/quic-go/http3"

	"buf.build/gen/go/connectrpc/eliza/connectrpc/go/connectrpc/eliza/v1/elizav1connect"
	elizav1 "buf.build/gen/go/connectrpc/eliza/protocolbuffers/go/connectrpc/eliza/v1"
)

var introResponses = []string{
	"Hi %s. I'm Eliza.",
	"Before we begin, %s, let me tell you something about myself.",
}

// A string array of facts about ELIZA.  Used in responses to Introduce, which is a server-stream.
var elizaFacts = []string{
	"I was created by Joseph Weizenbaum.",
	"I was created in the 1960s.",
	"I am a Rogerian psychotherapist.",
	"I am named after Eliza Doolittle from the play Pygmalion.",
	"I was originally written on an IBM 7094.",
	"I can be accessed in most Emacs implementations with the command M-x doctor.",
	"I was created at the MIT Artificial Intelligence Laboratory.",
	"I was one of the first programs capable of attempting the Turing test.",
	"I was designed as a method to show the superficiality of communication between man and machine.",
}

var _ elizav1connect.ElizaServiceHandler = (*server)(nil)

type server struct {
	elizav1connect.UnimplementedElizaServiceHandler
}

// Say implements elizav1connect.ElizaServiceHandler.
func (s *server) Say(ctx context.Context, req *connect.Request[elizav1.SayRequest]) (*connect.Response[elizav1.SayResponse], error) {
	slog.Info("Say()", "req", req)
	return connect.NewResponse(&elizav1.SayResponse{
		Sentence: req.Msg.GetSentence(),
	}), nil
}

func randomElementFrom(list []string) string {
	return list[rand.Intn(len(list))] //nolint:gosec
}

// GetIntroResponses returns a collection of introductory responses tailored to the given name.
func GetIntroResponses(name string) []string {
	intros := make([]string, 0, len(introResponses)+2)
	for _, n := range introResponses {
		intros = append(intros, fmt.Sprintf(n, name))
	}

	intros = append(intros, randomElementFrom(elizaFacts))
	intros = append(intros, "How are you feeling today?")
	return intros
}

// Introduce implements elizav1connect.ElizaServiceHandler.
func (e *server) Introduce(
	ctx context.Context,
	req *connect.Request[elizav1.IntroduceRequest],
	stream *connect.ServerStream[elizav1.IntroduceResponse],
) error {
	name := req.Msg.GetName()
	if name == "" {
		name = "Anonymous User"
	}
	intros := GetIntroResponses(name)
	for _, resp := range intros {
		if err := stream.Send(&elizav1.IntroduceResponse{Sentence: resp}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	mux := http.NewServeMux()
	mux.Handle(elizav1connect.NewElizaServiceHandler(&server{}))

	addr := "0.0.0.0:6660"
	log.Printf("Starting connectrpc on %s", addr)
	h3srv := http3.Server{
		Addr:    addr,
		Handler: mux,
	}
	if err := h3srv.ListenAndServeTLS("cert.crt", "cert.key"); err != nil {
		log.Fatalf("error: %s", err)
	}
}
