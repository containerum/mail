package storages

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

const testDB = "test.db"

func TestTemplateStorage(t *testing.T) {
	logrus.SetOutput(ioutil.Discard) // do not write logs during test
	Convey("Test template storage", t, func() {
		storage, err := NewTemplateStorage(testDB, nil)
		So(err, ShouldBeNil)

		Convey("Put and get one version", func() {
			So(storage.PutTemplate("tmpl1", "1", "a"), ShouldBeNil)
			v, err := storage.GetTemplate("tmpl1", "1")
			So(err, ShouldBeNil)
			So(v.Data, ShouldResemble, "a")

			_, err = storage.GetTemplate("unknown", "1")
			So(err, ShouldEqual, ErrTemplateNotExists)

			So(storage.PutTemplate("tmpl2", "1", "a"), ShouldBeNil)
			_, err = storage.GetTemplate("tmpl2", "-1")
			So(err, ShouldEqual, ErrVersionNotExists)
		})

		Convey("Put and get multiple versions", func() {
			So(storage.PutTemplate("tmpl3", "1", "a"), ShouldBeNil)
			So(storage.PutTemplate("tmpl3", "2", "b"), ShouldBeNil)

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
			So(err, ShouldEqual, ErrTemplateNotExists)
		})

		Convey("Deleting", func() {
			So(storage.PutTemplate("tmpl4", "1", "a"), ShouldBeNil)
			So(storage.PutTemplate("tmpl4", "2", "b"), ShouldBeNil)
			So(storage.PutTemplate("tmpl4", "3", "c"), ShouldBeNil)

			So(storage.DeleteTemplate("tmpl4", "1"), ShouldBeNil)
			_, err := storage.GetTemplate("tmpl4", "-1")
			So(err, ShouldEqual, ErrVersionNotExists)

			So(storage.DeleteTemplate("tmpl4", "-1"), ShouldEqual, ErrVersionNotExists)

			So(storage.DeleteTemplates("tmpl4"), ShouldBeNil)
			_, err = storage.GetTemplate("tmpl4", "2")
			So(err, ShouldEqual, ErrTemplateNotExists)
			_, err = storage.GetTemplates("tmpl4")
			So(err, ShouldEqual, ErrTemplateNotExists)

			So(storage.DeleteTemplates("tmpl4"), ShouldEqual, ErrTemplateNotExists)
		})

		// cleanup
		So(storage.Close(), ShouldBeNil)
		So(os.Remove(testDB), ShouldBeNil)
	})
}
