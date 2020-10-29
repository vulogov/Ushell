package piping

import (
  "fmt"
  "bytes"
  "github.com/vulogov/TSAK/internal/log"
  // "github.com/vulogov/TSAK/internal/conf"
  zmq "github.com/pebbe/zmq4"
)

var INCH = 0
var OUTCH = 1

var pipeIn = make(chan string, 1000000)
var pipeOut = make(chan string, 1000000)
var zmqS = make(map[string]*zmq.Socket)

var zmqCtx,_ = zmq.NewContext()
var zmqErr int64

func To(dst int, _data []byte) {
  var data = bytes.NewBuffer(_data)
  log.Trace(fmt.Sprintf("Sending %d bytes to  pipeline %d", data.Len(), dst))
  if dst == INCH {
    pipeIn <- data.String()
    log.Trace(fmt.Sprintf("%d element in pipeline IN", len(pipeIn)))
  } else if dst == OUTCH {
    pipeOut <- data.String()
    log.Trace(fmt.Sprintf("%d element in pipeline OUT", len(pipeOut)))
  } else {
    log.Error("Trying to send data to non-existent pipeline")
  }
}

func From(src int) []byte {
  if src == INCH && len(pipeIn) > 0 {
    return []byte(<-pipeIn)
  } else if src == OUTCH && len(pipeOut) > 0 {
    return []byte(<-pipeOut)
  } else {
    return []byte("")
  }
}

func Len(src int) int {
  if src == INCH {
    return len(pipeIn)
  } else if src == OUTCH {
    return len(pipeOut)
  } else {
    return 0
  }
}

func Shutdown() {
  log.Trace("Terminating Pipelines")
  for k, v := range zmqS {
    log.Trace(fmt.Sprintf("Closing ZMQ: %s", k))
    v.Close()
  }
  zmqCtx.Term()
  log.Trace("Pipelines are terminated")
}