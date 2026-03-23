module merge-geosite-ru-whitelist

go 1.22.6

require (
	github.com/adrg/xdg v0.5.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/v2fly/v2ray-core/v5 v5.21.0
	google.golang.org/protobuf v1.35.1
)

// Same as v2fly/v2ray-core@v5.21.0: dependency replace directives do not apply when building as a dependency.
replace github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40 => github.com/xiaokangwang/struc v0.0.0-20231031203518-0e381172f248

replace github.com/apernet/hysteria/core/v2 v2.4.5 => github.com/JimmyHuang454/hysteria/core/v2 v2.0.0-20240724161647-b3347cf6334d
