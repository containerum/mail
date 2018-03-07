package reqtests

import (
	"testing"

	"git.containerum.net/ch/kube-client/pkg/model"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	resourceTestVolume = "n1-volume"
)

func TestVolume(test *testing.T) {
	client := newMockClient(test)
	Convey("Test volume methods", test, func() {
		Convey("resource api", func() {
			var volumes []model.Volume
			Convey("get volume and list", func() {
				var err error
				volumes, err = client.GetVolumeList(nil)
				So(err, ShouldBeNil)
				So(len(volumes), ShouldBeGreaterThan, 0)
				gainedVolume, err := client.GetVolume(volumes[0].Label)
				So(err, ShouldBeNil)
				So(gainedVolume, ShouldResemble, volumes[0])
				So(client.RenameVolume(gainedVolume.Label, "abanamat"), ShouldBeNil)
				So(client.RenameVolume("abanamat", gainedVolume.Label), ShouldBeNil)
			})
		})
	})
}
