package arukas

import (
	"fmt"
	API "github.com/arukasio/cli"
	"github.com/stretchr/testify/assert"
	"github.com/yamamoto-febc/arukas-ship/message"
	"testing"
)

func TestCreateContainer(t *testing.T) {
	msg := &message.IncomingMessage{}
	msg.Repository.RepoName = "nginx:latest"

	c, err := NewArukasClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)

	err = c.createContainer("testcontainer", msg)
	assert.NoError(t, err)

	containers, err := c.listContainerByName("testcontainer")
	assert.NoError(t, err)
	assert.NotNil(t, containers)
	assert.Len(t, containers, 1)

	err = c.updateContainer(containers[0], msg)
	assert.NoError(t, err)

	var container API.Container
	c.client.Get(&container, fmt.Sprintf("/containers/%s", containers[0].ID))
	c.client.Delete(fmt.Sprintf("/apps/%s", container.AppID))

}
