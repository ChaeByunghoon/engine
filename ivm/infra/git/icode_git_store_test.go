/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package git_test

import (
	"errors"
	"os"
	"testing"

	"encoding/hex"
	"fmt"

	"github.com/it-chain/engine/ivm"
	"github.com/it-chain/engine/ivm/infra/git"
	"github.com/stretchr/testify/assert"
)

func TestICodeGitStoreApi_Clone(t *testing.T) {
	baseTempPath := "./.tmp"
	sshPath := "./id_rsa"
	err, teardown := generatePriKey(sshPath)
	assert.NoError(t, err, "err in generate pri key file")
	defer teardown()
	os.RemoveAll(baseTempPath)
	defer os.RemoveAll(baseTempPath)

	//given
	tests := map[string]struct {
		ID          string
		InputGitURL string
		OutputMeta  ivm.ICode
		OutputErr   error
	}{
		"success": {
			ID:          "1",
			InputGitURL: "github.com/junbeomlee/test_icode",
			OutputMeta:  ivm.ICode{RepositoryName: "test_icode", GitUrl: "github.com/junbeomlee/test_icode", Path: baseTempPath + "/" + "test_icode"},
			OutputErr:   nil,
		},
		"fail": {
			ID:          "2",
			InputGitURL: "github.com/nonono",
			OutputMeta:  ivm.ICode{},
			OutputErr:   errors.New("repository not found"),
		},
	}

	icodeApi := git.NewRepositoryService()

	for _, test := range tests {
		//when
		meta, err := icodeApi.Clone(test.ID, baseTempPath, test.InputGitURL, sshPath, "")

		if err == nil {
			// ivm ID 는 랜덤이기때문에 실데이터에서 주입
			// commit hash 는 repo 상황에따라 바뀌기 때문에 주입
			test.OutputMeta.ID = meta.ID
			test.OutputMeta.CommitHash = meta.CommitHash
		}

		//then
		assert.Equal(t, test.OutputMeta, meta)
		assert.Equal(t, test.OutputErr, err)
	}
}

func TestICodeGitStoreApi_CloneWithPassword(t *testing.T) {
	baseTempPath := "./.tmp"
	sshPathWithPassword := "./id_rsa_pwd"
	err, teardown := generatePriKeyWithPassword(sshPathWithPassword)
	assert.NoError(t, err, "err in generate pri key file")
	defer teardown()
	os.RemoveAll(baseTempPath)
	defer os.RemoveAll(baseTempPath)

	//given
	tests := map[string]struct {
		ID          string
		InputGitURL string
		InputPwd    string
		OutputMeta  ivm.ICode
		OutputErr   error
	}{
		"success": {
			ID:          "1",
			InputGitURL: "github.com/nesticat/test_icode",
			InputPwd:    "pwdtest",
			OutputMeta:  ivm.ICode{RepositoryName: "test_icode", GitUrl: "github.com/nesticat/test_icode", Path: baseTempPath + "/" + "test_icode"},
			OutputErr:   nil,
		},
		"fail": {
			ID:          "1",
			InputGitURL: "github.com/nesticat/test_icode",
			InputPwd:    "wrongpassword",
			OutputMeta:  ivm.ICode{},
			OutputErr:   errors.New("x509: decryption password incorrect"),
		},
	}

	icodeApi := git.NewRepositoryService()

	for _, test := range tests {
		//when
		meta, err := icodeApi.Clone(test.ID, baseTempPath, test.InputGitURL, sshPathWithPassword, test.InputPwd)

		if err == nil {
			// ivm ID 는 랜덤이기때문에 실데이터에서 주입
			// commit hash 는 repo 상황에따라 바뀌기 때문에 주입
			test.OutputMeta.ID = meta.ID
			test.OutputMeta.CommitHash = meta.CommitHash
		}

		//then
		assert.Equal(t, test.OutputMeta, meta)
		assert.Equal(t, test.OutputErr, err)
	}
}

