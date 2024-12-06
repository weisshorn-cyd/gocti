// package main_test contains examples of GoCTI usage.
//
// This example shows how to use some helper functions to switch between GraphQL interfaces
// and implementations.
package main_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/weisshorn-cyd/gocti"
	"github.com/weisshorn-cyd/gocti/api"
	"github.com/weisshorn-cyd/gocti/entity"
	"github.com/weisshorn-cyd/gocti/graphql"
)

func Example_decodeGraphQLInterface() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	// Get config from env
	cfg := struct {
		URL   string `envconfig:"URL"   required:"true"`
		Token string `envconfig:"TOKEN" required:"true"`
	}{}
	if err := envconfig.Process("OPENCTI", &cfg); err != nil {
		logger.Error("loading config", "error", err)
		os.Exit(1)
	}

	// Create client
	client, err := gocti.NewOpenCTIAPIClient(
		cfg.URL,
		cfg.Token,
		gocti.WithHealthCheck(),
		gocti.WithLogger(logger),
	)
	if err != nil {
		logger.Error("creating client", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	callback, ids := provisionObservables(ctx, client, logger)
	defer callback()

	for _, id := range ids {
		// Retrieve the observable. As it is of type StixCyberObservable, we don't have access
		// to the implementation-specific fields, we need to decode it using [api.DecodeInterface]
		obs, err := client.ReadStixCyberObservable(ctx, "", id)
		if err != nil {
			logger.Error("reading observable", "error", err)
		}

		switch graphql.StixCyberObservableType(obs.EntityType) {
		case graphql.StixCyberObservableTypeArtifact:
			artifact := graphql.Artifact{}
			if err := api.DecodeInterface(obs, &artifact); err != nil {
				logger.Error("decoding artifact", "error", err)
			}

			fmt.Printf("Artifact has MIME type %q and url %q\n", artifact.MimeType, artifact.URL)

		case graphql.StixCyberObservableTypeIPV4Addr:
			ipv4Addr := graphql.IPv4Addr{}
			if err := api.DecodeInterface(obs, &ipv4Addr); err != nil {
				logger.Error("decoding ipv4 address", "error", err)
			}

			fmt.Printf("IPv4 address has value %q\n", ipv4Addr.Value)

		case graphql.StixCyberObservableTypeEmailMessage:
			emailMessage := graphql.EmailMessage{}
			if err := api.DecodeInterface(obs, &emailMessage); err != nil {
				logger.Error("decoding email message", "error", err)
			}

			fmt.Printf("Email message has subject %q and body %q\n", emailMessage.Subject, emailMessage.Body)

			// If you need to fill in a field of type StixCyberObservable,
			// you can decode the implementation into the target interface:
			stixCyberObservable := graphql.StixCyberObservable{}

			if err := api.Decode(emailMessage, &stixCyberObservable); err != nil {
				logger.Error("decoding email into observable", "error", err)
			}

			subject, ok := stixCyberObservable.Remain["subject"]
			if !ok {
				logger.Error("cannot find the key \"subject\"")
			}

			fmt.Printf("Observable content: Remain: \"%v\"", subject)
		}
	}

	// Output:
	// Artifact has MIME type "application/json" and url "url.random.com"
	// IPv4 address has value "987.987.987.987"
	// Email message has subject "A cool Email message" and body "I bought a new hamster !"
	// Observable content: Remain: "A cool Email message"
}

func provisionObservables(ctx context.Context, client *gocti.OpenCTIAPIClient, logger *slog.Logger) (func(), []string) {
	type deletionTask struct {
		id         string
		deleteFunc func(ctx context.Context, id string) (string, error)
	}

	deletionTasks := []deletionTask{}
	ids := []string{}

	// List of some observables with fields not available on a generic StixCyberObservable
	observables := []entity.StixCyberObservableAddInput{
		{
			Type: graphql.StixCyberObservableTypeArtifact,
			Artifact: graphql.ArtifactAddInput{
				MimeType: "application/json",
				URL:      "url.random.com",
			},
		}, {
			Type: graphql.StixCyberObservableTypeIPV4Addr,
			IPv4Addr: graphql.IPv4AddrAddInput{
				Value: "987.987.987.987",
			},
		}, {
			Type: graphql.StixCyberObservableTypeEmailMessage,
			EmailMessage: graphql.EmailMessageAddInput{
				Subject: "A cool Email message",
				Body:    "I bought a new hamster !",
			},
		},
	}

	for _, observable := range observables {
		obs, err := client.CreateStixCyberObservable(ctx, "id", observable)
		if err != nil {
			logger.Error("provisioning ioc", "type", observable.Type, "error", err)
		}

		deletionTasks = append(deletionTasks, deletionTask{
			id:         obs.ID,
			deleteFunc: client.DeleteStixCyberObservable,
		})

		ids = append(ids, obs.ID)
	}

	return func() {
		for _, task := range deletionTasks {
			if _, err := task.deleteFunc(ctx, task.id); err != nil {
				logger.Error("cleaning up", "error", err)
			}
		}
	}, ids
}
