package work

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-logr/logr"

	c "github.com/ghostbaby/sentinel/pkg/client"
	"github.com/ghostbaby/sentinel/pkg/g"
	"github.com/ghostbaby/sentinel/pkg/models"
	ctrl "sigs.k8s.io/controller-runtime"
)

func PingCn() []*models.WorkResult {

	var list []*models.WorkResult

	config := g.Config()
	log := ctrl.Log.WithName("controllers").WithName("sentinel").WithName("ping.cn")

	chJobs := make(chan *models.Job, len(config.IPS))
	chResults := make(chan *models.ScopeResult, len(config.IPS))

	for w := 1; w <= config.Routines; w++ {
		go PinCnWork(chJobs, chResults)
	}

	for ip, _ := range config.IPS {
		job := &models.Job{
			Log:        log,
			IP:         ip,
			RetryTimes: config.RetryTimes,
		}
		chJobs <- job
	}
	close(chJobs)

	for i := 1; i <= len(config.IPS); i++ {
		res := <-chResults
		if res.Work == nil {
			log.Error(errors.New("result is nil"), fmt.Sprintf(" failed to get result for %s", res.IP))
			continue
		}
		list = append(list, res.Work)
	}

	return list
}

func PinCnWork(jobs <-chan *models.Job, results chan<- *models.ScopeResult) {
	for j := range jobs {

		var res models.ScopeResult
		rt := j.RetryTimes

		for {
			taskID, err := PinCnTask(j.Log, j.IP)
			if err != nil {
				j.Log.Error(err, "fail to exec PinCnTask")
				continue
			}

			time.Sleep(20 * time.Second)

			result, err := PinCnResult(j.Log, j.IP, taskID)
			res.Work = result
			res.IP = j.IP
			if err != nil {
				j.Log.Error(err, "fail to exec PinCnTask")
				rt--
				if rt == 0 {
					results <- &res
					break
				}
				continue
			}
			results <- &res
			break
		}

	}

}

func PinCnTask(log logr.Logger, ip string) (string, error) {
	client := c.NewBaseClient(models.PingCnUrl, 0)
	w := NewWork(client, log)
	log.Info("start to get task id", "url", models.PingCnUrl, "host", ip)

	task, err := w.GetPingCnTaskInfo(models.PingCnCheckPath, ip)
	if err != nil {
		msg := fmt.Sprintf("fail to get task id for %s", ip)
		log.Error(err, msg)
		return "", err
	}

	if task.Code != 1 {
		log.Error(errors.New("fail to create task"), "return wrong value task.code ",
			"host", ip, "code", task.Code)
		return "", err
	}

	log.Info("success to get task id", "url", models.PingCnUrl, "host", ip, "task", task.Data.TaskID)
	return task.Data.TaskID, nil
}

func PinCnResult(log logr.Logger, ip, taskID string) (*models.WorkResult, error) {

	config := g.Config()
	client := c.NewBaseClient(models.PingCnUrl, 0)
	w := NewWork(client, log)
	log.Info("start to get result", "host", ip, "task", taskID)

	result, err := w.GetPingCnResult(models.PingCnCheckPath, ip, taskID)
	if err != nil {
		msg := fmt.Sprintf("fail to get result for %s", ip)
		log.Error(err, msg, "host", ip, "task", taskID)
		return nil, err
	}

	if result.Code != 1 {
		log.Error(errors.New("fail to get result"), "return wrong value task.code ",
			"host", ip, "code", result.Code)
		return nil, err
	}

	if result.Data == nil {
		log.Error(errors.New("fail to get result"), "result.Data is nil ",
			"host", ip, "task", taskID)
		return nil, err
	}

	log.Info("success to get result", "url", models.PingCnUrl, "host", ip, "task", taskID)

	ipd, ok := config.IPS[ip]
	if !ok {
		log.Error(errors.New("fail to get link name"), "LineNameMap is not exist ",
			"host", ip, "task", taskID)
	}

	loss, details := PingCnLost(result.Data.InitData.Result)

	return &models.WorkResult{
		Provider: models.PingCnUrl,
		Name:     ipd.Describe,
		Host:     ip,
		Max:      result.Data.InitData.MinMaxAvg.Max.Cost,
		Min:      result.Data.InitData.MinMaxAvg.Min.Cost,
		Avg:      result.Data.InitData.MinMaxAvg.Avg.Cost,
		Loss:     loss,
		Details:  details,
	}, nil
}

func PingCnLost(result []*models.PingCnResultInfo) (float64, []*models.ResultDetails) {
	var (
		packets, received int
		details           []*models.ResultDetails
	)

	for _, info := range result {
		packets += info.Packets
		received += info.Received

		p, _ := strconv.ParseFloat(strconv.Itoa(info.Packets), 64)
		r, _ := strconv.ParseFloat(strconv.Itoa(info.Received), 64)

		d := &models.ResultDetails{
			Name:    info.NodeName,
			Loss:    1 - (r / p),
			Max:     info.Max,
			Min:     info.Min,
			Avg:     info.Avg,
			Area:    info.Area,
			IspName: info.IspName,
		}

		details = append(details, d)
	}

	p, _ := strconv.ParseFloat(strconv.Itoa(packets), 64)
	r, _ := strconv.ParseFloat(strconv.Itoa(received), 64)

	return 1 - (r / p), details
}
