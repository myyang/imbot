package commands

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/subcommands"

	botHttp "github.com/myyang/imbot/http"
	botLog "github.com/myyang/imbot/log"
)

// FIXME: change host
const jenkinsURL = "https://your.jenkins.com/buildByToken/buildWithParameters"

type jenkinsCmd struct {
	logger botLog.Logger

	src string
	env string
}

func (c *jenkinsCmd) Name() string     { return "jenkins" }
func (c *jenkinsCmd) Synopsis() string { return "trigger jenkins jobs" }
func (c *jenkinsCmd) Usage() string {
	return `jenkins <flags> <job>:
		trigger jenkins jobs
	-src string
		source branch/tag, ex: master/release/v1.0
	-env string
		target env, ex: prod/stage/dev
	===
	available jobs:
` + c.listJobs()
}

func (c *jenkinsCmd) listJobs() string {
	b := &strings.Builder{}
	// list jobs
	return b.String()
}

func (c *jenkinsCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.src, "src", "", "source branch/tage, ex: master/release/v1.0")
	f.StringVar(&c.env, "env", "", "target env(s), separated with ',', ex: prod/stage/dev,dev2")
}

func (c *jenkinsCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	c.logger = ctx.Value(botLog.CtxLogger).(botLog.Logger)

	if f.NArg() == 0 || f.Arg(0) == "help" {
		c.logger.Infof("%v", c.Usage())
		return subcommands.ExitSuccess
	}

	job := f.Arg(0)
	values := url.Values{
		"job":   []string{job},
		"cause": []string{"triggered from imbot"},
		// FIXME: following params are for parametered triggers
		"src":   []string{c.src},
		"token": []string{fmt.Sprintf("%x", md5.Sum([]byte(job)))},
	}

	envs := strings.Split(c.env, ",")

	result := &strings.Builder{}
	for i, e := range envs {
		if i > 0 {
			result.WriteString("===\n")
		}

		// FIXME: src is for parametered triggers
		values["env"] = []string{e}
		targetURL := jenkinsURL + "?" + values.Encode()

		httpRequest, err := http.NewRequest(http.MethodGet, targetURL, nil)
		if err != nil {
			result.WriteString(fmt.Sprintf(
				"job: %v, env: %v, src: %v\n\nprepare request error: %v\n",
				job, e, c.src, err,
			))
			continue
		}

		// FIXME: set auth if required
		// httpRequest.SetBasicAuth(user, password)

		raw, httpResp, err := botHttp.SendRequestDebug(httpRequest)
		if err != nil {
			result.WriteString(fmt.Sprintf(
				"job: %v, env: %v, src: %v\n\nrequest error: %v\nresponse: %s\n",
				job, e, c.src, err, raw,
			))
			continue
		}

		switch httpResp.StatusCode {
		case http.StatusCreated:
			result.WriteString(fmt.Sprintf(
				"job: %v, env: %v, src: %v\n\nstatus: triggered\n",
				job, e, c.src,
			))
		default:
			result.WriteString(fmt.Sprintf(
				"job: %v, env: %v, src: %v\n\nstsatus: failed, code: %v\n",
				job, e, c.src, httpResp.StatusCode,
			))
		}
	}
	c.logger.Infof(result.String())

	return subcommands.ExitSuccess
}
