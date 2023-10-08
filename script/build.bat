@echo off

cd ../

@rem ALib
cd ./src/ALib
go mod tidy
cd ../../


@rem Windows
go env -w GOOS=windows
@rem DockerCTL
cd ./src/DockerCTL
go mod tidy
go build -o ../../dist/dctl.exe main.go
cd ../../

@rem DomainHealthCTL
cd ./src/DomainHealthCTL
go mod tidy
go build -o ../../dist/domain.exe main.go
cd ../../

@rem GitlabCTL
cd ./src/GitlabCTL
go mod tidy
go build -o ../../dist/gctl.exe main.go
cd ../../

@rem KubernetesCTL
cd ./src/KubernetesCTL
go mod tidy
go build -o ../../dist/kctl.exe main.go
cd ../../

@rem NacosCTL
cd ./src/NacosCTL
go mod tidy
go build -o ../../dist/nctl.exe main.go
cd ../../


@rem ----------------------------------------


@rem Linux
go env -w GOOS=linux
@rem DockerCTL
cd ./src/DockerCTL
go mod tidy
go build -o ../../dist/dctl main.go
cd ../../

@rem DomainHealthCTL
cd ./src/DomainHealthCTL
go mod tidy
go build -o ../../dist/domain main.go
cd ../../


@rem GitlabCTL
cd ./src/GitlabCTL
go mod tidy
go build -o ../../dist/gctl main.go
cd ../../

@rem KubernetesCTL
cd ./src/KubernetesCTL
go mod tidy
go build -o ../../dist/kctl main.go
cd ../../

@rem NacosCTL
cd ./src/NacosCTL
go mod tidy
go build -o ../../dist/nctl main.go
cd ../../