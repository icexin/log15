/**
* Time based log rotation Writer.
* Could rotate log file every `rotateInterval` minutes.
**/

package log

import (
	"io"
	"os"
	"sync"
	"time"
)

type TimeRotateWriter struct{
	filename 		string
	maxBackups		int		// max log files
	rotateInterval	int		// in minutes

	file 			io.WriteCloser
	mutex			sync.Mutex
	rotateAt 		int64
	intervalInSeconds	int64
}

func NewTimeRotateWriter(filename string, interval int, backupCount int) (*TimeRotateWriter, error) {
	wr := TimeRotateWriter{
		filename: 		filename,
		maxBackups:		backupCount,
		rotateInterval:	interval,
	}

	// init rotate time
	wr.intervalInSeconds = int64(interval * 60)
	wr.calcNextRotateTime()

	// open file to write
	err := wr.openFile();
	return &wr, err
}

// implements Write interface of io.Writer
func (wr *TimeRotateWriter) Write(data []byte) (succBytes int, err error){
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	// Open log file
	if err := wr.openFile(); err != nil{
		return 0, err
	}

	if wr.shouldRotate(){
		if err := wr.rotate(); err != nil{
			return 0, err
		}
	}

	return wr.file.Write(data)
}

// Close of WriterCloser
func (wr *TimeRotateWriter) Close() (err error) {
	if err = wr.file.Close(); err != nil {
		return
	}
	wr.file = nil
	return
}

func (wr *TimeRotateWriter) openFile() error{
	if wr.file == nil{
		fd, err := os.OpenFile(wr.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil{
			return err
		}
		wr.file = fd
	}
	return nil
}

// check if should rotate the log
func (wr *TimeRotateWriter) shouldRotate() bool {
	return time.Now().Unix() >= wr.rotateAt
}

// calculate the next log rotation time
func (wr *TimeRotateWriter) calcNextRotateTime() {
	currentTime := time.Now().Unix()

	timestruct := time.Unix(currentTime, 0)
	//currentHour := timestruct.Hour()
	//currentMinute := timestruct.Minute()
	currentSecond := timestruct.Second()

	wr.rotateAt = int64(currentTime - int64(currentSecond) + wr.intervalInSeconds)
}

// do log rotation
func (wr *TimeRotateWriter) rotate() (err error) {
	if err = wr.Close(); err != nil{
		return err
	}

	dstTime := wr.rotateAt - wr.intervalInSeconds
	dstPath := wr.filename + "." + time.Unix(dstTime, 0).Format("200601021504")

	if _, err := os.Stat(dstPath); err == nil {
		os.Remove(dstPath)
	}

	if err = os.Rename(wr.filename, dstPath); err != nil{
		return err
	}

	if wr.maxBackups > 0 {
		wr.deleteExpiredFiles()
	}

	wr.calcNextRotateTime()
	
	err = wr.openFile()
	return err
}

// delete expired log files
func (wr *TimeRotateWriter) deleteExpiredFiles() {
	// TODO: implement deleting expired files 
}