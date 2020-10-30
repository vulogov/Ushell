package tsak

import (
  "sync"
  "github.com/sirupsen/logrus"
  "github.com/vulogov/TSAK/internal/log"
  "github.com/vulogov/TSAK/internal/script"
  "github.com/vulogov/TSAK/internal/signal"
  "github.com/vulogov/TSAK/internal/nr"
  "github.com/vulogov/TSAK/internal/conf"
)

func InProc() {
  var start = nr.NowMillisec()
  go func(wg *sync.WaitGroup) {
    defer wg.Done()
    inproc()
    log.Trace("In thread exiting")
    signal.ExitRequest()
    nr.RecordDuration("In() duration", start)
  }(signal.WG())
}

func inproc() {
  if conf.Conf != "" {
    res := script.RunScript("in", conf.Conf)
    log.Event(
      "Bootstrap is loaded for In()",
      logrus.Fields{
        "result":     res,
        "confSource": conf.Conf,
    })
  }
  script.RunScript("in", conf.In)
}
