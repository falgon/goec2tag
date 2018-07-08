package utils

import (
	"os"
	"fmt"
	"flag"
	"strings"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
)

var (
	// Region 設定
	ArgRegion = flag.String("region", "ap-northeast-1", "Region name")
	// Endpoint 設定
	ArgEndpoint = flag.String("endpoint", "", "Endpoint.")
	// インスタンス ID またはそのタグの設定
	ArgInstances = flag.String("instances", "", "Instance id or instance tag name.")
	// 設定対象タグの指定
	ArgTags = flag.String("tags", "", "Tag Key(Use Key=) and Tag Value(Use Value=)\n[Example]:\n\t ... -tags='Key=foo,Value=bar Key=hoge,Value=piyo...'")
	// タグ追加フラグ
	ArgAdd = flag.Bool("addT", false, "Give the tag to the instance.")
	// タグ削除フラグ
	ArgDel = flag.Bool("rmT", false, "Remove tag from instance.")
	// タグ表示フラグ
	ArgShowTags = flag.Bool("showtags", false, "DescribeTags API operation for EC2.\n Describes one or more of the tags for your EC2 resources. filter=...")
	// タグ表示の際のフィルタ
	ArgShowTagsFilter = flag.String("filter", "", "This flag is used in conjunction with the showtags flag to filter tags by describing filter statements.\n[Example]:\n\t ... -filter 'name:resource-id,values:i-xxxxxxxx i-yyyyyyyy'")
)

// このプログラムを動作させた EC2 インスタンスのインスタンスIDとエラー情報を返す.
//
// 使用 API: <https://docs.aws.amazon.com/sdk-for-go/api/aws/ec2metadata/#EC2Metadata.GetInstanceIdentityDocument>
func GetThisInstanceId() (r string, err error) {
	svc := ec2metadata.New(session.Must(session.NewSession()))
	var doc ec2metadata.EC2InstanceIdentityDocument
	if doc, err = svc.GetInstanceIdentityDocument(); err == nil {
		r = doc.InstanceID
	}
	return
}

// * msg: エラー情報の文字列
//
// * args...: エラーオブジェクトなど
//
// エラー情報を出力し, プログラムを終了する.
func ExitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+msg+"\nDetali:\n", args...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

// * s: エラー情報の文字列
//
// * err: エラーオブジェクト
//
// エラーオブジェクトが `nil` であれば何もせず, そうでなければ `ExitErrorf` を実行する.
func Unwrap(s string, err error) {
	if err != nil {
		ExitErrorf(s, err)
	}
}

// * reg: Region 情報の文字列
// 
// 設定された Region と Endpoint から ec2.EC2 オブジェクトを構築してそのポインタを返す.
//
// 使用 API: 
//
// * <https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#New>
func NewEc2Client (reg string) *ec2.EC2 {
	cfg := aws.Config{
		Region: aws.String(reg),
		Endpoint: aws.String(*ArgEndpoint),
	}
	return ec2.New(session.New(&cfg))
}

func splitxs(instances string, s string) (r []*string) {
	for _, i := range strings.Split(instances, s) {
		r = append(r, aws.String(i))
	}
	return
}

// * cli: 初期化済みの *ec2.EC2
//
// * instances: (単数 ∈)複数のインスタンス情報の文字列
//
// * tags: 付与するタグ
//
// 指定した EC2 リソースに 1 つ以上のタグを追加または上書きし, エラー情報を返す.
// 各リソースは、最大 50 のタグを持つことができ, 各タグはキーとオプションの値で構成されており,
// タグキーはリソースごとに一意でなければならない. 
// 
// 使用 API: 
//
// * <https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#EC2.CreateTags>
// 
// * <https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#CreateTagsInput>
//
// * <https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#CreateTagsOutput>
func CreateTag(cli *ec2.EC2, instances string, tags []*ec2.Tag) (err error) {
	input := &ec2.CreateTagsInput {
		Resources: splitxs(instances, ","),
		Tags: tags,
	}
	var r *ec2.CreateTagsOutput
	if r, err = cli.CreateTags(input); err == nil {
		fmt.Println(r)
	}
	return
}


// * cli: 初期化済みの *ec2.EC2
//
// * filter: フィルタ文字列. Name と Values が利用可能. 各コンテンツに関する詳細は[公式ドキュメント](https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/APIReference/API_Filter.html)を参照.
// 
// フィルタ規則に従ってエラーがなければ 1 つ以上のタグ情報に関する情報を出力し, エラー情報を返す.
// フィルタ規則は, [aws cli](https://docs.aws.amazon.com/cli/latest/reference/ec2/describe-tags.html) とほぼ同様であるが少し異なり,
// 例えば 'name:resource-id, values:i-xxxxxxxx i-yyyyyyyy' という文字列を期待する.
//
// 使用 API: 
//
// * <https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#EC2.DescribeTags>
//
// * <https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeTagsInput>
//
// * <https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeTagsOutput>
func DescribeTag(cli *ec2.EC2, filter string) (err error) {
	var input *ec2.DescribeTagsInput
	if filter == "" {
		input = &ec2.DescribeTagsInput {}
	} else {
		var (
			name string
			values []*string
		)
		for _, t := range strings.Split(filter, ",") {
			section := strings.Split(t, ":")
			switch section[0] {
				case "name" : name = section[1]
				case "values": values = splitxs(section[1], " ")
				default: err = errors.New("Fliter parse error.")
			}
		}

		input = &ec2.DescribeTagsInput {
			Filters: []*ec2.Filter{
				{
					Name: aws.String(name),
					Values: values,
				},
			},
		}
	}

	var r *ec2.DescribeTagsOutput
	if r, err = cli.DescribeTags(input); err == nil {
		fmt.Println(r)
	}
	return
}

// * cli: 初期化済みの *ec2.EC2
//
// * instances: (単数 ∈)複数のインスタンス情報の文字列
//
// * tags: 削除対象となるタグ
//
// 指定した EC2 リソースから 1 つ以上のタグを除き, エラー情報を返す.
//
// 使用 API: 
//
// * <https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#EC2.DeleteTags>
//
// * <https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DeleteTagsInput>
//
// * <https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DeleteTagsOutput>
func DeleteTag(cli *ec2.EC2, instances string, tags []*ec2.Tag) (err error) {
	input := &ec2.DeleteTagsInput {
		Resources: splitxs(instances, ","),
		Tags: tags,
	}
	var r *ec2.DeleteTagsOutput
	if r, err = cli.DeleteTags(input); err == nil {
		fmt.Println(r)
	}
	return
}

// * tags: タグ情報の文字列
//
// スペース区切りのタグ文字列(Ex. 'Key=foo',Value=bar Key=hoge, Value=piyo')からそれぞれ ec2.Tag を生成し
// それらのスライスとエラー情報を返す.
//
// 使用 API:
//
// * <https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#Tag>
func GenerateTags(tags string) (r []*ec2.Tag, err error) {
	for _, tag := range strings.Split(tags, " ") {
		var tagKey, tagValue string
		for _, t := range strings.Split(tag, ",") {
			val := strings.Split(t, "=")
			switch val[0] {
				case "Key": tagKey = val[1]
				case "Value": tagValue = val[1]
				default: err = errors.New("Tags parse error.")
			}
		}
		Tag := &ec2.Tag {
			Key: aws.String(tagKey),
			Value: aws.String(tagValue),
		}
		r = append(r, Tag)
	}
	return
}