func generatePriKey(path string) (error, func() error) {
	_, err := os.Stat(path)
	if err == nil {
		os.Remove(path)
	}

	file, err := os.Create(path)
	if err != nil {
		return err, nil
	}

	src := "2d2d2d2d2d424547494e205253412050524956415445204b45592d2d2d2d2d0a4d4949456f67494241414b4341514541784237474d566a7135376e314f49646f736950354341677141506b53334641417231416b4c4f6e48424331484d636c530a66586e7732745a54722b41386474764e37442b71364c5a793666634551623555546271706e4a51706a39304e356268366f4f317a7538433649586c52363237330a586f525536334f784170766f5a772f4d5a337a2f3467533679324343466c7a59426c743635336663787a754e7a3359374275325a476f5334732f4135783958520a656962622b6c74596732655658514562627936794d74455862557850652f7757664c583268544630316645374a67517a35766f526644356e3164396c326f37490a53734137533135464535334254673954334f6f553161553964586368316f32786e6d5142647a5465534a466d54766a6f624f69426650315665444d6d6a71624e0a4c335174617a3352574a4266756d6e45647969737870746d55346d6142624a7a536e2b634b51494441514142416f4942414632436a7431596d436945386664530a47516c58505a596d7a6d4249596b584a6e346e336e45674e37325a2b63454f38796967707a44324c6b3774344831784d305a4b6a694d6f4d744233364f5831660a55724c394859496134765a4659437234477742414e373539316b472f70742b717553664830505779342b4e716b7855513431553074497a2f3146444559304a6d0a596c6f6c7043525a636c744d65674642546b4f765a69444f78344b44596b4a6c436d6f65616e36325a384d5a775a6d7532597a37756a32455a7666764f586a730a7a3749504b5a6a6f524c2b447a577142466e5347706d6b5876324c7a5334546c363041747545354768494558392f777053346b583568376734314b4d2f3375620a336145764a7949486841545855626373344b436a70634a75757a61463038476d63614b79424439416a4d6b6f595553684945716d2b566d53677533524d5877520a43372f77344145436759454135794167786a744e4a364c5435346a446c5a6c4132544f61735a6d665654686c504c5674503866547847426650493377454954370a4a6c3355526e53392f6e6855684d4c33477465684355786b6f7173304e6256672f6159516b78696156394c65622f357439444e31636e482b47616141697656410a5149334d35304c4a6a344b5251646b664d6e6467614f46475a7150684e3645307439413070764f7a357832456555534958546669597345436759454132546f780a2f5675646f63434558356b385271616c775157486e7946455664564c567076474e7739576f5837394f3847315338356d784870622b4951736c4f31795954692b0a65654b5a4239496e77483532637168453470547637592f6b76413379334947456d4a75345772576f77697768746d5873586c73544b532b7431424b35526364460a352f744b6f693937554f383848527936446b4a3863645268544a7839756b57716579437132326b43675941736c58503943543975332b674568387443746c64650a4471684f6a6858414f4b712b7454696e7a77493470575a3570642b6a4d43504b574e737a354230715530667166446c7967686e635531497556747778614257580a6d45736d4e4e3742426a70475845775669542b6b6e66796f4d67676c7866316f396e474b51735869327772754b74587278443969752b4836747134684c7757650a56356c77677834322f4f6972413939534c412b4e67514b4267446d36425637573465554354437337685a45673642754c5a4b63644b4250485175595a4c3275690a58397336374144645556683732554f4e594c4f434c4862485177596a466a7339784830586c416a4c6b703656715069747137547438464d7051636a6e676c30720a784b6f5762477074582b67673364655654466f396d5777714c61496c65715a5457566f5156433046356d7532487074376636616755647353477a644e48436a730a58587442416f474146476f6559706c39562b7a54525974303366582f63504b34546c335a666a48325a4b67704372484951786e6a4957426f37706d31734159590a62764f5743757a704b486d6d4e794443504154707851564d6d417752416762635664736365342f7930624b467161784c387a6c3731614a6e324e3769344567540a6a6262346155396c39534d70314366744933355766394d76714341496b5476372b675351674a726d4161636e696a436f6431343d0a2d2d2d2d2d454e44205253412050524956415445204b45592d2d2d2d2d0a"
	p, _ := hex.DecodeString(src)

	file.WriteString(fmt.Sprintf("%s", p))
	return nil, func() error {

		if err := os.RemoveAll(path); err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}

		return nil
	}
}

