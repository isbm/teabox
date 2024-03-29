package teaboxlib

import (
	"fmt"
	"os"
	"path"
	"sort"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_utils "github.com/infra-whizz/wzlib/utils"
	"github.com/isbm/go-nanoconf"
	"github.com/karrick/godirwalk"
)

/*
TeaConf is a global config of everything.
1. Application is reading a startup config (either /etc/<appname>.conf or ~/<appname>rc or ./appname.conf)
2. There is a "content" directive where to recursively parse everything else, starting from "init.conf" configuration
3. Each module is in its own directory and also has "init.conf"
4. All this information is then put into a common tree
*/

type TeaConf struct {
	title              string
	contentPath        string
	initConfPath       string
	callbackSocketPath string
	rootConf           *nanoconf.Config

	modIndex []TeaConfComponent
	wzlib_logger.WzLogger
}

func NewTeaConf(appname string) (*TeaConf, error) {
	tc := new(TeaConf)

	configFileName := fmt.Sprintf("%s.conf", appname)
	configPath := configFileName
	if !wzlib_utils.FileExists(configPath) {
		configPath = path.Join("/etc", configFileName)
		if !wzlib_utils.FileExists(configPath) {
			return nil, fmt.Errorf("no config file found as \"./%[1]s\" neither as \"/etc/%[1]s\"", configFileName)
		}
	}

	tc.rootConf = nanoconf.NewConfig(configPath)
	tc.contentPath = tc.GetRootConfig().Root().String("content", "")
	tc.callbackSocketPath = tc.GetRootConfig().Root().String("callback", "")
	tc.initConfPath = path.Join(tc.contentPath, "init.conf")

	environ, exists := tc.GetRootConfig().Root().Raw()["env"]
	if exists {
		for envKey, envVar := range environ.(map[interface{}]interface{}) {
			os.Setenv(fmt.Sprintf("%v", envKey), fmt.Sprintf("%v", envVar))
		}
	}

	return tc, nil
}

// GetRootConfig returns a root configuration of the application
func (tc *TeaConf) GetRootConfig() *nanoconf.Config {
	return tc.rootConf
}

// GetContentPath returns root path to all modules
func (tc *TeaConf) GetContentPath() string {
	return tc.contentPath
}

func (tc *TeaConf) GetTitle() string {
	return tc.title
}

func (tc *TeaConf) GetSocketPath() string {
	return tc.callbackSocketPath
}

func (tc *TeaConf) GetModuleStructure() []TeaConfComponent {
	return tc.modIndex
}

// Returns a group by its ID
func (tc *TeaConf) getGroup(id string) TeaConfComponent {
	for _, c := range tc.modIndex {
		if c.GetGroup() == id {
			return c
		}
	}

	c := NewTeaConfGroup(id)
	tc.modIndex = append(tc.modIndex, c)
	return c
}

func (tc *TeaConf) InitConfig() error {
	if tc.modIndex != nil {
		return fmt.Errorf("config already initalised")
	}

	tc.modIndex = []TeaConfComponent{}

	for _, p := range []string{tc.initConfPath, tc.contentPath} {
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			return err
		}
	}

	initconf := nanoconf.NewConfig(tc.initConfPath)
	tc.title = initconf.Root().String("title", "")

	err := godirwalk.Walk(tc.contentPath, &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback: func(pth string, de *godirwalk.Dirent) error {
			if pth == tc.initConfPath {
				return nil
			}

			if path.Base(pth) == "init.conf" {
				defer func() {
					if err := recover(); err != nil {
						tc.GetLogger().Errorf("Error loading %s: %s", pth, err)
						panic(err)
					}
				}()
				c := nanoconf.NewConfig(pth)
				title := c.Root().String("title", "")
				if title == "" {
					fmt.Printf("Skipping module: %s (insufficient configuration)\n", pth)
					return nil
				}

				m := NewTeaConfModule(title)
				m.modulePath = path.Dir(pth)
				if c.Root().Raw()["commands"] != nil {
					m.SetCondition(c.Root().Raw()["conditions"]).
						SetCommands(c.Root().Raw()["commands"]).
						SetCallbackPath(tc.GetSocketPath()).
						SetLandingPageType(c.Root().String("landing", "")).
						SetSetupCommand(c.Root().String("setup", ""))
				}

				groupId := c.Root().String("group", "")
				if groupId != "" { // Belongs to group and IS a group. Group has no commands, otherwise it is a module that only belongs to a group
					tc.getGroup(groupId).Add(m)
				} else {
					tc.modIndex = append(tc.modIndex, m)
				}
			}
			return nil
		},
		Unsorted: true,
	})

	// Sort all the items, except "Exit", which should be always at the end.
	sort.Slice(tc.modIndex, func(i, j int) bool {
		return tc.modIndex[i].GetTitle() < tc.modIndex[j].GetTitle()
	})

	tc.modIndex = append(tc.modIndex, NewTeaConfCmd("exit", LABEL_EXIT))

	return err
}
