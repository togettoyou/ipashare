package sign

import (
	"os"
	"path"
	"supersign/pkg/conf"
	"sync"
	"time"

	"go.uber.org/zap"
)

var signJob *job

type job struct {
	logger    *zap.Logger
	streamCh  chan *Stream
	doneCache map[string]*done
	mu        sync.Mutex
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
}

func Setup(logger *zap.Logger, maxJob int) {
	signJob = &job{
		logger:    logger,
		streamCh:  make(chan *Stream, maxJob),
		doneCache: make(map[string]*done, 0),
	}
	go func() {
		for {
			select {
			case stream, ok := <-signJob.streamCh:
				if !ok {
					return
				}
				go func() {
					time.Sleep(1 * time.Hour)
					signJob.mu.Lock()
					defer signJob.mu.Unlock()
					signJob.logger.Info("开始清理旧数据:" + stream.ProfileUUID)
					delete(signJob.doneCache, stream.ProfileUUID)
					os.RemoveAll(path.Join(conf.Apple.TemporaryFilePath, stream.ProfileUUID))
				}()
				go func() {
					signJob.logger.Info("开始执行重签名任务......")
					err := run(
						path.Join(conf.Apple.AppleDeveloperPath, stream.Iss, "pem.pem"),
						path.Join(conf.Apple.AppleDeveloperPath, stream.Iss, "key.key"),
						stream.MobileprovisionPath,
						path.Join(conf.Apple.TemporaryFilePath, stream.ProfileUUID, "ipa.ipa"),
						path.Join(conf.Apple.UploadFilePath, stream.IpaUUID, "ipa.ipa"),
					)
					if err != nil {
						signJob.mu.Lock()
						defer signJob.mu.Unlock()
						signJob.doneCache[stream.ProfileUUID] = &done{
							Success: false,
							Msg:     "重签名任务执行失败:" + err.Error(),
						}
						signJob.logger.Error("重签名任务执行失败:" + err.Error())
						return
					}
					signJob.mu.Lock()
					defer signJob.mu.Unlock()
					signJob.doneCache[stream.ProfileUUID] = &done{
						Success:          true,
						Msg:              "重签名任务执行成功",
						BundleIdentifier: stream.BundleIdentifier,
						Version:          stream.Version,
						Name:             stream.Name,
						Summary:          stream.Summary,
						IpaUUID:          stream.IpaUUID,
					}
					signJob.logger.Info("重签名任务执行成功")
				}()
			}
		}
	}()
}

func Push(stream *Stream) {
	signJob.streamCh <- stream
}

func DoneCache(ProfileUUID string) (done *done, ok bool) {
	signJob.mu.Lock()
	defer signJob.mu.Unlock()
	done, ok = signJob.doneCache[ProfileUUID]
	return
}
