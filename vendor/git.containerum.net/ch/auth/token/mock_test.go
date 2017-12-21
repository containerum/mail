package token

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMockIssuerValidator(t *testing.T) {
	Convey("Test mock token issuer-validator", t, func() {
		mockiv := NewMockIssuerValidator(time.Hour)
		Convey("generate and validate tokens", func() {
			accessToken, refreshToken, err := mockiv.IssueTokens(ExtensionFields{})
			So(err, ShouldBeNil)
			So(accessToken.LifeTime, ShouldEqual, mockiv.returnedLifeTime)
			So(refreshToken.LifeTime, ShouldEqual, mockiv.returnedLifeTime)
			So(accessToken.Id, ShouldResemble, refreshToken.Id)

			valid, err := mockiv.ValidateToken(accessToken.Value)
			So(err, ShouldBeNil)
			So(valid.Id, ShouldResemble, accessToken.Id)
			So(valid.Valid, ShouldBeTrue)
			So(valid.Kind, ShouldEqual, KindAccess)

			valid, err = mockiv.ValidateToken(refreshToken.Value)
			So(err, ShouldBeNil)
			So(valid.Id, ShouldResemble, refreshToken.Id)
			So(valid.Valid, ShouldBeTrue)
			So(valid.Kind, ShouldEqual, KindRefresh)
		})
		Convey("validate invalid token", func() {
			_, err := mockiv.ValidateToken("invalid")
			So(err, ShouldNotBeNil)
		})
	})
}