func generatePriKeyWithPassword(path string) (error, func() error) {
	_, err := os.Stat(path)
	if err == nil {
		os.Remove(path)
	}

	file, err := os.Create(path)
	if err != nil {
		return err, nil
	}

	src := "2d2d2d2d2d424547494e205253412050524956415445204b45592d2d2d2d2d0a50726f632d547970653a20342c454e435259505445440a44454b2d496e666f3a204145532d3132382d4342432c36383741323545374230353641444431343139343343334643373831453638380a0a4b666e3757374d6f5777737a494c5a4d5552476b734448742b46414751445a324e3069734530626e73537245544a5a66303650715a367a4a5878353477324f510a5744354541715253366776336d4c4564364c37336243586a69484b6875304933506b52626f2f35476d4a5432572f2b4867764d6c4945796533796b74677430650a70385078356769652b7765742b4f76524c79493958374145685465745a6c35683871474652516c38596361495546345776623854384143596a543552456554680a62505a324968513672396f467063354b627a304362785841326c3931316537615055504c506165766f49573377384d7274676d79776e554d62426e45345a32310a4457477444745447334c4b35534141747930615445353678666d6a33566649463972524e5459534761783947754e54502f6f475847514e63706e7266334875470a64574335766d517166567a42544e32744847545a57322b494a516b4b34766f446d6b69736f62385434563049416f46683272764544684b51732b6374546e56630a646b634e5537304e6b5869573562534e53676179416553347062425379377248574774423739335853456251586e31556f7037384e4b4755686a795137542b460a635544747066417577735569365833644d7669776c5769322b78493073595475507553376230646566584c44417a78634e72496d4a7232704d503579363169340a474c69536238417361586268553778363134476e5844316b2b397a6f434755504a41742f796436367957616c4f74535176447a57476d71756855714877486a450a6e703469552f33585476344452744c6169715955534a4456616835544641536464797a38774c6e52476f5a4b37476c347555586468456b696f5a5576793641320a6c434b555a4441396c376c796a78385534664c7a636f6d636f2f78464f5830304936784b4579353650382b71502b35454f7648384a4d79394f462b4c433837660a6d764158373031432b306446794d3532357243504b446e776e34626e654277796848663754542b614862466934666a70707650504a5a683770446571454c6e4d0a6e2b4330663050776b656949546330643373527a6645374b4b786f514232474a2b2f754565552f45612f71726564562f674b45584f4b366152357855522f6a450a545342364e656b577746773669747561536c666e47593949466738574e3942357079336e3558347034746f5a4e4e737231674c6d43326d794a5a6d563470457a0a4d665275704f433268395a4d3243372b4b72626a65473448647679737765625331484263594c71354f4d77506173516f4e7a54365454684d3764724242762f710a7631594c2f385355664e37736e4d79547272394346396a6b6c7a676c6e3242346a7a68766b2b76564c46337a6437794b755748614e553836305166416a4e35330a355a4175356930593279694e756b503978775938484841673461746c38703576594f667672516a633470515766304d744f344a386a37783034797a6e586138530a62664f783730306959367a572f6f4433363679766d6870764369494a7369384b56546b654f4d433761626a543878505258703656526361776e667247766466440a306b554e67774c72544c6b4b76426e796937577a3130646a37323850653431334f503161396d7a56796277397037304c6b7a50426948726f32352f5550486a7a0a61717159774f715a33352f553073573570517736584872776533434235334378436e416d7168734e56757939644575485a57465167306e31696132714c42504c0a4b6c384157452f49796e692f4d4c54376870726a446a39796a6676414d555833577a613736435570426f506943624e67504a32785338556b676e6a46476b71410a4771383168367a7647326f5a4661417645345470557a36454b77504269386842373147584d575771504a675637524253467a5062526a5577373650414f5657450a52394f4f39766b59776231365742394a4a586a586177766f6958484b4e6836677366344e4f764e6c6d494d2f6d655364385a53666167575573384c59345a78590a70574a554e4973795a792b2f346d7333316b7033453948436a636c78356e4352315a6742314f32505438454362584867714936416251774f79385857315436680a4675683949424a7339783959505776336c5937506655564951732f63557a514d3071666e5363565a54433835694f45554177623958766c6b4646346567774b720a2d2d2d2d2d454e44205253412050524956415445204b45592d2d2d2d2d0a"
	p, _ := hex.DecodeString(src)

	file.WriteString(fmt.Sprintf("%s", p))
	return nil, func() error {

		if err := os.RemoveAll(path); err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}

		return nil
	}
}
