// +build !skippackr
// Code generated by github.com/gobuffalo/packr/v2. DO NOT EDIT.

// You can use the "packr2 clean" command to clean up this,
// and any other packr generated files.
package packrd

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/packr/v2/file/resolver"
)

var _ = func() error {
	const gk = "bf099a2c0ff496e0b22de2af5c310957"
	g := packr.New(gk, "")
	hgr, err := resolver.NewHexGzip(map[string]string{
		"2a3d1e318604e12c1da8f4d2166de45c": "1f8b08000000000000ffec584d8fd33010bde757585e8e5d0227a49ef782c40189235aaddc66d27a37b1cd780c2aa8ff1d2569ab3871dcd0a6888f5c7ab033336f9edfccd8fd9130c65fd9f5164ac1978c6f89cc324d9fad56f7cdea6b8d9b344391d3fd9b7769b376c717955d06b95492a456962f59e58a312e952551149f08cc6991314e3b035500bd7a8635d5f6f5ba416d0049826d7dcd187f71967429bf83b73cec28e6acdecbc0ae519a0a6d6fb3e20021afdcdea5adac524b601e5a869ed57ee14750a2eca2f5105b42a936711f081a33c06bdd1cd91315ea27a98ca3984b8128767cd1dd9604659fc91098ae2963bc94eafdc1fe6d67777f06fb0a5001817d2a859239587ad28ecea43086956f42467dacb42e409c3965e1483f8029f46e8cab1ea719e4c215158c5c1416a2a12c44e1062a807599ef6d8a2cabb52d8a8fed622174712c5f45e1026535a5887e49258d228280c2a57c34887a350245090418cd7488f8016eafcd3c19c03b14af1696f721c217271132be649f033d6b31dc2947f495f1e5dbfaf031e9a6774a2d8cb635100e6b8f2783180f49cb3777668322837936f510cf43e55256c60e957f66124cd479e7ae16ee6a5ec3b224c8d9ffb05f211420ec8d2f1c5735985bea370bd334b9bc9c9a5f6bf3449c7a22fe49c3ca38dc448ff9b2f1fcbb8aff2f1a547ee9057a4947c3d5797f00b5a1ed49ba4767fd1bc64009054a675866c3fd26dc0b229af7cfdcc21a03921f6dff02fdbf1262c649c0cd48e179877bb8269d125834508eda68e254bf759c10b1214a8f7f005e767ae1c6de9e523e3aff613765ccf65b3118b3b99b4d19b275db0b67a96ec0ad7f071838fbe80b7b9ffc0c0000ffffb7c3216946160000",
	})
	if err != nil {
		panic(err)
	}
	g.DefaultResolver = hgr

	func() {
		b := packr.New("schema", "./schema")
		b.SetResolver("kustomize.json", packr.Pointer{ForwardBox: gk, ForwardPath: "2a3d1e318604e12c1da8f4d2166de45c"})
		}()

	return nil
}()