package manifest

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

type Service struct {
	Name string

	Build       ServiceBuild
	Environment []string
	Image       string
	Test        string
	Volumes     []string
}

type Services []Service

type ServiceBuild struct {
	Args []string
	Path string
}

func (s Service) BuildHash() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("build[path=%q, args=%v] image=%q", s.Build.Path, s.Build.Args, s.Image))))
}

func (s Service) Env(env []string) (map[string]string, error) {
	full := map[string]string{}

	penv := parseEnv(env)

	for _, e := range s.Environment {
		parts := strings.SplitN(e, "=", 2)

		switch len(parts) {
		case 1:
			v, ok := penv[parts[0]]
			if !ok {
				return nil, fmt.Errorf("required env: %s", parts[0])
			}
			full[parts[0]] = v
		case 2:
			v, ok := penv[parts[0]]
			if ok {
				full[parts[0]] = v
			} else {
				full[parts[0]] = parts[1]
			}
		default:
			return nil, fmt.Errorf("invalid environment")
		}
	}

	return full, nil
}

func parseEnv(env []string) map[string]string {
	parsed := map[string]string{}

	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)

		if len(parts) == 2 {
			parsed[parts[0]] = parts[1]
		}
	}

	return parsed
}

func (ss Services) Find(name string) (*Service, error) {
	for _, s := range ss {
		if s.Name == name {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("could not find service: %s", name)
}
