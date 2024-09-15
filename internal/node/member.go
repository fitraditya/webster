package node

import (
	"context"
	"strings"

	"github.com/fitraditya/webster/config"
	"github.com/hashicorp/memberlist"
	"github.com/obrel/go-lib/pkg/log"
)

func CreateMemberList(ctx context.Context, delegate *Delegate, port int, join string) (*memberlist.Memberlist, error) {
	c := memberlist.DefaultLocalConfig()
	c.Name = config.GetNodeName()
	c.BindPort = port
	c.Delegate = delegate
	c.Logger.SetOutput(config.GetLogLevel())

	list, err := memberlist.Create(c)

	if err != nil {
		return nil, err
	}

	if join != "" {
		hosts := strings.Split(join, ",")

		if len(hosts) > 0 {
			log.For("member", "create").Infof("%v joining cluster", hosts)

			if _, err := list.Join(hosts); err != nil {
				return nil, err
			}
		}
	}

	return list, nil
}
