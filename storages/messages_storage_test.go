package storages

import (
	"io/ioutil"
	"os"
	"testing"

	"time"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

const testMessagesDB = "test_messages.db"

func TestMessagesStorage(t *testing.T) {
	logrus.SetOutput(ioutil.Discard) // do not write logs during test
	Convey("Test messages storage", t, func() {
		storage, err := NewMessagesStorage(testMessagesDB, nil)
		So(err, ShouldBeNil)

		testValue := &MessagesStorageValue{
			UserId:       "user",
			TemplateName: "template",
			Variables:    map[string]string{"a": "1", "b": "2"},
			CreatedAt:    time.Now().UTC(),
			Message:      "message",
		}

		So(storage.PutValue("id", testValue), ShouldBeNil)
		v, err := storage.GetValue("id")
		So(err, ShouldBeNil)
		So(v, ShouldResemble, testValue)

		_, err = storage.GetValue("blah")
		So(err, ShouldEqual, ErrMessageNotExists)

		// cleanup
		So(storage.Close(), ShouldBeNil)
		So(os.Remove(testMessagesDB), ShouldBeNil)
	})
}
