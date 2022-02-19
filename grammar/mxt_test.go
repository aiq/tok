package grammar

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aiq/tok"
)

func TestMXTHeader(t *testing.T) {

}

func TestMXT(t *testing.T) {
	posCases := []struct {
		mxt   string
		names []string
	}{
		{
			`// name.of.chunk -->`,
			[]string{"name.of.chunk"},
		},
		{
			`//salt++++++++++++++++++++++++++ Σ-element ++++++++++++++++++++++++++++++++++-->`,
			[]string{"Σ-element"},
		},
		{
			`//---------------------------------------------------------------- user.json -->
{
    "user": "alucard",
    "password": "C:SotN1997"
}
//--------------------------------------------------------------- connection.ini
// comment line that is not part of the ini file,
// comment lines will be joined with a space character
//
// empty comment lines will generate a newline character in the comment
//----------------------------------------------------------------------------->
request: GET

[url]
schema=http
host=localhost
port=8080
path=/db/add
// user.pgp --> XYZ
-----BEGIN PGP MESSAGE-----

hQEMA8p144+Gi+YpAQf/VeFG9Zb+8w9aldWll8n2g3jqpE613LKg2XAJgwXQmSQL
R4O+TlQakJ+Mz5vM4IxxubPgYCyt6cyL7qM3oJIuk7vsqMbl5t7c/dOfXjj7goIC
IskIX+9e5qrr8jRG/KZYSdBJtFEI9oNtZTLlnv3yeV3OWNTbUnjdTWrk/h1kavJE
5psOaTgH+Yg0utcF9Y0AqTo/o6EmiAn+uQlE5l9GTSzq51jRNiMoZRaXfUeznneF
SDn/eFve7VlLV60j8LcXSrXgiIXEYi9ZgvWHHV+h9yorQXu24fgrEd0AhE/adiJm
dzEHVCb+Q0i0p14hE2QQYdDE0zLxZN+FZJy/MiWwHtLA0QFZ+Zh47Zmiq14WyMPp
4ASQ/wNpwmZmpsiKDweUTlng3MTjifopmTeglPmRKsIOvUbTnA3yetznViugSYrA
jPuSltUZpdLTrOPrOrbPorxKmurKmswZje1hoXIrpPN2KzTnA0UNH6GKOnMjLtRH
uXMlS3q/GkIcXYDJIVFsEU3dGWMJFOtF4HPlyHmB8Gko8cDep6ugPbgJt8LBu3+p
zE00pkAd4Kf7yodSUR6mRJm5qK+vqp7dnkZd/WhTNRihMoqRY2oKiH2rHjFSa76s
D8nD4suo6ckVzYGJpknGSIAwaCFlW0aqR/3SWO4wi6ibbfub8LA73V90Ll3/S/Ph
xU15HYmdCATnVX1sp1PWmyz972bMvl8txyIKMUueVw+w0C19ZTfWXjuFSguF7zt7
RY+I3to2lbyVJbcI9Dyz04GOJZ2vIhG9eq65FxeweAKDa7L+iH1NA5L2lYd9DEr1
ro/CU6vIqkOSNRUrNYDwqz1g3Z3eAQB/8t9Y4WsV4KL0M229rsFrtl26i7+quYfg
uuTd
=WxK9
-----END PGP MESSAGE-----
//XYZ hello-world.h -->
//---------------------------------------------------------- hello-world.c --> X
// this is part of hello-world.c
#include<stdio.h>

int main(void) {
printf("Hello World\n");
return 0;
}`,
			[]string{"user.json", "connection.ini", "user.pgp", "hello-world.h", "hello-world.c"},
		},
	}
	for i, c := range posCases {
		sca := tok.NewScanner(c.mxt)
		r := MXT()
		names := []string{}
		r.Name.Map(func(t tok.Token) {
			names = append(names, sca.Get(t))
		})
		err := sca.Use(r)
		if err != nil {
			t.Errorf("%d unexpected error: %v", i, err)
		}
		if !sca.AtEnd() {
			t.Errorf("did not read the whole mxt: %s", sca.Tail())
		}
		if len(names) != len(c.names) {
			t.Errorf("unexpected number of names: %d != %d: %s", len(names), len(c.names), strings.Join(names, " "))
		} else {
			for i, s := range names {
				if s != c.names[i] {
					t.Errorf("unexpected name at %d: %q != %q", i, s, c.names[i])
				}
			}
		}
	}
	fmt.Println(strings.Join(MXT().Grammar().Lines(), "\n"))

}
