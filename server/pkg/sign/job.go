package sign

import (
	"fmt"
	"ipashare/pkg/ali"
	"ipashare/pkg/conf"
	"os"
	"path"
	"sync"
	"time"

	"go.uber.org/zap"
)

var signJob *job

type job struct {
	logger    *zap.Logger
	streamCh  chan struct{}
	doneCache map[string]*done
	mu        sync.RWMutex
}

type Stream struct {
	ProfileUUID         string
	Iss                 string
	MobileprovisionPath string
	IpaUUID             string
	BundleIdentifier    string
	Version             string
	Name                string
	Summary             string
}

type done struct {
	Success          bool
	Msg              string
	BundleIdentifier string
	Version          string
	Name             string
	Summary          string
	IpaUUID          string
	IpaURL           string
}

func Setup(logger *zap.Logger, maxJob int) {
	signJob = &job{
		logger:    logger,
		streamCh:  make(chan struct{}, maxJob),
		doneCache: make(map[string]*done, 0),
	}
}

func Push(stream *Stream) {
	signJob.streamCh <- struct{}{}
	go func() {
		cleanOldData := 1
		if conf.Server.CleanOldData > 1 {
			cleanOldData = conf.Server.CleanOldData
		}
		time.Sleep(time.Duration(cleanOldData) * time.Hour)
		signJob.mu.Lock()
		delete(signJob.doneCache, stream.ProfileUUID)
		signJob.mu.Unlock()
		signJob.logger.Info("开始清理旧数据:" + stream.ProfileUUID)
		os.RemoveAll(path.Join(conf.Apple.TemporaryFilePath, stream.ProfileUUID))
		ali.DelFile(stream.ProfileUUID + ".ipa")
	}()
	go func() {
		var err error
		defer func() {
			if e := recover(); e != nil {
				signJob.logger.Named("Push").Error(fmt.Sprintf("%v", e))
			}
			if err != nil {
				signJob.mu.Lock()
				signJob.doneCache[stream.ProfileUUID] = &done{
					Success: false,
					Msg:     "打包任务执行失败:" + err.Error(),
				}
				signJob.mu.Unlock()
				signJob.logger.Error("打包任务执行失败:" + err.Error())
			}
			<-signJob.streamCh
		}()
		signJob.logger.Info("开始执行打包任务......")
		err = run(
			path.Join(conf.Apple.AppleDeveloperPath, stream.Iss, "pem.pem"),
			path.Join(conf.Apple.AppleDeveloperPath, stream.Iss, "key.key"),
			stream.MobileprovisionPath,
			path.Join(conf.Apple.TemporaryFilePath, stream.ProfileUUID, "ipa.ipa"),
			path.Join(conf.Apple.UploadFilePath, stream.IpaUUID, "ipa.ipa"),
		)
		if err != nil {
			return
		}
		ipaURL, _ := ali.UploadFile(
			stream.ProfileUUID+".ipa",
			path.Join(conf.Apple.TemporaryFilePath, stream.ProfileUUID, "ipa.ipa"),
		)
		if ipaURL == "" {
			ipaURL = fmt.Sprintf("%s/api/v1/download/tempipa/%s", conf.Server.URL, stream.ProfileUUID)
		}
		signJob.mu.Lock()
		signJob.doneCache[stream.ProfileUUID] = &done{
			Success:          true,
			Msg:              "打包任务执行成功",
			BundleIdentifier: stream.BundleIdentifier,
			Version:          stream.Version,
			Name:             stream.Name,
			Summary:          stream.Summary,
			IpaUUID:          stream.IpaUUID,
			IpaURL:           ipaURL,
		}
		signJob.mu.Unlock()
		signJob.logger.Info("打包任务执行成功")
	}()
}

func DoneCache(ProfileUUID string) (done *done, ok bool) {
	signJob.mu.RLock()
	defer signJob.mu.RUnlock()
	done, ok = signJob.doneCache[ProfileUUID]
	return
}
