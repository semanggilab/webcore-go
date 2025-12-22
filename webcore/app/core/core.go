package core

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/semanggilab/webcore-go/app/config"
	"github.com/semanggilab/webcore-go/app/logger"
)

// Context represents shared dependencies that can be injected into modules
type AppContext struct {
	Context  context.Context
	Config   *config.Config
	Web      *fiber.App
	Root     fiber.Router
	EventBus *EventBus
}

func (a *AppContext) Start() error {
	libmanager := Instance().LibraryManager

	// Initialize database if configured
	if a.Config.Database.Host != "" {
		// a.SetupDatabase("default", a.Config.Database)
		lName := "db:" + a.Config.Database.Driver
		loader, ok := libmanager.GetLoader(lName)
		if !ok {
			return fmt.Errorf("LibraryLoader '%s' tidak ditemukan", lName)
		}

		_, err := libmanager.LoadSingletonFromLoader(loader, a.Context, a.Config.Database)
		if err != nil {
			return err
		}
	}

	// Initialize Redis if configured
	if a.Config.Redis.Host != "" {
		// a.SetupRedis(a.Config.Redis)
		loader, ok := libmanager.GetLoader("redis")
		if !ok {
			return fmt.Errorf("LibraryLoader 'redis' tidak ditemukan")
		}
		_, err := libmanager.LoadSingletonFromLoader(loader, a.Context, a.Config.Database)
		if err != nil {
			return err
		}
	}

	// Initialize PubSub if configured
	if a.Config.PubSub.ProjectID != "" && a.Config.PubSub.Topic != "" {
		// a.SetupPubSub("default", a.Config.PubSub)
		loader, ok := libmanager.GetLoader("pubsub")
		if !ok {
			return fmt.Errorf("LibraryLoader 'pubsub' tidak ditemukan")
		}
		_, err := libmanager.LoadSingletonFromLoader(loader, a.Context, a.Config.Database)
		if err != nil {
			return err
		}
	}

	return nil
}

// Destroy release all resources
func (a *AppContext) Destroy() error {
	// Shutdown Fiber app
	if a.Web != nil {
		return a.Web.Shutdown()
	}

	return nil
}

func CheckSingleLoader[L any](name string, loaders []L) []L {
	newLoaders := []L{}
	list := []string{}
	for _, loader := range loaders {
		lType := reflect.TypeOf(loader)
		if lType.Kind() == reflect.Ptr {
			lType = lType.Elem()
		}

		lName := lType.Name()
		if slices.Contains(list, lName) {
			logger.Fatal(name+" is registered multiple times", "name", lName)
		}

		newLoaders = append(newLoaders, loader)
	}

	return newLoaders
}

func AppendRouteToArray(routes []*ModuleRoute, route *ModuleRoute) []*ModuleRoute {
	route.Root.Add(route.Method, route.Path, route.Handler)

	routes = append(routes, route)
	return routes
}
