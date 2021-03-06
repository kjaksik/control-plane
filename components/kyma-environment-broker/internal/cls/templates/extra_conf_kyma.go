package templates

const (
	FluentBitExtraConf = `[OUTPUT]
  Name              http
  Match             *
  Host              {{.FluentdEndPoint}}
  Port              443
  HTTP_User         {{.FluentdUsername}}
  HTTP_Passwd       {{.FluentdPassword}}
  tls               true
  tls.verify        true
  tls.debug         1
  URI               /
  Format            json`
)
