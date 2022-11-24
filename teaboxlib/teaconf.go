package teaboxlib

import (
	"fmt"
	"os"
	"path"
	"sort"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
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

	modIndex []TeaConfComponent
	wzlib_logger.WzLogger
}

func NewTeaConf(appname string) *TeaConf {
	tc := new(TeaConf)
	tc.modIndex = []TeaConfComponent{}

	cfg := nanoconf.NewConfig(appname + ".conf")

	tc.contentPath = cfg.Root().String("content", "")
	tc.callbackSocketPath = cfg.Root().String("callback", "")
	tc.initConfPath = path.Join(tc.contentPath, "init.conf")

	if err := tc.initConfig(); err != nil {
		panic(fmt.Sprintf("Unable to initialise modules: %s", err))
	}

	return tc
}

func (tc *TeaConf) GetTitle() string {
	return tc.title
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

func (tc *TeaConf) initConfig() error {
	for _, p := range []string{tc.initConfPath, tc.contentPath} {
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			return err
		}
	}

	initconf := nanoconf.NewConfig(tc.initConfPath)
	tc.title = initconf.Root().String("title", "")

	err := godirwalk.Walk(tc.contentPath, &godirwalk.Options{
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
						SetCallbackPath(tc.callbackSocketPath).
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
