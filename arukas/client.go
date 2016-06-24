package arukas

import (
	"fmt"
	API "github.com/arukasio/cli"
	"github.com/yamamoto-febc/arukas-ship/message"
	"time"
)

type ArukasClient struct {
	client *API.Client
}

func NewArukasClient() (*ArukasClient, error) {
	client, err := API.NewClient()
	if err != nil {
		return nil, err
	}
	client.UserAgent = "arukas-ship by github.com/yamamoto-febc/arukas-ship"
	client.Debug = true

	return &ArukasClient{client: client}, nil
}

func (c *ArukasClient) HandleRequest(appName string, msg *message.IncomingMessage) error {
	containers, err := c.listContainerByName(appName)
	if err != nil {
		return err
	}
	if containers == nil || len(containers) == 0 {
		err = c.createContainer(appName, msg)
		if err != nil {
			return err
		}
	} else {
		for _, container := range containers {
			err = c.updateContainer(container, msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ArukasClient) listContainerByName(appName string) ([]API.Container, error) {
	var containers []API.Container

	err := c.client.Get(&containers, "/containers")
	if err != nil {
		return nil, err
	}
	res := []API.Container{}
	for _, con := range containers {
		var app API.App

		err := c.client.Get(&app, fmt.Sprintf("/apps/%s", con.AppID))
		if err != nil {
			return nil, err
		}
		if app.Name == appName {
			res = append(res, con)
		}
	}
	return res, nil
}

func (c *ArukasClient) createContainer(appName string, msg *message.IncomingMessage) error {

	var appSet API.AppSet

	// create an app
	newApp := API.App{Name: appName}

	ports := []API.Port{API.Port{Protocol: "tcp", Number: 80}}

	newContainer := API.Container{
		ImageName: msg.Repository.RepoName,
		Mem:       256,
		Instances: 1,
		Ports:     ports,
		//Cmd:       "",
		//Name:      "",
	}
	newAppSet := API.AppSet{
		App:       newApp,
		Container: newContainer,
	}

	// create
	if err := c.client.Post(&appSet, "/app-sets", newAppSet); err != nil {
		return err
	}

	// start container
	if err := c.client.Post(nil, fmt.Sprintf("/containers/%s/power", appSet.Container.ID), nil); err != nil {
		return err
	}

	return nil
}

func (c *ArukasClient) updateContainer(container API.Container, msg *message.IncomingMessage) error {

	container.ImageName = msg.Repository.RepoName
	// update
	if err := c.client.Patch(nil, fmt.Sprintf("/containers/%s", container.ID), container); err != nil {
		return err
	}
	// shutdown container
	if err := c.client.Delete(fmt.Sprintf("/containers/%s/power", container.ID)); err != nil {
		return err
	}

	if err := sleepUntilDown(c.client, container.ID, 300*time.Second); err != nil {
		return err
	}

	// start container
	if err := c.client.Post(nil, fmt.Sprintf("/containers/%s/power", container.ID), nil); err != nil {
		return err
	}
	return nil
}

func sleepUntilDown(client *API.Client, containerID string, timeout time.Duration) error {
	current := 0 * time.Second
	interval := 5 * time.Second
	for {
		var container API.Container
		if err := client.Get(&container, fmt.Sprintf("/containers/%s", containerID)); err != nil {
			return err
		}

		if container.StatusText == "stopped" {
			return nil
		}
		time.Sleep(interval)
		current += interval

		if timeout > 0 && current > timeout {
			return fmt.Errorf("Timeout: sleepUntilUp")
		}
	}
}
