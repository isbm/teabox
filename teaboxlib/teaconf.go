package teaboxlib

import (
	"fmt"
	"os"
	"path"

	"github.com/davecgh/go-spew/spew"
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
	title        string
	contentPath  string
	initConfPath string

	modIndex []TeaConfComponent
}

func NewTeaConf(appname string) *TeaConf {
	tc := new(TeaConf)
	tc.modIndex = []TeaConfComponent{}

	tc.contentPath = nanoconf.NewConfig(appname+".conf").Root().String("content", "")
	tc.initConfPath = path.Join(tc.contentPath, "init.conf")
	if err := tc.initConfig(); err != nil {
		panic(fmt.Sprintf("Unable to initialise modules: %s", err))
	}

	return tc
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
				c := nanoconf.NewConfig(pth)
				groupId := c.Root().String("group", "")
				if groupId != "" {
					group := tc.getGroup(groupId)
					group.Add(NewTeaConfModule(c.Root().String("title", "")))
				} else {
					fmt.Printf("%s %s\n", de.ModeType(), pth)
				}
			}
			return nil
		},
		Unsorted: true,
	})

	spew.Dump(tc.modIndex)
	return err
}
