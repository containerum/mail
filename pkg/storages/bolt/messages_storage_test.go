package bolt

import (
	"io/ioutil"
	"os"
	"testing"

	"time"

	"git.containerum.net/ch/mail-templater/pkg/models"
	"git.containerum.net/ch/mail-templater/pkg/mtErrors"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

const testMessagesDB = "test_messages.db"

func TestMessagesStorage(t *testing.T) {
	logrus.SetOutput(ioutil.Discard) // do not write logs during test
	Convey("Test messages storage", t, func() {
		storage, err := NewBoltMessagesStorage(testMessagesDB, nil)
		So(err, ShouldBeNil)

		testValue := &models.MessagesStorageValue{
			UserId:       "user",
			TemplateName: "template",
			Variables:    map[string]interface{}{"a": "1", "b": "2"},
			CreatedAt:    time.Now().UTC(),
			Message:      "message",
		}

		So(storage.PutMessage("id", testValue), ShouldBeNil)
		v, err := storage.GetMessage("id")
		So(err, ShouldBeNil)
		So(v, ShouldResemble, testValue)

		_, err = storage.GetMessage("blah")
		So(err.Error(), ShouldEqual, mtErrors.ErrMessageNotExist().Error())

		// cleanup
		So(storage.Close(), ShouldBeNil)
		So(os.Remove(testMessagesDB), ShouldBeNil)
	})
}
