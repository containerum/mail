package bolt

import (
	"io/ioutil"
	"os"
	"testing"

	"git.containerum.net/ch/mail-templater/pkg/mtErrors"
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
			So(storage.PutTemplate("tmpl1", "1", "a", "s", true), ShouldBeNil)
			v, err := storage.GetTemplate("tmpl1", "1")
			So(err, ShouldBeNil)
			So(v.Data, ShouldResemble, "a")

			_, err = storage.GetTemplate("unknown", "1")
			So(err.Error(), ShouldEqual, mterrors.ErrTemplateNotExist().Error())

			So(storage.PutTemplate("tmpl2", "1", "a", "s", true), ShouldBeNil)
			_, err = storage.GetTemplate("tmpl2", "-1")
			So(err.Error(), ShouldEqual, mterrors.ErrTemplateVersionNotExist().Error())
		})

		Convey("Put and get multiple versions", func() {
			So(storage.PutTemplate("tmpl3", "1", "a", "s", true), ShouldBeNil)
			So(storage.PutTemplate("tmpl3", "2", "b", "s", true), ShouldBeNil)

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
			So(err.Error(), ShouldEqual, mterrors.ErrTemplateNotExist().Error())
		})

		Convey("Deleting", func() {
			So(storage.PutTemplate("tmpl4", "1", "a", "s", true), ShouldBeNil)
			So(storage.PutTemplate("tmpl4", "2", "b", "s", true), ShouldBeNil)
			So(storage.PutTemplate("tmpl4", "3", "c", "s", true), ShouldBeNil)

			So(storage.DeleteTemplate("tmpl4", "1"), ShouldBeNil)
			_, err := storage.GetTemplate("tmpl4", "-1")
			So(err.Error(), ShouldEqual, mterrors.ErrTemplateVersionNotExist().Error())

			So(storage.DeleteTemplate("tmpl4", "-1").Error(), ShouldEqual, mterrors.ErrTemplateVersionNotExist().Error())

			So(storage.DeleteTemplates("tmpl4"), ShouldBeNil)
			_, err = storage.GetTemplate("tmpl4", "2")
			So(err.Error(), ShouldEqual, mterrors.ErrTemplateNotExist().Error())
			_, err = storage.GetTemplates("tmpl4")
			So(err.Error(), ShouldEqual, mterrors.ErrTemplateNotExist().Error())

			So(storage.DeleteTemplates("tmpl4").Error(), ShouldEqual, mterrors.ErrTemplateNotExist().Error())
		})

		// cleanup
		So(storage.Close(), ShouldBeNil)
		So(os.Remove(testDB), ShouldBeNil)
	})
}
