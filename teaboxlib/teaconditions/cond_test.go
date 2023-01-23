package teaconditions

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TeaCondTestSuite struct {
	conditions []map[string][]string
	suite.Suite
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TeaCondTestSuite))
}

func (suite *TeaCondTestSuite) SetupTest() {
	suite.conditions = []map[string][]string{}
}

func (suite *TeaCondTestSuite) TestPresentBlocks() {
	present_condition := map[string][]string{
		"present": {
			"/bin/who-knows-what",
		},
		"message": {"no such file"},
	}

	suite.conditions = append(suite.conditions, present_condition)
	p, e := NewTeaConditionsProcessor(suite.conditions)

	suite.Equal(nil, e)
	suite.False(false, p.Satisfied())
	suite.Equal("no such file", p.GetInfoMessage())
}

func (suite *TeaCondTestSuite) TestPresentPasses() {
	fn := "/tmp/a.txt"
	ioutil.WriteFile(fn, []byte("test"), 0600)
	defer os.Remove(fn)

	present_condition := map[string][]string{
		"present": {fn},
		"message": {"no such file"},
	}

	suite.conditions = append(suite.conditions, present_condition)
	p, e := NewTeaConditionsProcessor(suite.conditions)

	suite.Equal(nil, e)
	suite.True(true, p.Satisfied())
	suite.Equal("", p.GetInfoMessage())
}

func (suite *TeaCondTestSuite) TestAbsentBlocks() {
	fn := "/tmp/b.txt"
	ioutil.WriteFile(fn, []byte("test"), 0600)
	defer os.Remove(fn)

	present_condition := map[string][]string{
		"absent":  {fn},
		"message": {"file present"},
	}

	suite.conditions = append(suite.conditions, present_condition)
	p, e := NewTeaConditionsProcessor(suite.conditions)

	suite.Equal(nil, e)
	suite.False(false, p.Satisfied())
	suite.Equal("file present", p.GetInfoMessage())
}

func (suite *TeaCondTestSuite) TestAbsentPasses() {
	present_condition := map[string][]string{
		"absent":  {"/dev/darth-wader"},
		"message": {"file present"},
	}

	suite.conditions = append(suite.conditions, present_condition)
	p, e := NewTeaConditionsProcessor(suite.conditions)

	suite.Equal(nil, e)
	suite.True(true, p.Satisfied())
	suite.Equal("", p.GetInfoMessage())
}

func (suite *TeaCondTestSuite) TestPresentCrAbsentCr() {
	/*
		present: /tmp/present.txt (created)
		absent: /tmp/absent.txt (created)
	*/
	fna := "/tmp/present.txt"
	fnb := "/tmp/absent.txt"

	for _, fn := range []string{fna, fnb} {
		ioutil.WriteFile(fn, []byte("test"), 0600)
		defer os.Remove(fn)
	}

	present_condition := map[string][]string{
		"present": {fna},
		"message": {"file is absent"},
	}

	absent_condition := map[string][]string{
		"absent":  {fna},
		"message": {"file is present"},
	}

	suite.conditions = append(suite.conditions, present_condition)
	suite.conditions = append(suite.conditions, absent_condition)
	p, e := NewTeaConditionsProcessor(suite.conditions)

	suite.Equal(nil, e)
	suite.False(false, p.Satisfied())
	suite.Equal("file is present", p.GetInfoMessage())
}

func (suite *TeaCondTestSuite) TestPresentCrAbsentNcr() {
	/*
		present: /tmp/present.txt (created)
		absent: /tmp/absent.txt (not created)
	*/
	fna := "/tmp/present.txt"
	fnb := "/tmp/absent.txt"
	ioutil.WriteFile(fna, []byte("test"), 0600)
	defer os.Remove(fna)

	present_condition := map[string][]string{
		"present": {fna},
		"message": {"file is absent"},
	}

	absent_condition := map[string][]string{
		"absent":  {fnb},
		"message": {"file is present"},
	}

	suite.conditions = append(suite.conditions, present_condition)
	suite.conditions = append(suite.conditions, absent_condition)
	p, e := NewTeaConditionsProcessor(suite.conditions)

	suite.Equal(nil, e)
	suite.True(true, p.Satisfied())
	suite.Equal("", p.GetInfoMessage())
}

func (suite *TeaCondTestSuite) TestPresentMultipleAbsentMultiple() {
	/*
		present: /tmp/present-1.txt (created)
		         /tmp/present-2.txt (created)
		absent:  /dev/vader (not created)
		         /dev/luke  (not created)
	*/
	fna := "/tmp/present-1.txt"
	fnb := "/tmp/present-2.txt"
	for _, fn := range []string{fna, fnb} {
		ioutil.WriteFile(fn, []byte("test"), 0600)
		defer os.Remove(fn)
	}

	present_condition := map[string][]string{
		"present": {fna, fnb},
		"message": {"file is absent"},
	}

	absent_condition := map[string][]string{
		"absent":  {"/dev/vader", "/dev/luke"},
		"message": {"file is present"},
	}

	suite.conditions = append(suite.conditions, present_condition)
	suite.conditions = append(suite.conditions, absent_condition)
	p, e := NewTeaConditionsProcessor(suite.conditions)

	suite.Equal(nil, e)
	suite.True(true, p.Satisfied())
	suite.Equal("", p.GetInfoMessage())
}

func (suite *TeaCondTestSuite) TestAllPresentMultiple() {
	/*
		all-present: /tmp/present-1.txt (created)
		             /tmp/present-2.txt (created)
	*/
	fna := "/tmp/present-1.txt"
	fnb := "/tmp/present-2.txt"
	for _, fn := range []string{fna, fnb} {
		ioutil.WriteFile(fn, []byte("test"), 0600)
		defer os.Remove(fn)
	}

	present_condition := map[string][]string{
		"all-present": {fna, fnb},
		"message":     {"files are absent"},
	}

	suite.conditions = append(suite.conditions, present_condition)
	p, e := NewTeaConditionsProcessor(suite.conditions)

	suite.Equal(nil, e)
	suite.True(true, p.Satisfied())
	suite.Equal("", p.GetInfoMessage())
}

func (suite *TeaCondTestSuite) TestAllAbsentMultiple() {
	/*
		all-present: /tmp/present-1.txt (created)
		             /tmp/present-2.txt (created)
	*/
	fna := "/tmp/present-1.txt"
	fnb := "/tmp/present-2.txt"

	present_condition := map[string][]string{
		"all-absent": {fna, fnb},
		"message":    {"files are present!??"},
	}

	suite.conditions = append(suite.conditions, present_condition)
	p, e := NewTeaConditionsProcessor(suite.conditions)

	suite.Equal(nil, e)
	suite.True(true, p.Satisfied())
	suite.Equal("", p.GetInfoMessage())
}
