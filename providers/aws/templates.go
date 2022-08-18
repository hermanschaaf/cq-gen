package aws

// PaginatorTemplate provides an example template from AWS.
// Example values:
//  Required:
//    Package: xray
//    ServiceName: Xray
//    ReadAttribute: Groups
//  Optional:
//    Input: GetGroupsInput
//    InputAttributes: {"Example": "\"value\""}
//    NewPaginator: NewGetGroupsPaginator
const PaginatorTemplate = `
{{$NewPaginator := .NewPaginator}}
{{- if not $NewPaginator -}}
	{{$NewPaginator = print "NewGet" .ReadAttribute "Paginator"}}
{{- end -}}
{{- if .Input -}}
	var input *{{.Package}}.{{.Input}}
	{{if .InputAttributes}}
		{{- range $attr, $val := .InputAttributes -}}
		input.{{ $attr }} = {{ $val }}
		{{ end }}
	{{- end -}}
{{- end -}}
paginator := {{.Package}}.{{$NewPaginator}}(meta.(*client.Client).Services().{{.ServiceName}}, {{ if .Input }}input{{ else }}nil{{end}})
for paginator.HasMorePages() {
	v, err := paginator.NextPage(ctx)
	if err != nil {
		return diag.WrapError(err)
	}
	res <- v.{{.ReadAttribute}}
}
return nil
`
