# goec2tag

```sh
go get goec2tag
```

## Usage

```sh
 ./dst/main --help
Usage of ./dst/q3:
  -addT
        Give the tag to the instance.
  -endpoint string
        Endpoint.
  -filter string
        This flag is used in conjunction with the showtags flag to filter tags by describing filter statements.
        [Example]:
                 ... -filter 'name:resource-id,values:i-xxxxxxxx i-yyyyyyyy'
  -instances string
        Instance id or instance tag name.
  -region string
        Region name (default "ap-northeast-1")
  -rmT
        Remove tag from instance.
  -showtags
        DescribeTags API operation for EC2.
         Describes one or more of the tags for your EC2 resources. filter=...
  -tags string
        Tag Key(Use Key=) and Tag Value(Use Value=)
        [Example]:
                 ... -tags='Key=foo,Value=bar Key=hoge,Value=piyo...'
```

## Example

```sh
$ ./dst/main -showtags
...

$ ./dst/main -showtags -filter "name:resource-id,values:i-xxxxxxxxxxxxxxxxx" # filtering
...

$ ./dst/main -instances=i-xxxxxxxxxxxxxxxxx -tags='Key=test,Value=hoge' -addT # adding tag
...

$ ./dst/q3 -instances=i-xxxxxxxxxxxxxxxxx -tags='Key=test,Value=hoge' -rmT # remove tag
...

```
