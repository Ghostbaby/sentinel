package work

import (
	"context"

	"github.com/ghostbaby/sentinel/pkg/models"

	"github.com/pkg/errors"

	"github.com/ghostbaby/sentinel/pkg/client"
	"github.com/go-logr/logr"
)

type Work struct {
	*client.BaseClient
	CTX context.Context
	Log logr.Logger
}

func NewWork(cli *client.BaseClient, log logr.Logger) *Work {
	return &Work{
		BaseClient: cli,
		CTX:        context.Background(),
		Log:        log,
	}
}

func (w *Work) RebuildCTX() *Work {
	w.CTX = context.Background()
	return w
}

func (w *Work) GetPingCnTaskInfo(path, ip string) (*models.PingCnTask, error) {
	data := models.PingCnPayload{
		Host:       ip,
		Type:       models.PingCnPayloadType,
		CreateTask: models.PingCnPayloadCreateTask,
	}

	list := &models.PingCnTask{}
	if err := w.Post(w.CTX, path, &data, &list); err != nil {
		return nil, errors.Wrap(err, "fail to call /check api to get task info")
	}

	return list, nil
}

func (w *Work) GetPingCnResult(path, ip, taskID string) (*models.PingCnResult, error) {
	data := models.PingCnPayload{
		Host:       ip,
		Type:       models.PingCnPayloadType,
		CreateTask: models.PingCnPayloadExecTask,
		TaskID:     taskID,
	}

	list := &models.PingCnResult{}
	if err := w.Post(w.CTX, path, &data, &list); err != nil {
		return nil, errors.Wrap(err, "fail to call /check api to get result")
	}

	return list, nil
}
