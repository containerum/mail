package bolt

import (
	"io/ioutil"
	"os"
	"testing"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

const testDB = "test.db"

func TestTemplateStorage(t *testing.T) {
	logrus.SetOutput(ioutil.Discard) // do not write logs during test
	Convey("Test template storage", t, func() {
		storage, err := NewBoltTemplateStorage(testDB, nil)
		So(err, ShouldBeNil)

		Convey("Put and get one version", func() {
			So(storage.PutTemplate("tmpl1", "1", "a", "s"), ShouldBeNil)
			v, err := storage.GetTemplate("tmpl1", "1")
			So(err, ShouldBeNil)
			So(v.Data, ShouldResemble, "a")

			_, err = storage.GetTemplate("unknown", "1")
			So(err, ShouldEqual, mttypes.ErrTemplateNotExists)

			So(storage.PutTemplate("tmpl2", "1", "a", "s"), ShouldBeNil)
			_, err = storage.GetTemplate("tmpl2", "-1")
			So(err, ShouldEqual, mttypes.ErrVersionNotExists)
		})

		Convey("Put and get multiple versions", func() {
			So(storage.PutTemplate("tmpl3", "1", "a", "s"), ShouldBeNil)
			So(storage.PutTemplate("tmpl3", "2", "b", "s"), ShouldBeNil)

			v, err := storage.GetTemplate("tmpl3", "1")
			So(err, ShouldBeNil)
			So(v.Data, ShouldEqual, "a")

			v, err = storage.GetTemplate("tmpl3", "2")
			So(err, ShouldBeNil)
			So(v.Data, ShouldEqual, "b")

			mp, err := storage.GetTemplates("tmpl3")
			So(err, ShouldBeNil)
			So(mp["1"].Data, ShouldEqual, "a")
			So(mp["2"].Data, ShouldEqual, "b")

			_, err = storage.GetTemplates("unknown")
			So(err, ShouldEqual, mttypes.ErrTemplateNotExists)
		})

		Convey("Deleting", func() {
			So(storage.PutTemplate("tmpl4", "1", "a", "s"), ShouldBeNil)
			So(storage.PutTemplate("tmpl4", "2", "b", "s"), ShouldBeNil)
			So(storage.PutTemplate("tmpl4", "3", "c", "s"), ShouldBeNil)

			So(storage.DeleteTemplate("tmpl4", "1"), ShouldBeNil)
			_, err := storage.GetTemplate("tmpl4", "-1")
			So(err, ShouldEqual, mttypes.ErrVersionNotExists)

			So(storage.DeleteTemplate("tmpl4", "-1"), ShouldEqual, mttypes.ErrVersionNotExists)

			So(storage.DeleteTemplates("tmpl4"), ShouldBeNil)
			_, err = storage.GetTemplate("tmpl4", "2")
			So(err, ShouldEqual, mttypes.ErrTemplateNotExists)
			_, err = storage.GetTemplates("tmpl4")
			So(err, ShouldEqual, mttypes.ErrTemplateNotExists)

			So(storage.DeleteTemplates("tmpl4"), ShouldEqual, mttypes.ErrTemplateNotExists)
		})

		// cleanup
		So(storage.Close(), ShouldBeNil)
		So(os.Remove(testDB), ShouldBeNil)
	})
}
