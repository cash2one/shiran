package masterslave

import (
	"github.com/williammuji/shiran/shiran"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"path/filepath"
	"os/exec"
	"os"
	"fmt"
	"strings"
	"bytes"
	"syscall"
	"time"
	"io/ioutil"
	"github.com/golang/glog"
)

type SlaveService struct {
	appManager		*AppManager
	name			string
}

func NewSlaveService(name string) *SlaveService {
	service := &SlaveService{
		appManager:		NewAppManager(),
		name:			name,
	}
	go service.appManager.eventLoop()
	return service
}

func (service *SlaveService) HandleAddApplicationRequest(msg *AddApplicationRequest, session *shiran.Session) {
	event := Event{
		eventType:		ADD,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleRemoveApplicationsRequest(msg *RemoveApplicationsRequest, session *shiran.Session) {
	event := Event{
		eventType:		REMOVE,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleStartApplicationsRequest(msg *StartApplicationsRequest, session *shiran.Session) {
	event := Event{
		eventType:		START,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleStopApplicationRequest(msg *StopApplicationRequest, session *shiran.Session) {
	event := Event{
		eventType:		STOP,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleRestartApplicationRequest(msg *RestartApplicationRequest, session *shiran.Session) {
	event := Event{
		eventType:		RESTART,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleListApplicationsRequest(msg *ListApplicationsRequest, session *shiran.Session) {
	event := Event{
		eventType:		LIST,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleGetApplicationsRequest(msg *GetApplicationsRequest, session *shiran.Session) {
	event := Event{
		eventType:		GET,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleGetHardwareRequest(msg *GetHardwareRequest, session *shiran.Session) {
	response := &GetHardwareResponse{
		SlaveCommander: msg.GetSlaveCommander(),
	}
	if msg.GetLshw() == true {
		out, err := exec.Command("lshw").CombinedOutput()
		if err != nil {
			glog.Errorf("HandleGetHardwareRequest lshw failed %s", err)
		} else {
			data := string(out)
			response.Lshw = &data
		}
	}

	out, err := exec.Command("lspci").CombinedOutput()
	if err != nil {
		glog.Errorf("HandleGetHardwareRequest lspci failed %s", err)
	} else {
		data := string(out)
		response.Lspci = &data
	}

	out, err = exec.Command("/sbin/ifconfig").CombinedOutput()
	if err != nil {
		glog.Errorf("HandleGetHardwareRequest ifconfig failed %s", err)
	} else {
		data := string(out)
		response.Ifconfig = &data
	}

	out, err = exec.Command("lscpu").CombinedOutput()
	if err != nil {
		glog.Errorf("HandleGetHardwareRequest lscpu failed %s", err)
	} else {
		data := string(out)
		response.Lscpu = &data
	}
	session.SendMessage("MasterService", "HandleGetHardwareResponse", response)
}

func (service *SlaveService) HandleGetFileContentRequest(msg *GetFileContentRequest, session *shiran.Session) {
	response := &GetFileContentResponse{
		SlaveCommander: msg.GetSlaveCommander(),
	}

	dir, err := os.Stat(msg.GetFileName())
	if err != nil {
		glog.Errorf("GetFileContentRequest stat %s failed %s", msg.GetFileName(), err)
		errorCode := fmt.Sprintf("%s", err)
		response.ErrorCode = &errorCode
	} else {
		if dir.IsDir() == true {
			errorCode := "Is Directory"
			response.ErrorCode = &errorCode
		} else if dir.Mode().IsRegular() == true {
			fileSize := int64(dir.Size())
			response.FileSize = &fileSize

			modifyTime := int64(dir.ModTime().Unix())
			response.ModifyTime = &modifyTime

			buffer, err := ioutil.ReadFile(msg.GetFileName())
			if err != nil {
				glog.Errorf("GetFileContentRequest ReadFile %s failed %s", msg.GetFileName(), err)
				errorCode := fmt.Sprintf("%s", err)
				response.ErrorCode = &errorCode
			} else {
				if msg.GetMaxSize() != 0 && msg.GetMaxSize() < dir.Size() {
					response.Content = buffer[:msg.GetMaxSize()]
				} else {
					response.Content = buffer
				}
				errorCode := "SUCCESS"
				response.ErrorCode = &errorCode
			}
		}
	}

	session.SendMessage("MasterService", "HandleGetFileContentResponse", response)
}

func (service *SlaveService) HandleGetFileChecksumRequest(msg *GetFileChecksumRequest, session *shiran.Session) {
	response := &GetFileChecksumResponse{
		SlaveCommander: msg.GetSlaveCommander(),
	}

	if len(msg.GetFiles()) != 0 {
		out, err := exec.Command("md5sum", msg.GetFiles()...).CombinedOutput()
		if err != nil {
			glog.Errorf("HandleGetFileChecksumRequest md5sum %v failed %s", msg.GetFiles(), err)
		} else {
			response.Md5Sums = make([]string, 0)
			res := strings.Split(string(out), "\n")
			for i := 0; i < len(res); i++ {
				response.Md5Sums = append(response.Md5Sums, res[i])
			}
		}
	}

	session.SendMessage("MasterService", "HandleGetFileChecksumResponse", response)
}

func (service *SlaveService) HandleRunCommandRequest(msg *RunCommandRequest, session *shiran.Session) {
	opened := uint64(0)
	fdPath := fmt.Sprintf("/proc/%d/fd", os.Getpid())
	err := filepath.Walk(fdPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			glog.Errorf("HandleRunCommandRequest Command %s %v filepath.Walk failed err:%s", msg.GetCommand(), msg.GetArgs(), err)
		} else {
			if fi.IsDir() == false {
				opened++
			}
		}
		return nil 
	})
	if err != nil {
		glog.Errorf("HandleRunCommandRequest Command %s %v filepath.Walk failed err:%s", msg.GetCommand(), msg.GetArgs(), err)
	}

	var rlimit syscall.Rlimit
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		glog.Errorf("HandleRunCommandRequest Command %s %v syscall.Getrlimit failed err:%s", msg.GetCommand(), msg.GetArgs(), err)
	}

	if rlimit.Cur < opened + 20 {
		errCode := "not enough fds"
		glog.Errorf("HandleRunCommandRequest start Command %s %v failed %s", msg.GetCommand(), msg.GetArgs(), errCode)
		response := &RunCommandResponse{
			SlaveCommander: msg.GetSlaveCommander(),
		}
		response.ErrorCode = &errCode
		session.SendMessage("MasterService", "HandleRunCommandResponse", response)
		return
	}

	go func(msg *RunCommandRequest, session *shiran.Session) {
		response := &RunCommandResponse{
			SlaveCommander: msg.GetSlaveCommander(),
		}

		startTime := time.Now().UnixNano()/1000

		cmd := exec.Command(msg.GetCommand(), msg.GetArgs()...)
		var stdOut, stdErr bytes.Buffer
		cmd.Stdout = &stdOut
		cmd.Stderr = &stdErr
		err := cmd.Start()
		if err != nil {
			glog.Errorf("HandleRunCommandRequest Start Command %s %v failed err:%s", msg.GetCommand(), msg.GetArgs(), err)
			errStr := fmt.Sprintf("%s", err)
			response.ErrorCode = &errStr
			session.SendMessage("MasterService", "HandleRunCommandResponse", response)
			return
		}

		var childStartTime uint64
		var ppid int64
		statPath := fmt.Sprintf("/proc/%d/stat", cmd.Process.Pid)
		processStat, err := linuxproc.ReadProcessStat(statPath)
		if err != nil {
			glog.Errorf("HandleRunCommandRequest ReadProcessStat %s %v failed err:%s", msg.GetCommand(), msg.GetArgs(), err)
			errStr := "cannot open /proc/pid/stat"
			response.ErrorCode = &errStr
		} else {
			childStartTime = processStat.Starttime
			ppid = processStat.Ppid
		}

		//timeout
		if msg.GetTimeout() != 0 {
			go func(timeout int32, childStartTime uint64, ppid int64, pid int) {
				for {
					select {
					case <-time.After(time.Duration(timeout) * time.Second):
						statPath := fmt.Sprintf("/proc/%d/stat", pid)
						processStat, err := linuxproc.ReadProcessStat(statPath)
						if err != nil {
							glog.Errorf("HandleRunCommandRequest timeout ReadProcessStat %s failed err:%s", statPath, err)
						} else {
							fmt.Printf("HandleRunCommandRequest timeout (%d,%d) (%d,%d) %d\n", childStartTime, processStat.Starttime, ppid, processStat.Ppid, pid)
							if childStartTime == processStat.Starttime && ppid == processStat.Ppid {
								syscall.Kill(pid, syscall.SIGTERM)
							}
						}
						return
					}
				}
			}(msg.GetTimeout(), childStartTime, ppid, cmd.Process.Pid)
		}

		path := fmt.Sprintf("/proc/%d/exe", cmd.Process.Pid)
		exeFile, err := os.Readlink(path)
		if err != nil {
			glog.Errorf("HandleRunCommandRequest Readlink %s %v failed err:%s", msg.GetCommand(), msg.GetArgs(), err)
			errStr := "cannot open /proc/pid/exe"
			response.ErrorCode = &errStr
		} else {
			response.ExecutableFile = &exeFile
		}

		err = cmd.Wait()
		if err != nil {
			glog.Errorf("HandleRunCommandRequest Wait %s %v failed err:%s", msg.GetCommand(), msg.GetArgs(), err)
			glog.Errorf("Command(%d) status:%s stdout:%s stderr:%s", cmd.ProcessState.Pid(), cmd.ProcessState, stdOut.String(), stdErr.String())
		}
		var pid int32 = int32(cmd.ProcessState.Pid())
		response.Pid = &pid
		status := fmt.Sprintf("%s", cmd.ProcessState)
		response.Status = &status

		if msg.GetMaxStdout() != 0 && msg.GetMaxStdout() < int32(len(stdOut.Bytes())) {
			response.StdOutput = stdOut.Bytes()[:msg.GetMaxStdout()]
		} else {
			response.StdOutput = stdOut.Bytes()
		}
		if msg.GetMaxStderr() != 0 && msg.GetMaxStderr() < int32(len(stdErr.Bytes())) {
			response.StdError = stdErr.Bytes()[:msg.GetMaxStderr()]
		} else {
			response.StdError = stdErr.Bytes()
		}

		response.StartTimeUs = &startTime
		endTime := time.Now().UnixNano()/1000
		response.FinishTimeUs = &endTime

		su := cmd.ProcessState.SysUsage()
		if sysUsage, ok := su.(*syscall.Rusage); ok {
			response.UserTime = &sysUsage.Utime.Usec
			response.SystemTime = &sysUsage.Stime.Usec
			response.MemoryMaxrssKb = &sysUsage.Maxrss
		}

		waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
		es := int32(waitStatus.ExitStatus())
		response.ExitStatus = &es
		sg := int32(waitStatus.Signal())
		response.Signaled = &sg
		cd := waitStatus.CoreDump()
		response.Coredump = &cd

		if response.ErrorCode == nil {
			errStr := "SUCCESS"
			response.ErrorCode = &errStr
		}

		session.SendMessage("MasterService", "HandleRunCommandResponse", response)

	}(msg, session)
}

func (service *SlaveService) HandleRunScriptRequest(msg *RunScriptRequest, session *shiran.Session) {
	go func() {
		prefix := fmt.Sprintf("%s-runScript-%d-%d-", service.name, time.Now().Unix(), os.Getpid())
		f, err := ioutil.TempFile("/tmp", prefix)
		if err != nil {
			glog.Errorf("HandleRunScriptRequest %v failed err:%s", msg, err)
		}

		os.Chmod(f.Name(), 0755)
		f.Write(msg.GetScript())
		f.Close()

		request := &RunCommandRequest{
			SlaveCommander: msg.GetSlaveCommander(),
		}
		fname := f.Name()
		request.Command = &fname
		request.MaxStdout = msg.MaxStdout
		request.MaxStderr = msg.MaxStderr
		request.Timeout = msg.Timeout
		service.HandleRunCommandRequest(request, session)
	}()
}
