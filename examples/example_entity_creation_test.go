// package main_test contains examples of GoCTI usage.
//
// This example shows the basic functionality and usage principles of GoCTI.
package main_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/kelseyhightower/envconfig"

	"github.com/weisshorn-cyd/gocti"
	"github.com/weisshorn-cyd/gocti/entity"
)

func Example_entityCreation() {
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

	// Create a malware entity
	// The customAttributes function parameter is the empty string because we
	// are not interested in customizing which attributes of entity.Malware
	// should be set. The empty string means that the default attributes will be set.
	malware, err := client.CreateMalware(ctx, "", entity.MalwareAddInput{
		Name:        "Example Malware",
		Description: "An example of a very malicious malware.",
	})
	if err != nil {
		logger.Error("creating malware", "error", err)
	}

	// Since this is an example, we clean up after ourselves
	defer func() {
		if _, err := client.DeleteMalware(ctx, malware.ID); err != nil {
			logger.Error("deleting", "malware", malware, "error", err)
		}
	}()

	fmt.Printf("We created a malware named: '%s'\n", malware.Name)

	// List all malwares
	// We want only the ID to be returned; we override the default properties by using the customAttributes field
	// There are few entites in our platform; we can specify that we want all values returned instead of using pagination.
	// We do not require an options argument (such as filters, etc)
	malwares, err := client.ListMalwares(ctx, "id, name", true, nil)
	if err != nil {
		logger.Error("listing malwares", "error", err)
	}

	found := false

	for _, malware := range malwares {
		if malware.Name == "Example Malware" {
			found = true
		}
	}

	if found {
		fmt.Println("We found our example malware.")
	} else {
		fmt.Println("We didn't find our example malware.")
	}

	// Create an instrusion set
	intrusionSet, err := client.CreateIntrusionSet(ctx, "", entity.IntrusionSetAddInput{
		Name: "Example Intrusion Set",
	})
	if err != nil {
		logger.Error("creating intrusion set", "error", err)
	}

	defer func() {
		if _, err := client.DeleteIntrusionSet(ctx, intrusionSet.ID); err != nil {
			logger.Error("deleting", "intrusion set", intrusionSet, "error", err)
		}
	}()

	fmt.Printf("We created an intrusion set named: '%s'\n", intrusionSet.Name)

	// Mark the intrusion set as using the example malware by creating the
	// required relationship between the two entities
	_, err = client.CreateStixCoreRelationship(ctx, "", entity.StixCoreRelationshipAddInput{
		FromID:           intrusionSet.ID,
		ToID:             malware.ID,
		RelationshipType: "uses",
	})
	if err != nil {
		logger.Error("creating relationship", "error", err)
	}

	// Update the existing intrusion set
	// Reuse the key fields to force the update
	//
	// Note: An actual update function will be added in a future release of GoCTI
	updatedIntrusionSet, err := client.CreateIntrusionSet(ctx, "", entity.IntrusionSetAddInput{
		Name:        intrusionSet.Name,      // Required field
		Description: "I have been updated.", // An updated description
	})
	if err != nil {
		logger.Error("updating intrusion set", "error", err)
	}

	fmt.Printf("The updated intrusion set has description: '%s'\n", updatedIntrusionSet.Description)

	// If the ID is known, entities can be read directly
	// Since we do not provide values for customAttributes, all the fields
	// contained in the default attributes will be set.
	intrusionSet, err = client.ReadIntrusionSet(ctx, "", intrusionSet.ID)
	if err != nil {
		logger.Error("reading intrusion set", "error", err)
	}

	fmt.Printf("Intrusion set has description: '%s'\n", intrusionSet.Description)

	// Output:
	// We created a malware named: 'Example Malware'
	// We found our example malware.
	// We created an intrusion set named: 'Example Intrusion Set'
	// The updated intrusion set has description: 'I have been updated.'
	// Intrusion set has description: 'I have been updated.'
}
